package builders

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*GetBuildersApplication)

// GetBuildersApplication is an application service
type GetBuildersApplication struct {
	buildersStore repository.BuildersStorerReader
	imagesStore   repository.ImagesStorerReader
	selectors     map[string]repository.BuildersSelector
	output        repository.BuildersOutputter
	filterFactory FilterFactorier
}

// NewGetBuildersApplication creats a new application service
func NewGetBuildersApplication(options ...OptionsFunc) *GetBuildersApplication {

	app := &GetBuildersApplication{}
	app.Options(options...)

	return app
}

// WithBuildersStore sets the images store
func WithBuildersStore(store repository.BuildersStorerReader) OptionsFunc {
	return func(a *GetBuildersApplication) {
		a.buildersStore = store
	}
}

// WithImagesStore sets the images store
func WithImagesStore(store repository.ImagesStorerReader) OptionsFunc {
	return func(a *GetBuildersApplication) {
		a.imagesStore = store
	}
}

// WithSelector set the images selector
func WithSelector(selectors map[string]repository.BuildersSelector) OptionsFunc {
	return func(a *GetBuildersApplication) {
		a.selectors = selectors
	}
}

// WithOutput set the output
func WithOutput(output repository.BuildersOutputter) OptionsFunc {
	return func(a *GetBuildersApplication) {
		a.output = output
	}
}

// WithFilterFactory set the output
func WithFilterFactory(filterFactory FilterFactorier) OptionsFunc {
	return func(a *GetBuildersApplication) {
		a.filterFactory = filterFactory
	}
}

// Options configure the service
func (a *GetBuildersApplication) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Run method carries out the application tasks
func (a *GetBuildersApplication) Run(ctx context.Context, options *Options, optionsFunc ...OptionsFunc) error {

	var err error
	var buildersList, buildersFromImages []*builder.Builder
	var imagesList []*image.Image

	errContext := "(application::get::builders::Run)"

	if a.buildersStore == nil {
		return errors.New(errContext, "On get builders application, builders store must be provided")
	}

	if a.imagesStore == nil {
		return errors.New(errContext, "On get builders application, images store must be provided")
	}

	if a.selectors == nil {
		return errors.New(errContext, "On get builders application, selectors must be provided")
	}

	if a.filterFactory == nil {
		return errors.New(errContext, "On get builders application, filter factory must be provided")
	}

	if a.output == nil {
		return errors.New(errContext, "On get builders application, builders output must be provided")
	}

	imagesList, err = a.imagesStore.List()

	buildersFromImages, err = getInLineBuildersFromImages(imagesList)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildersList, err = a.buildersStore.List()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildersList = append(buildersList, buildersFromImages...)

	for _, filter := range options.Filter {
		operation := a.filterFactory.FilterOperation()
		err := operation.ParseFilterOpration(filter)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if operation.IsDefined() {
			selector, valid := a.selectors[operation.Attribute()]
			if !valid {
				continue
			}
			// filterOperation.operation is ignored
			buildersList, err = selector.Select(buildersList, operation.Operation(), operation.Item().(string))
			if err != nil {
				return errors.New(errContext, "Builders selection does not finish properly", err)
			}
		}
	}

	err = a.output.Output(buildersList)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func getInLineBuildersFromImages(images []*image.Image) ([]*builder.Builder, error) {

	list := []*builder.Builder{}

	for _, image := range images {
		switch image.Builder.(type) {
		case string:
			continue
		case *builder.Builder:
			list = append(list, image.Builder.(*builder.Builder))
		default:
			continue
		}
	}

	return list, nil
}
