package build

import (
	"context"
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"gopkg.in/yaml.v2"
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
	credentials    repository.CredentialsStorer
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

func WithCredentials(credentials repository.CredentialsStorer) OptionsFunc {
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

// Build starts the building process
func (a *Application) Build(ctx context.Context, buildPlan Planner, name string, version []string, options *Options, optionsFunc ...OptionsFunc) error {

	var err error
	var steps []*plan.Step
	var wg sync.WaitGroup
	buildWorkerErrs := []func() error{}

	errContext := "(build::Build)"

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
	errContext := "(build::build)"

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
	if options.ImageName != "" {
		i.Name = options.ImageName
	}

	if options.ImageRegistryHost != "" {
		i.RegistryHost = options.ImageRegistryHost
	}

	if options.ImageRegistryNamespace != "" {
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
	// add persistent vars defined service options
	// options definition has precedence over parent and image ones
	for k, v := range options.PersistentVars {
		i.PersistentVars[k] = v
	}

	if i.Parent != nil && i.Parent.PersistentVars != nil {
		for k, v := range i.Parent.PersistentVars {
			_, exist := i.PersistentVars[k]
			if !exist {
				i.PersistentVars[k] = v
			}
		}
	}

	// add persistent vars defined on the image
	for varKey, varValue := range i.PersistentVars {
		_, exist := i.PersistentVars[varKey]
		if !exist {
			i.PersistentVars[varKey] = varValue
		}
	}

	if i.Vars == nil {
		i.Vars = map[string]interface{}{}
	}
	// add vars defined on service options
	// options defintion has precedence over the image one
	for k, v := range options.Vars {
		i.Vars[k] = v
	}
	// add vars defined on the image
	for varKey, varValue := range i.Vars {
		_, exist := i.Vars[varKey]
		if !exist {
			i.Vars[varKey] = varValue
		}
	}

	if i.Labels == nil {
		i.Labels = map[string]string{}
	}
	// add lables defined on service options
	// options defintion has precedence over the image one
	for k, v := range options.Labels {
		i.Labels[k] = v
	}
	// add persistent lables defined on the image
	for k, v := range i.Labels {
		i.Labels[k] = v
	}

	if i.Parent != nil {
		pullAuth, err := a.getCredentials(i.Parent.RegistryHost)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if pullAuth != nil {
			// TODO allow other auth methods than user-pass
			buildOptions.PullAuthUsername = pullAuth.Username
			buildOptions.PullAuthPassword = pullAuth.Password
		}
	}

	pushAuth, err := a.getCredentials(i.RegistryHost)
	if err != nil {
		return errors.New(errContext, "", err)
	}
	if pushAuth != nil {
		buildOptions.PushAuthUsername = pushAuth.Username
		buildOptions.PushAuthPassword = pushAuth.Password
	}

	imageBuilder, err := a.getBuilder(i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildOptions.BuilderOptions = imageBuilder.Options
	// TODO is it populated by default?
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
	errContext := "(build::job)"

	if a.jobFactory == nil {
		return nil, errors.New(errContext, "To create a build job, is required a job factory")
	}

	return a.jobFactory.New(cmd), nil
}

func (a *Application) command(driver repository.BuildDriverer, i *image.Image, options *image.BuildDriverOptions) (job.Commander, error) {
	errContext := "(build::command)"

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

func (a *Application) getCredentials(registry string) (*credentials.UserPasswordAuth, error) {

	errContext := "(build::getCredentials)"

	if a.credentials == nil {
		return nil, errors.New(errContext, "To get credentials, is required a credentials store")
	}

	auth, _ := a.credentials.Get(registry)

	return auth, nil
}

func (a *Application) getDriver(builder *builder.Builder, options *Options) (repository.BuildDriverer, error) {
	errContext := "(build::getDriver)"

	if a.driverFactory == nil {
		return nil, errors.New(errContext, "To create a build driver, is required a driver factory")
	}

	driverName := builder.Driver
	if options.DryRun {
		driverName = "dry-run"
	}

	driver, err := a.driverFactory.Get(driverName)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return driver, nil
}

func (a *Application) getBuilder(i *image.Image) (*builder.Builder, error) {

	errContext := "(build::builder)"

	if i == nil {
		return nil, errors.New(errContext, "To generate a builder, is required an image definition")
	}

	if i.Builder == nil {
		return nil, errors.New(errContext, "To generate a builder, is required a builder definition")
	}

	if a.builders == nil {
		return nil, errors.New(errContext, "To generate a builder, is required a builder store defined on build service")
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
