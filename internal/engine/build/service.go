package build

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/schedule/job"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*Service)

// Service is an application service to build docker images
type Service struct {
	builders       BuildersStorer
	commandFactory BuildCommandFactorier
	driverFactory  DriverFactorier
	jobFactory     JobFactorier
	dispatch       Dispatcher
	semver         Semverser
	credentials    CredentialsStorer
}

// NewService creates a Service to build docker images
func NewService(options ...OptionsFunc) *Service {

	service := &Service{}
	service.Options(options...)

	return service
}

// WithBuilders sets the builders storer
func WithBuilders(builders BuildersStorer) OptionsFunc {
	return func(s *Service) {
		s.builders = builders
	}
}

// WithCommandFactory sets the command factory
func WithCommandFactory(commandFactory BuildCommandFactorier) OptionsFunc {
	return func(s *Service) {
		s.commandFactory = commandFactory
	}
}

// WithDriverFactory sets the driver factory
func WithDriverFactory(driverFactory DriverFactorier) OptionsFunc {
	return func(s *Service) {
		s.driverFactory = driverFactory
	}
}

// WithJobFactory sets the job factory
func WithJobFactory(jobFactory JobFactorier) OptionsFunc {
	return func(s *Service) {
		s.jobFactory = jobFactory
	}
}

// WithDispatch sets the dispatcher
func WithDispatch(dispatch Dispatcher) OptionsFunc {
	return func(s *Service) {
		s.dispatch = dispatch
	}
}

// WithSemver sets the semver
func WithSemver(semver Semverser) OptionsFunc {
	return func(s *Service) {
		s.semver = semver
	}
}

func WithCredentials(credentials CredentialsStorer) OptionsFunc {
	return func(s *Service) {
		s.credentials = credentials
	}
}

// Options configure the service
func (s *Service) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

// Build starts the building process
func (s *Service) Build(ctx context.Context, buildPlan Planner, name string, version []string, options *ServiceOptions, optionsFunc ...OptionsFunc) error {

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
		return errors.New(errContext, err.Error())
	}

	// configure service options before start build
	s.Options(optionsFunc...)

	// future promise which triggers the image build
	buildWorkerFunc := func(ctx context.Context, step PlanSteper, options *ServiceOptions) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			// defer notify to plans subscribed to this plan
			defer step.Notify()
			image := step.Image()

			// wait to be notified before start building
			step.Wait()

			err = s.build(ctx, image, options)
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

func (s *Service) build(ctx context.Context, image *image.Image, options *ServiceOptions) error {
	errContext := "(build::build)"

	if options == nil {
		return errors.New(errContext, "Build worker requires service options")
	}

	if image == nil {
		return errors.New(errContext, "Build worker requires an image specification")
	}

	if s.dispatch == nil {
		return errors.New(errContext, "Build worker requires a dispatcher")
	}

	if s.driverFactory == nil {
		return errors.New(errContext, "Build worker requires a driver factory")
	}

	if s.semver == nil {
		return errors.New(errContext, "Build worker requires a semver generator")
	}

	if s.credentials == nil {
		return errors.New(errContext, "Build worker requires a credentials store")
	}

	// Enrich options with image information

	// An originalOptions' copy is kept because it will be passed to children build on cascade mode.
	buildOptions := &driver.BuildDriverOptions{}

	// Image name could be overwritten by options
	if options.ImageName != "" {
		image.Name = options.ImageName
	}

	if options.ImageRegistryHost != "" {
		image.RegistryHost = options.ImageRegistryHost
	}

	if options.ImageRegistryNamespace != "" {
		image.RegistryNamespace = options.ImageRegistryNamespace
	}

	if options.EnableSemanticVersionTags {
		// semantically versions are generated by all tags and the image version
		semVerTags, _ := s.semver.GenerateSemverList(append(options.Tags, image.Version), options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			image.Tags = append(image.Tags, semVerTags...)
		}
	}

	image.Tags = append(image.Tags, options.Tags...)

	if image.PersistentVars == nil {
		image.PersistentVars = map[string]interface{}{}
	}
	// add persistent vars defined service options
	// options definition has precedence over parent and image ones
	for k, v := range options.PersistentVars {
		image.PersistentVars[k] = v
	}

	if image.Parent != nil && image.Parent.PersistentVars != nil {
		for k, v := range image.Parent.PersistentVars {
			_, exist := image.PersistentVars[k]
			if !exist {
				image.PersistentVars[k] = v
			}
		}
	}

	// add persistent vars defined on the image
	for varKey, varValue := range image.PersistentVars {
		_, exist := image.PersistentVars[varKey]
		if !exist {
			image.PersistentVars[varKey] = varValue
		}
	}

	if image.Vars == nil {
		image.Vars = map[string]interface{}{}
	}
	// add vars defined on service options
	// options defintion has precedence over the image one
	for k, v := range options.Vars {
		image.Vars[k] = v
	}
	// add vars defined on the image
	for varKey, varValue := range image.Vars {
		_, exist := image.Vars[varKey]
		if !exist {
			image.Vars[varKey] = varValue
		}
	}

	if image.Labels == nil {
		image.Labels = map[string]string{}
	}
	// add lables defined on service options
	// options defintion has precedence over the image one
	for k, v := range options.Labels {
		image.Labels[k] = v
	}
	// add persistent lables defined on the image
	for k, v := range image.Labels {
		image.Labels[k] = v
	}

	if image.Parent != nil {
		pullAuth := s.getCredentials(image.Parent.RegistryHost)
		if pullAuth != nil {
			// TODO allow other auth methods than user-pass
			buildOptions.PullAuthUsername = pullAuth.Username
			buildOptions.PullAuthPassword = pullAuth.Password
		}
	}

	pushAuth := s.getCredentials(image.RegistryHost)
	if pushAuth != nil {
		buildOptions.PushAuthUsername = pushAuth.Username
		buildOptions.PushAuthPassword = pushAuth.Password
	}

	imageBuilder, err := s.builder(image.Builder)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	buildOptions.BuilderOptions = imageBuilder.Options
	// TODO is it populated by default?
	buildOptions.BuilderVarMappings = imageBuilder.VarMapping

	driver, err := s.getDriver(imageBuilder, options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	// used by ansible driver
	buildOptions.AnsibleConnectionLocal = options.AnsibleConnectionLocal
	if options.AnsibleIntermediateContainerName != "" {
		buildOptions.AnsibleIntermediateContainerName = options.AnsibleIntermediateContainerName
	} else {
		buildOptions.AnsibleIntermediateContainerName = strings.Join([]string{"builder", imageBuilder.Driver, image.RegistryNamespace, image.Name, image.Version}, "_")
	}
	buildOptions.AnsibleInventoryPath = options.AnsibleInventoryPath
	buildOptions.AnsibleLimit = options.AnsibleLimit

	buildOptions.PullParentImage = options.PullParentImage

	buildOptions.PushImageAfterBuild = options.PushImageAfterBuild

	buildOptions.RemoveImageAfterBuild = options.RemoveImagesAfterPush

	cmd, err := s.command(driver, image, buildOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	// End options enrichment
	job, err := s.job(ctx, cmd)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	s.dispatch.Enqueue(job)

	err = job.Wait()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

func (s *Service) getCredentials(registry string) *credentials.RegistryUserPassAuth {
	auth, _ := s.credentials.GetCredentials(registry)

	return auth
}

func (s *Service) getDriver(builder *builder.Builder, options *ServiceOptions) (driver.BuildDriverer, error) {
	errContext := "(build::getDriver)"

	driverName := builder.Driver
	if options.DryRun {
		driverName = "dry-run"
	}

	driver, err := s.driverFactory.Get(driverName)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return driver, nil
}

func (s *Service) job(ctx context.Context, cmd job.Commander) (schedule.Jobber, error) {
	errContext := "(build::job)"

	if s.jobFactory == nil {
		return nil, errors.New(errContext, "To create a build job, is required a job factory")
	}

	return s.jobFactory.New(cmd), nil
}

func (s *Service) command(driver driver.BuildDriverer, image *image.Image, options *driver.BuildDriverOptions) (job.Commander, error) {
	errContext := "(build::command)"

	if s.commandFactory == nil {
		return nil, errors.New(errContext, "To create a build command, is required a command factory")
	}

	if driver == nil {
		return nil, errors.New(errContext, "To create a build command, is required a driver")
	}

	if image == nil {
		return nil, errors.New(errContext, "To create a build command, is required a image")
	}

	if options == nil {
		return nil, errors.New(errContext, "To create a build command, is required a service options")
	}

	return s.commandFactory.New(driver, image, options), nil
}

func (s *Service) builder(builderDefinition interface{}) (*builder.Builder, error) {

	errContext := "(build::builder)"

	if builderDefinition == nil {
		return nil, errors.New(errContext, "To generate a builder, is required a builder definition")
	}

	if s.builders == nil {
		return nil, errors.New(errContext, "To generate a builder, is required a builder store defined on build service")
	}

	switch builderDefinition.(type) {
	case string:
		return s.builders.Find(builderDefinition.(string))
	case *builder.Builder:
		builderAux := builderDefinition.(*builder.Builder)

		return builder.NewBuilder(builderAux.Name, builderAux.Driver, builderAux.Options, builderAux.VarMapping), nil
	default:
		// In-line builder definition
		return builder.NewBuilderFromIOReader(bytes.NewBuffer([]byte(builderDefinition.(string))))
	}
}
