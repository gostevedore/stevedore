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
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/schedule"
)

// Service is an application service to build docker images
type Service struct {
	plan Planner
	// images         ImagesStorer
	builders       BuildersStorer
	commandFactory BuildCommandFactorier
	driverFactory  DriverFactorier
	jobFactory     JobFactorier
	dispatch       Dispatcher
	semver         Semverser
	credentials    CredentialsStorer
}

// NewService creates a Service to build docker images
func NewService(plans Planner, images ImagesStorer, builders BuildersStorer, commandFactory BuildCommandFactorier, jobFactory JobFactorier, dispatch Dispatcher, semver Semverser, credentials CredentialsStorer) *Service {
	return &Service{
		plan: plans,
		// images:         images,
		builders:       builders,
		commandFactory: commandFactory,
		jobFactory:     jobFactory,
		dispatch:       dispatch,
		semver:         semver,
		credentials:    credentials,
	}
}

// Build starts the building process
func (s *Service) Build(ctx context.Context, options *ServiceOptions) error {

	var err error
	var steps []Steper
	var wg sync.WaitGroup
	buildWorkerErrs := []func() error{}

	errContext := "(build::Build)"

	if options == nil {
		return errors.New(errContext, "Options are required on build service")
	}

	if s.plan == nil {
		return errors.New(errContext, "Plan storer is required on build service")
	}

	// buildImageList, err = s.generateImagesList(options.ImageName, options.ImageVersions)
	// if err != nil {
	// 	return errors.New(errContext, err.Error())
	// }

	steps, err = s.plan.Plan(options.ImageName, options.ImageVersions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

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

	errMsg := ""
	for _, buildWorkerErr := range buildWorkerErrs {
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

// func (s *Service) generateImagesList(name string, versions []string) ([]*image.Image, error) {
// 	errContext := "(build::generateImagesList)"
// 	var list []*image.Image
// 	var imageAux *image.Image
// 	var err error

// 	if name == "" {
// 		return nil, errors.New(errContext, "Image name is required to build an image")
// 	}

// 	if versions == nil || len(versions) < 1 {
// 		list, err = s.images.All(name)
// 		if err != nil {
// 			return nil, errors.New(errContext, err.Error())
// 		}
// 	} else {
// 		for _, version := range versions {
// 			imageAux, err = s.images.Find(name, version)
// 			if err != nil {
// 				return nil, errors.New(errContext, err.Error())
// 			}
// 			list = append(list, imageAux)
// 		}
// 	}

// 	return list, nil
// }

func (s *Service) worker(ctx context.Context, image *image.Image, options *ServiceOptions) error {
	// var wg sync.WaitGroup
	// childBuildErrs := []func() error{}

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

	buildOptions.Tags = append(options.Tags, image.Tags...)

	// add persistent vars defined service options
	// options definition has precedence over the image one
	for k, v := range options.PersistentVars {
		buildOptions.PersistentVars[k] = v
	}
	// add persistent vars defined on the image
	for varKey, varValue := range image.PersistentVars {
		_, exist := buildOptions.PersistentVars[varKey]
		if !exist {
			buildOptions.PersistentVars[varKey] = varValue
		}
	}

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

	// if image.ParentName != "" {
	// 	parent, err := s.images.Find(image.ParentName, image.ParentVersion)
	// 	if err != nil {
	// 		return errors.New(errContext, err.Error())
	// 	}

	// 	buildOptions.ImageFromName = parent.Name
	// 	buildOptions.ImageFromVersion = parent.Version
	// 	buildOptions.ImageFromRegistryHost = parent.RegistryHost
	// 	buildOptions.ImageFromRegistryNamespace = parent.RegistryNamespace

	// 	pullAuth := s.getCredentials(parent.RegistryHost)
	// 	if pullAuth != nil {
	// 		// TODO allow other auth methods than user-pass
	// 		buildOptions.PullAuthUsername = pullAuth.Username
	// 		buildOptions.PullAuthPassword = pullAuth.Password
	// 	}
	// }

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

	select {
	case <-job.Done():
	case jobErr := <-job.Err():
		return errors.New(errContext, jobErr.Error())
	}

	// if options.Cascade && options.CascadeDepth != 0 {

	// 	childBuildFunc := func(ctx context.Context, options *ServiceOptions) func() error {
	// 		var err error

	// 		c := make(chan struct{}, 1)
	// 		go func() {
	// 			defer close(c)
	// 			err = s.Build(ctx, options)
	// 			wg.Done()
	// 		}()

	// 		return func() error {
	// 			<-c
	// 			return err
	// 		}
	// 	}

	// 	for childName, childVersions := range image.Children {

	// 		childServiceOptions := options.Copy()
	// 		childServiceOptions.ImageName = childName
	// 		childServiceOptions.ImageVersions = childVersions
	// 		childServiceOptions.Tags = []string{}

	// 		childServiceOptions.PersistentVars = make(map[string]interface{})
	// 		// Copy the parent persistent vars
	// 		for k, v := range buildOptions.PersistentVars {
	// 			childServiceOptions.PersistentVars[k] = v
	// 		}
	// 		childServiceOptions.Vars = map[string]interface{}{}
	// 		childServiceOptions.Labels = map[string]string{}
	// 		childServiceOptions.CascadeDepth--

	// 		wg.Add(1)
	// 		childBuildErrs = append(childBuildErrs, childBuildFunc(ctx, childServiceOptions))
	// 	}

	// 	wg.Wait()

	// 	errMsg := ""
	// 	for _, childBuildErr := range childBuildErrs {
	// 		err = childBuildErr()
	// 		if err != nil {
	// 			errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
	// 		}
	// 	}
	// 	if errMsg != "" {
	// 		return errors.New(errContext, errMsg)
	// 	}
	// }

	return nil
}

func (s *Service) getCredentials(registry string) *credentials.RegistryUserPassAuth {
	auth, _ := s.credentials.GetCredentials(registry)

	return auth
}

func (s *Service) job(ctx context.Context, cmd BuildCommander) (schedule.Jobber, error) {
	return s.jobFactory.New(cmd), nil
}

func (s *Service) command(driver driver.BuildDriverer, options *driver.BuildDriverOptions) (BuildCommander, error) {
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
	// case *builder.Builder:
	// 	return builderDefinition.(*builder.Builder), nil
	default:
		// In-line builder definition
		return builder.NewBuilderFromIOReader(bytes.NewBuffer([]byte(builderDefinition.(string))))
	}

}
