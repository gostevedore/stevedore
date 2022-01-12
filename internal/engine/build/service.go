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
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/schedule/job"
)

// Service is an application service to build docker images
type Service struct {
	plan           Planner
	builders       BuildersStorer
	commandFactory BuildCommandFactorier
	driverFactory  DriverFactorier
	jobFactory     JobFactorier
	dispatch       Dispatcher
	semver         Semverser
	credentials    CredentialsStorer
}

// NewService creates a Service to build docker images
func NewService(plans Planner, builders BuildersStorer, commandFactory BuildCommandFactorier, driverFactory DriverFactorier, jobFactory JobFactorier, dispatch Dispatcher, semver Semverser, credentials CredentialsStorer) *Service {

	return &Service{
		plan:           plans,
		builders:       builders,
		commandFactory: commandFactory,
		driverFactory:  driverFactory,
		jobFactory:     jobFactory,
		dispatch:       dispatch,
		semver:         semver,
		credentials:    credentials,
	}
}

// Build starts the building process
func (s *Service) Build(ctx context.Context, name string, version []string, options *ServiceOptions) error {

	var err error
	var steps []*plan.Step
	var wg sync.WaitGroup
	buildWorkerErrs := []func() error{}

	errContext := "(build::Build)"

	if options == nil {
		return errors.New(errContext, "To build an image, service options are required")
	}

	if s.plan == nil {
		return errors.New(errContext, "To build an image, execution plan is required")
	}

	steps, err = s.plan.Plan(name, version)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	// future promise which triggers the image build
	buildWorkerFunc := func(ctx context.Context, step Steper, options *ServiceOptions) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			// defer notify to plans subscribed to this plan
			defer step.Notify()
			image := step.Image()

			// wait to be notified before start building
			step.Wait()

			err = s.worker(ctx, image, options)
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

func (s *Service) worker(ctx context.Context, image *image.Image, options *ServiceOptions) error {
	errContext := "(build::worker)"

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
	if options.ImageName == "" {
		buildOptions.ImageName = image.Name
	} else {
		buildOptions.ImageName = options.ImageName
	}
	buildOptions.ImageVersion = image.Version

	if options.ImageRegistryHost != "" {
		buildOptions.RegistryHost = options.ImageRegistryHost
	} else {
		// TODO image domain need to ensure that image registry host is set
		buildOptions.RegistryHost = image.RegistryHost
	}

	if options.ImageRegistryNamespace != "" {
		buildOptions.RegistryNamespace = options.ImageRegistryNamespace
	} else {
		buildOptions.RegistryNamespace = image.RegistryNamespace
	}

	if options.EnableSemanticVersionTags {
		// semantically versions are generated by all tags and the image version
		semVerTags, _ := s.semver.GenerateSemverList(append(options.Tags, buildOptions.ImageVersion), options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			buildOptions.Tags = append(buildOptions.Tags, semVerTags...)
		}
	}

	buildOptions.Tags = append(buildOptions.Tags, options.Tags...)
	buildOptions.Tags = append(buildOptions.Tags, image.Tags...)

	buildOptions.PersistentVars = map[string]interface{}{}
	// add persistent vars defined service options
	// options definition has precedence over parent and image ones
	for k, v := range options.PersistentVars {
		buildOptions.PersistentVars[k] = v
	}

	if image.Parent != nil {
		for k, v := range image.Parent.PersistentVars {
			_, exist := buildOptions.PersistentVars[k]
			if !exist {
				buildOptions.PersistentVars[k] = v
			}
		}
	}

	// add persistent vars defined on the image
	for varKey, varValue := range image.PersistentVars {
		_, exist := buildOptions.PersistentVars[varKey]
		if !exist {
			buildOptions.PersistentVars[varKey] = varValue
		}
	}

	buildOptions.Vars = map[string]interface{}{}
	// add vars defined on service options
	// options defintion has precedence over the image one
	for k, v := range options.Vars {
		buildOptions.Vars[k] = v
	}
	// add vars defined on the image
	for varKey, varValue := range image.Vars {
		_, exist := buildOptions.Vars[varKey]
		if !exist {
			buildOptions.Vars[varKey] = varValue
		}
	}

	buildOptions.Labels = map[string]string{}
	// add lables defined on service options
	// options defintion has precedence over the image one
	for k, v := range options.Labels {
		buildOptions.Labels[k] = v
	}
	// add persistent lables defined on the image
	for k, v := range image.Labels {
		buildOptions.Labels[k] = v
	}

	if image.Parent != nil {
		buildOptions.ImageFromName = image.Parent.Name
		buildOptions.ImageFromVersion = image.Parent.Version
		buildOptions.ImageFromRegistryHost = image.Parent.RegistryHost
		buildOptions.ImageFromRegistryNamespace = image.Parent.RegistryNamespace

		pullAuth := s.getCredentials(image.Parent.RegistryHost)
		if pullAuth != nil {
			// TODO allow other auth methods than user-pass
			buildOptions.PullAuthUsername = pullAuth.Username
			buildOptions.PullAuthPassword = pullAuth.Password
		}
	}

	pushAuth := s.getCredentials(buildOptions.RegistryHost)
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

	driver, err := s.driverFactory.Get(imageBuilder.Driver)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	buildOptions.BuilderName = strings.Join([]string{"builder", imageBuilder.Driver, buildOptions.RegistryNamespace, buildOptions.ImageName, buildOptions.ImageVersion}, "_")

	// used by ansible driver
	buildOptions.ConnectionLocal = options.ConnectionLocal

	buildOptions.PullParentImage = options.PullParentImage

	buildOptions.PushImageAfterBuild = options.PushImageAfterBuild

	buildOptions.RemoveImageAfterBuild = options.RemoveAfterBuild

	cmd, err := s.command(driver, buildOptions)
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

func (s *Service) job(ctx context.Context, cmd job.Commander) (schedule.Jobber, error) {
	errContext := "(build::job)"

	if s.jobFactory == nil {
		return nil, errors.New(errContext, "To create a build job, is required a job factory")
	}

	return s.jobFactory.New(cmd), nil
}

func (s *Service) command(driver driver.BuildDriverer, options *driver.BuildDriverOptions) (job.Commander, error) {
	errContext := "(build::command)"

	if s.commandFactory == nil {
		return nil, errors.New(errContext, "To create a build command, is required a command factory")
	}

	if driver == nil {
		return nil, errors.New(errContext, "To create a build command, is required a driver")
	}

	if options == nil {
		return nil, errors.New(errContext, "To create a build command, is required a service options")
	}

	return s.commandFactory.New(driver, options), nil
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
