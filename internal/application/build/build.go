package build

import (
	"context"
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	authmethodbasic "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"gopkg.in/yaml.v3"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*Application)

// Application is an application service to build docker images
type Application struct {
	builders       repository.BuildersStorer
	commandFactory BuildCommandFactorier
	driverFactory  DriverFactorier
	jobFactory     JobFactorier
	dispatch       Dispatcher
	semver         Semverser
	credentials    repository.AuthFactorier
}

// NewApplication creates a Service to build docker images
func NewApplication(options ...OptionsFunc) *Application {

	service := &Application{}
	service.Options(options...)

	return service
}

// WithBuilders sets the builders storer
func WithBuilders(builders repository.BuildersStorer) OptionsFunc {
	return func(a *Application) {
		a.builders = builders
	}
}

// WithCommandFactory sets the command factory
func WithCommandFactory(commandFactory BuildCommandFactorier) OptionsFunc {
	return func(a *Application) {
		a.commandFactory = commandFactory
	}
}

// WithDriverFactory sets the driver factory
func WithDriverFactory(driverFactory DriverFactorier) OptionsFunc {
	return func(a *Application) {
		a.driverFactory = driverFactory
	}
}

// WithJobFactory sets the job factory
func WithJobFactory(jobFactory JobFactorier) OptionsFunc {
	return func(a *Application) {
		a.jobFactory = jobFactory
	}
}

// WithDispatch sets the dispatcher
func WithDispatch(dispatch Dispatcher) OptionsFunc {
	return func(a *Application) {
		a.dispatch = dispatch
	}
}

// WithSemver sets the semver
func WithSemver(semver Semverser) OptionsFunc {
	return func(a *Application) {
		a.semver = semver
	}
}

func WithCredentials(credentials repository.AuthFactorier) OptionsFunc {
	return func(a *Application) {
		a.credentials = credentials
	}
}

// Options configure the service
func (a *Application) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Build method carries out the application tasks
func (a *Application) Build(ctx context.Context, buildPlan Planner, name string, version []string, options *Options, optionsFunc ...OptionsFunc) error {

	var err error
	var steps []*plan.Step
	var wg sync.WaitGroup
	buildWorkerErrs := []func() error{}

	errContext := "(application::build::Build)"

	if options == nil {
		return errors.New(errContext, "To build an image, service options are required")
	}

	if buildPlan == nil {
		return errors.New(errContext, "To build an image, a build plan is required")
	}

	steps, err = buildPlan.Plan(name, version)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	// configure service options before start build
	a.Options(optionsFunc...)

	// future promise which triggers the image build
	buildWorkerFunc := func(ctx context.Context, step PlanSteper, options *Options) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			// defer notify to plans subscribed to this plan
			defer step.Notify()
			image := step.Image()

			// wait to be notified before start building
			step.Wait()

			err = a.build(ctx, image, options)
			wg.Done()
		}()

		return func() error {
			<-c
			return err
		}
	}

	// execute build workers as future promises
	for _, step := range steps {
		wg.Add(1)
		buildWorkerErrs = append(buildWorkerErrs, buildWorkerFunc(ctx, step, options))
	}

	wg.Wait()

	// Wait for all workers to finish
	errMsg := ""
	for _, buildWorkerErr := range buildWorkerErrs {
		// it is blocking
		err = buildWorkerErr()
		if err != nil {
			errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errContext, errMsg)
	}

	return nil
}

func (a *Application) build(ctx context.Context, i *image.Image, options *Options) error {
	var parent *image.Image
	errContext := "(application::build::build)"

	if options == nil {
		return errors.New(errContext, "Build worker requires service options")
	}

	if i == nil {
		return errors.New(errContext, "Build worker requires an image specification")
	}

	if a.dispatch == nil {
		return errors.New(errContext, "Build worker requires a dispatcher")
	}

	if a.driverFactory == nil {
		return errors.New(errContext, "Build worker requires a driver factory")
	}

	if a.semver == nil {
		return errors.New(errContext, "Build worker requires a semver generator")
	}

	if a.credentials == nil {
		return errors.New(errContext, "Build worker requires a credentials store")
	}

	// Enrich options with image information

	// An originalOptions' copy is kept because it will be passed to children build on cascade mode.
	buildOptions := &image.BuildDriverOptions{}

	// Image name could be overwritten by options
	if options.ImageName != image.UndefinedStringValue {
		i.Name = options.ImageName
	}

	if options.ImageRegistryHost != image.UndefinedStringValue {
		i.RegistryHost = options.ImageRegistryHost
	}

	if options.ImageRegistryNamespace != image.UndefinedStringValue {
		i.RegistryNamespace = options.ImageRegistryNamespace
	}

	if options.EnableSemanticVersionTags {
		// semantically versions are generated by all tags and the image version
		semVerTags, _ := a.semver.GenerateSemverList(append(options.Tags, i.Version), options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			i.Tags = append(i.Tags, semVerTags...)
		}
	}

	i.Tags = append(i.Tags, options.Tags...)

	if i.PersistentVars == nil {
		i.PersistentVars = map[string]interface{}{}
	}

	for k, v := range options.PersistentVars {
		i.PersistentVars[k] = v
	}

	if i.Vars == nil {
		i.Vars = map[string]interface{}{}
	}

	for k, v := range options.Vars {
		i.Vars[k] = v
	}

	if i.Labels == nil {
		i.Labels = map[string]string{}
	}

	for k, v := range options.Labels {
		i.Labels[k] = v
	}

	if i.PersistentLabels == nil {
		i.PersistentLabels = map[string]string{}
	}

	for k, v := range options.PersistentLabels {
		i.PersistentLabels[k] = v
	}

	if options.ImageFromName != image.UndefinedStringValue {
		if i.Parent == nil {
			if parent == nil {
				parent = &image.Image{}
			}
			parent.Name = options.ImageFromName
		} else {
			i.Parent.Name = options.ImageFromName
		}
	}

	if options.ImageFromVersion != image.UndefinedStringValue {
		if i.Parent == nil {
			if parent == nil {
				parent = &image.Image{}
			}
			parent.Version = options.ImageFromVersion
		} else {
			i.Parent.Version = options.ImageFromVersion
		}
	}

	if options.ImageFromRegistryHost != image.UndefinedStringValue {
		if i.Parent == nil {
			if parent == nil {
				parent = &image.Image{}
			}
			parent.RegistryHost = options.ImageFromRegistryHost
		} else {
			i.Parent.RegistryHost = options.ImageFromRegistryHost
		}
	}

	if options.ImageFromRegistryNamespace != image.UndefinedStringValue {
		if i.Parent == nil {
			if parent == nil {
				parent = &image.Image{}
			}
			parent.RegistryNamespace = options.ImageFromRegistryNamespace
		} else {
			i.Parent.RegistryNamespace = options.ImageFromRegistryNamespace
		}
	}

	if i.Parent == nil && parent != nil {
		i.Parent = parent
	}

	if i.Parent != nil && i.Parent.RegistryHost != "" && i.Parent.RegistryHost != image.UndefinedStringValue {
		auth, err := a.getCredentials(i.Parent.RegistryHost)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if auth != nil {
			pullAuth, isBasicAuth := auth.(*authmethodbasic.BasicAuthMethod)
			if !isBasicAuth {
				return errors.New(errContext, fmt.Sprintf("Invalid credentials method for '%s'. Found '%s' when is expected basic auth method", i.Parent.RegistryHost, auth.Name()))
			}

			buildOptions.PullAuthUsername = pullAuth.Username
			buildOptions.PullAuthPassword = pullAuth.Password
		}
	}

	if i.RegistryHost != image.UndefinedStringValue {
		auth, err := a.getCredentials(i.RegistryHost)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if auth != nil {
			pushAuth, isBasicAuth := auth.(*authmethodbasic.BasicAuthMethod)
			if !isBasicAuth {
				return errors.New(errContext, fmt.Sprintf("Invalid credentials method for '%s'. Found '%s' when is expected basic auth method", i.RegistryHost, auth.Name()))
			}

			buildOptions.PushAuthUsername = pushAuth.Username
			buildOptions.PushAuthPassword = pushAuth.Password
		}
	}

	imageBuilder, err := a.getBuilder(i)
	if err != nil {
		return errors.New(errContext, "", err) // TODO is it populated by default?
	}

	buildOptions.BuilderOptions = imageBuilder.Options
	buildOptions.BuilderVarMappings = imageBuilder.VarMapping

	driver, err := a.getDriver(imageBuilder, options)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	// used by ansible driver
	buildOptions.AnsibleConnectionLocal = options.AnsibleConnectionLocal
	if options.AnsibleIntermediateContainerName != "" {
		buildOptions.AnsibleIntermediateContainerName = options.AnsibleIntermediateContainerName
	} else {
		buildOptions.AnsibleIntermediateContainerName = strings.Join([]string{"builder", imageBuilder.Driver, i.RegistryNamespace, i.Name, i.Version}, "_")
	}
	buildOptions.AnsibleInventoryPath = options.AnsibleInventoryPath
	buildOptions.AnsibleLimit = options.AnsibleLimit

	buildOptions.PullParentImage = options.PullParentImage

	buildOptions.PushImageAfterBuild = options.PushImageAfterBuild

	buildOptions.RemoveImageAfterBuild = options.RemoveImagesAfterPush

	err = i.Sanetize()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	cmd, err := a.command(driver, i, buildOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	// End options enrichment
	job, err := a.job(ctx, cmd)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	a.dispatch.Enqueue(job)

	err = job.Wait()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (a *Application) job(ctx context.Context, cmd job.Commander) (scheduler.Jobber, error) {
	errContext := "(application::build::job)"

	if a.jobFactory == nil {
		return nil, errors.New(errContext, "To create a build job, is required a job factory")
	}

	return a.jobFactory.New(cmd), nil
}

func (a *Application) command(driver repository.BuildDriverer, i *image.Image, options *image.BuildDriverOptions) (job.Commander, error) {
	errContext := "(application::build::command)"

	if a.commandFactory == nil {
		return nil, errors.New(errContext, "To create a build command, is required a command factory")
	}

	if driver == nil {
		return nil, errors.New(errContext, "To create a build command, is required a driver")
	}

	if i == nil {
		return nil, errors.New(errContext, "To create a build command, is required a image")
	}

	if options == nil {
		return nil, errors.New(errContext, "To create a build command, is required a service options")
	}

	return a.commandFactory.New(driver, i, options), nil
}

func (a *Application) getCredentials(registry string) (repository.AuthMethodReader, error) {

	errContext := "(application::build::getCredentials)"

	if a.credentials == nil {
		return nil, errors.New(errContext, "To get credentials, is required a credentials store")
	}

	auth, _ := a.credentials.Get(registry)

	return auth, nil
}

func (a *Application) getDriver(builder *builder.Builder, options *Options) (repository.BuildDriverer, error) {
	errContext := "(application::build::getDriver)"

	var factoryFunc factory.BuildDriverFactoryFunc
	var err error
	var driver repository.BuildDriverer

	if a.driverFactory == nil {
		return nil, errors.New(errContext, "To create a build driver, is required a driver factory")
	}

	driverName := builder.Driver

	factoryFunc, err = a.driverFactory.Get(driverName)
	if err != nil {
		factoryFunc, err = a.driverFactory.Get(image.DefaultDriverName)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
	}

	driver, err = factoryFunc()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return driver, nil
}

func (a *Application) getBuilder(i *image.Image) (*builder.Builder, error) {

	errContext := "(application::build::builder)"

	if i == nil {
		return nil, errors.New(errContext, "To generate a builder, is required an image definition")
	}

	if a.builders == nil {
		return nil, errors.New(errContext, "To generate a builder, is required a builder store defined on build service")
	}

	if i.Builder == nil {
		return builder.NewBuilder(i.Name, image.DefaultDriverName, nil, nil), nil
	}

	switch i.Builder.(type) {
	case string:
		return a.builders.Find(i.Builder.(string))
	case *builder.Builder:
		builderAux := i.Builder.(*builder.Builder)
		return builder.NewBuilder(builderAux.Name, builderAux.Driver, builderAux.Options, builderAux.VarMapping), nil
	default:
		builderDefinitionBytes, err := yaml.Marshal(i.Builder)
		if err != nil {
			return nil, errors.New(errContext, fmt.Sprintf("There is an error marshaling '%s:%s' builder", i.Name, i.Version), err)
		}

		b, err := builder.NewBuilderFromByteArray(builderDefinitionBytes)
		if err != nil {
			return nil, errors.New(errContext, fmt.Sprintf("There is an error creating the builder for '%s:%s'", i.Name, i.Version), err)
		}

		return b, nil
	}
}
