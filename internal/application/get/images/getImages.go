package images

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*GetImagesApplication)

// GetImagesApplication is an application service
type GetImagesApplication struct {
	store     repository.ImagesStorerReader
	selectors map[string]repository.ImagesSelector
	output    repository.ImagesOutputter
}

// NewGetImagesApplication creats a new application service
func NewGetImagesApplication(options ...OptionsFunc) *GetImagesApplication {

	service := &GetImagesApplication{}
	service.Options(options...)

	return service
}

// WithStore sets the images store
func WithStore(store repository.ImagesStorerReader) OptionsFunc {
	return func(a *GetImagesApplication) {
		a.store = store
	}
}

// WithSelector set the images selector
func WithSelector(selectors map[string]repository.ImagesSelector) OptionsFunc {
	return func(a *GetImagesApplication) {
		a.selectors = selectors
	}
}

// WithOutput set the output
func WithOutput(output repository.ImagesOutputter) OptionsFunc {
	return func(a *GetImagesApplication) {
		a.output = output
	}
}

// Options configure the service
func (a *GetImagesApplication) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Build starts the building process
func (a *GetImagesApplication) Run(ctx context.Context, options *Options, optionsFunc ...OptionsFunc) error {

	var err error
	var imagesList []*image.Image
	errContext := "(application::get::images::Run)"

	if a.store == nil {
		return errors.New(errContext, "On get images application, images store must be provided")
	}

	if a.output == nil {
		return errors.New(errContext, "On get images application, images output must be provided")
	}

	imagesList, err = a.store.List()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	for _, filter := range options.Filter {
		filterOperation := NewFilterOperation(filter)
		if filterOperation.IsDefined() {
			selector, valid := a.selectors[filterOperation.attribute]
			if !valid {
				continue
			}
			// filterOperation.operation is ignored
			imagesList, err = selector.Select(imagesList, filterOperation.operation, filterOperation.item.(string))
			if err != nil {
				return errors.New(errContext, "Images selection does not finish properly", err)
			}
		}
	}

	err = a.output.Output(imagesList)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
