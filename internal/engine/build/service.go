package build

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/image"
)

// Service is an application service to build docker images
type Service struct {
	images         ImagesStorer
	builders       BuildersStorer
	commandFactory BuildCommandFactorier
	jobFactory     JobFactorier
	dispatch       Dispatcher
}

// NewService creates a Service to build docker images
func NewService(images ImagesStorer, builders BuildersStorer, commandFactory BuildCommandFactorier, jobFactory JobFactorier, dispatch Dispatcher) *Service {
	return &Service{
		images:         images,
		builders:       builders,
		commandFactory: commandFactory,
		jobFactory:     jobFactory,
		dispatch:       dispatch,
	}
}

// Build starts the building process
func (s *Service) Build(ctx context.Context, options *ServiceOptions) error {

	var err error
	var buildImageList []*image.Image
	var wg sync.WaitGroup
	buildWorkerErrs := []func() error{}

	errContext := "(build::Build)"

	if options == nil {
		return errors.New(errContext, "Options are required on build service")
	}

	if s.dispatch == nil {
		return errors.New(errContext, "Build worker requires a dispatcher")
	}

	buildImageList, err = s.generateImagesList(options.ImageName, options.ImageVersions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	buildWorkerFunc := func(ctx context.Context, image *image.Image, options *ServiceOptions) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			err = s.worker(ctx, image, options)
			wg.Done()
		}()

		return func() error {
			<-c
			return err
		}
	}

	for _, image := range buildImageList {
		wg.Add(1)
		buildWorkerErrs = append(buildWorkerErrs, buildWorkerFunc(ctx, image, options))

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

func (s *Service) generateImagesList(name string, versions []string) ([]*image.Image, error) {
	errContext := "(build::generateImagesList)"
	var list []*image.Image
	var imageAux *image.Image
	var err error

	if name == "" {
		return nil, errors.New(errContext, "Image name is required to build an image")
	}

	if versions == nil || len(versions) < 1 {
		list, err = s.images.All(name)
		if err != nil {
			return nil, errors.New(errContext, err.Error())
		}
	} else {
		for _, version := range versions {
			imageAux, err = s.images.Find(name, version)
			if err != nil {
				return nil, errors.New(errContext, err.Error())
			}
			list = append(list, imageAux)
		}
	}

	return list, nil
}

func (s *Service) worker(ctx context.Context, image *image.Image, options *ServiceOptions) error {
	var wg sync.WaitGroup
	childBuildErrs := []func() error{}

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

	job, err := s.job(ctx, image, options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	s.dispatch.Enqueue(job)

	select {
	case <-job.Done():
	case jobErr := <-job.Err():
		return errors.New(errContext, jobErr.Error())
	}

	if options.Cascade && options.CascadeDepth != 0 {

		childBuildFunc := func(ctx context.Context, options *ServiceOptions) func() error {
			var err error

			c := make(chan struct{}, 1)
			go func() {
				defer close(c)
				err = s.Build(ctx, options)
				wg.Done()
			}()

			return func() error {
				<-c
				return err
			}
		}

		for childName, childVersions := range image.Children {

			// TODO copy options
			childServiceOptions := options.Copy()
			childServiceOptions.ImageName = childName
			childServiceOptions.ImageVersions = childVersions
			childServiceOptions.Tags = []string{}
			childServiceOptions.Vars = map[string]interface{}{}
			childServiceOptions.Labels = map[string]string{}
			childServiceOptions.CascadeDepth--

			wg.Add(1)
			childBuildErrs = append(childBuildErrs, childBuildFunc(ctx, childServiceOptions))
		}

		wg.Wait()

		errMsg := ""
		for _, childBuildErr := range childBuildErrs {
			err = childBuildErr()
			if err != nil {
				errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
			}
		}
		if errMsg != "" {
			return errors.New(errContext, errMsg)
		}
	}

	return nil
}

func (s *Service) job(ctx context.Context, image *image.Image, options *ServiceOptions) (Jobber, error) {

	errContext := "(build::job)"

	cmd, err := s.command(ctx, image, options)

	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return s.jobFactory.New(cmd), nil
}

func (s *Service) command(ctx context.Context, image *image.Image, options *ServiceOptions) (BuildCommander, error) {
	errContext := "(build::Command)"

	if s.commandFactory == nil {
		return nil, errors.New(errContext, "Build worker requires a command factory")
	}

	b, err := s.builder(image.Builder)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	buildDriverOptions := &driver.BuildDriverOptions{
		BuilderOptions: b.Options,
	}
	_ = buildDriverOptions

	return nil, nil
}

func (s *Service) builder(builder interface{}) (*builder.Builder, error) {
	return nil, nil
}
