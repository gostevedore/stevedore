package images

import (
	"context"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/images"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/get/images"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	filter "github.com/gostevedore/stevedore/internal/infrastructure/filters/images"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	output "github.com/gostevedore/stevedore/internal/infrastructure/output/images"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	store "github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *GetImagesEntrypoint)

// GetImagesEntrypoint defines the entrypoint for the application
type GetImagesEntrypoint struct {
	fs            afero.Fs
	writer        io.Writer
	compatibility Compatibilitier
}

// NewGetImagesEntrypoint returns a new entrypoint
func NewGetImagesEntrypoint(opts ...OptionsFunc) *GetImagesEntrypoint {
	e := &GetImagesEntrypoint{}
	e.Options(opts...)

	return e
}

// WithWriter sets the writer for the entrypoint
func WithWriter(w io.Writer) OptionsFunc {
	return func(e *GetImagesEntrypoint) {
		e.writer = w
	}
}

// WithFileSystem sets the writer for the entrypoint
func WithFileSystem(fs afero.Fs) OptionsFunc {
	return func(e *GetImagesEntrypoint) {
		e.fs = fs
	}
}

// WithCompatibility set the
func WithCompatibility(c Compatibilitier) OptionsFunc {
	return func(e *GetImagesEntrypoint) {
		e.compatibility = c
	}
}

// Options provides the options for the entrypoint
func (e *GetImagesEntrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute is a pseudo-main method for the command
func (e *GetImagesEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, inputEntrypointOptions *Options, inputHandlerOptions *handler.Options) error {

	var err error
	var getImagesService *application.GetImagesApplication
	var getImagesHandler *handler.GetImagesHandler
	var imageRender *render.ImageRender
	var graphTemplateFactory *graph.GraphTemplateFactory
	var imagesGraphTemplatesStore *imagesgraphtemplate.ImagesGraphTemplate
	var imagesStore *store.Store

	errContext := "(get::images::entrypoint::Execute)"

	imageRender, err = e.createImageRender(now.NewNow())
	if err != nil {
		return errors.New(errContext, "", err)
	}

	graphTemplateFactory, err = e.createGraphTemplateFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	imagesGraphTemplatesStore, err = e.createImagesGraphTemplatesStorer(graphTemplateFactory)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	imagesStore, err = e.createImagesStore(conf, imageRender, imagesGraphTemplatesStore)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	writer := e.createtOutput(inputEntrypointOptions)

	// missing: assign outputer (tree, text) and assign store, assign filters
	getImagesService = application.NewGetImagesApplication(
		application.WithStore(imagesStore),
		application.WithSelector(map[string]repository.ImagesSelector{
			image.NameFilterAttribute:              filter.NewImageNameFilter(),
			image.VersionFilterAttribute:           filter.NewImageVersionFilter(),
			image.RegistryHostFilterAttribute:      filter.NewImageRegistryFilter(),
			image.RegistryNamespaceFilterAttribute: filter.NewImageNamespaceFilter(),
		}),
		application.WithOutput(writer),
	)
	getImagesHandler = handler.NewGetImagesHandler(
		handler.WithApplication(getImagesService),
	)

	err = getImagesHandler.Handler(ctx, inputHandlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// // prepareEntrypointOptions
// func (e *GetImagesEntrypoint) prepareEntrypointOptions(conf *configuration.Configuration, inputEntrypointOptions *Options) (*Options, error) {
// 	options := &Options{}
// 	errContext := "(get::images::entrypoint::prepareEntrypointOptions)"

// 	if conf == nil {
// 		return nil, errors.New(errContext, "To prepare get images entrypoint options, configuration is required")
// 	}

// 	if inputEntrypointOptions == nil {
// 		return nil, errors.New(errContext, "To prepare get images entrypoint options, entrypoint options are required")
// 	}

// 	return options, nil
// }

// func (e *GetImagesEntrypoint) prepareHandlerOptions(conf *configuration.Configuration, inputHandlerOptions *handler.Options) (*handler.Options, error) {
// 	options := &handler.Options{}
// 	errContext := "(get::images::entrypoint::prepareHandlerOptions)"

// 	if conf == nil {
// 		return nil, errors.New(errContext, "To prepare handler options in get images entrypoint, configuration is required")
// 	}

// 	if inputHandlerOptions == nil {
// 		return nil, errors.New(errContext, "To prepare handler options in get images entrypoint, handler options are required")
// 	}

// 	return options, nil
// }

func (e *GetImagesEntrypoint) createImageRender(now render.Nower) (*render.ImageRender, error) {
	errContext := "(get::images::entrypoint::createImageRender)"

	if now == nil {
		return nil, errors.New(errContext, "To create an image render in get images entrypoint, a nower is required")
	}

	return render.NewImageRender(now), nil
}

func (e *GetImagesEntrypoint) createImagesStore(conf *configuration.Configuration, render repository.Renderer, graph imagesconfiguration.ImagesGraphTemplatesStorer) (*store.Store, error) {

	errContext := "(get::images::entrypoint::createImagesStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create an images store in get images entrypoint, a filesystem is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create an images store in get images entrypoint, configuration is required")
	}

	if render == nil {
		return nil, errors.New(errContext, "To create an images store in get images entrypoint, image render is required")
	}

	if graph == nil {
		return nil, errors.New(errContext, "To create an images store in get images entrypoint, images graph templates storer is required")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create an images store in get images entrypoint, compatibility is required")
	}

	if conf.ImagesPath == "" {
		return nil, errors.New(errContext, "To create an images store in get images entrypoint, images path must be provided in configuration")
	}

	store := store.NewStore(render)
	imagesConfiguration := imagesconfiguration.NewImagesConfiguration(e.fs, graph, store, render, e.compatibility)
	err := imagesConfiguration.LoadImagesToStore(conf.ImagesPath)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return store, nil
}

func (e *GetImagesEntrypoint) createImagesGraphTemplatesStorer(factory *graph.GraphTemplateFactory) (*imagesgraphtemplate.ImagesGraphTemplate, error) {
	errContext := "(get::images::entrypoint::createImagesGraphTemplatesStorer)"

	if factory == nil {
		return nil, errors.New(errContext, "To create an images graph templates storer in get images entrypoint, a graph template factory is required")
	}

	graphtemplate := imagesgraphtemplate.NewImagesGraphTemplate(factory)

	return graphtemplate, nil
}

func (e *GetImagesEntrypoint) createGraphTemplateFactory() (*graph.GraphTemplateFactory, error) {
	return graph.NewGraphTemplateFactory(false), nil
}

func (e *GetImagesEntrypoint) createtOutput(options *Options) repository.ImagesOutputter {

	if options.Tree {
		return e.createTreeOutput()
	}

	return e.createPlainTextOutput()
}

func (e *GetImagesEntrypoint) createPlainTextOutput() repository.ImagesOutputter {
	c := console.NewConsole(e.writer, nil)
	output := output.NewPlainOutput(
		output.WithWriter(c),
	)

	return output
}

func (e *GetImagesEntrypoint) createTreeOutput() repository.ImagesOutputter {
	return nil
}
