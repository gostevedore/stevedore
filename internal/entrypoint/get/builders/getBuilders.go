package builders

import (
	"context"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/builders"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/get/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	buildersconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/builders"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	filter "github.com/gostevedore/stevedore/internal/infrastructure/filters/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/filters/operation"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	output "github.com/gostevedore/stevedore/internal/infrastructure/output/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	buildersstore "github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	imagesstore "github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *GetBuildersEntrypoint)

// GetBuildersEntrypoint defines the entrypoint for the application
type GetBuildersEntrypoint struct {
	fs            afero.Fs
	writer        io.Writer
	compatibility Compatibilitier
}

// NewGetBuildersEntrypoint returns a new entrypoint
func NewGetBuildersEntrypoint(opts ...OptionsFunc) *GetBuildersEntrypoint {
	e := &GetBuildersEntrypoint{}
	e.Options(opts...)

	return e
}

// WithWriter sets the writer for the entrypoint
func WithWriter(w io.Writer) OptionsFunc {
	return func(e *GetBuildersEntrypoint) {
		e.writer = w
	}
}

// WithFileSystem sets the writer for the entrypoint
func WithFileSystem(fs afero.Fs) OptionsFunc {
	return func(e *GetBuildersEntrypoint) {
		e.fs = fs
	}
}

// WithCompatibility set the
func WithCompatibility(c Compatibilitier) OptionsFunc {
	return func(e *GetBuildersEntrypoint) {
		e.compatibility = c
	}
}

// Options provides the options for the entrypoint
func (e *GetBuildersEntrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute is a pseudo-main method for the command
func (e *GetBuildersEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, inputHandlerOptions *handler.Options) error {

	var err error
	var getBuildersApplication *application.GetBuildersApplication
	var getBuildersHandler *handler.GetBuildersHandler
	var imageRender *render.ImageRender
	var graphTemplateFactory *graph.GraphTemplateFactory
	var imagesGraphTemplatesStore *imagesgraphtemplate.ImagesGraphTemplate
	var imagesStore *imagesstore.Store
	var buildersStore *buildersstore.Store
	var writer repository.BuildersOutputter
	var filterFactory *operation.FilterOperationFactory

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

	buildersStore, err = e.createBuildersStore(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	writer, err = e.createOutput()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	filterFactory = operation.NewFilterOperationFactory()

	getBuildersApplication = application.NewGetBuildersApplication(
		application.WithImagesStore(imagesStore),
		application.WithBuildersStore(buildersStore),
		application.WithSelector(map[string]repository.BuildersSelector{
			builder.NameFilterAttribute:   filter.NewBuilderNameFilter(),
			builder.DriverFilterAttribute: filter.NewBuilderDriverFilter(),
		}),
		application.WithFilterFactory(filterFactory),
		application.WithOutput(writer),
	)
	getBuildersHandler = handler.NewGetBuildersHandler(
		handler.WithApplication(getBuildersApplication),
	)

	err = getBuildersHandler.Handler(ctx, inputHandlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (e *GetBuildersEntrypoint) createImageRender(now render.Nower) (*render.ImageRender, error) {
	errContext := "(entrypoint::get::builders::createImageRender)"

	if now == nil {
		return nil, errors.New(errContext, "To create an image render in get builders entrypoint, a nower is required")
	}

	return render.NewImageRender(now), nil
}

func (e *GetBuildersEntrypoint) createImagesStore(conf *configuration.Configuration, render repository.Renderer, graph imagesconfiguration.ImagesGraphTemplatesStorer) (*imagesstore.Store, error) {

	errContext := "(entrypoint::get::builders::createImagesStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create an images store in get builders entrypoint, a filesystem is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create an images store in get builders entrypoint, configuration is required")
	}

	if render == nil {
		return nil, errors.New(errContext, "To create an images store in get builders entrypoint, image render is required")
	}

	if graph == nil {
		return nil, errors.New(errContext, "To create an images store in get builders entrypoint, images graph templates storer is required")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create an images store in get builders entrypoint, compatibility is required")
	}

	if conf.ImagesPath == "" {
		return nil, errors.New(errContext, "To create an images store in get builders entrypoint, images path must be provided in configuration")
	}

	store := imagesstore.NewStore(render)
	imagesConfiguration := imagesconfiguration.NewImagesConfiguration(e.fs, graph, store, render, e.compatibility)
	err := imagesConfiguration.LoadImagesToStore(conf.ImagesPath)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return store, nil
}

func (e *GetBuildersEntrypoint) createImagesGraphTemplatesStorer(factory *graph.GraphTemplateFactory) (*imagesgraphtemplate.ImagesGraphTemplate, error) {
	errContext := "(entrypoint::get::builders::createImagesGraphTemplatesStorer)"

	if factory == nil {
		return nil, errors.New(errContext, "To create an images graph templates storer in get builders entrypoint, a graph template factory is required")
	}

	graphtemplate := imagesgraphtemplate.NewImagesGraphTemplate(factory)

	return graphtemplate, nil
}

func (e *GetBuildersEntrypoint) createGraphTemplateFactory() (*graph.GraphTemplateFactory, error) {
	return graph.NewGraphTemplateFactory(false), nil
}

func (e *GetBuildersEntrypoint) createBuildersStore(conf *configuration.Configuration) (*buildersstore.Store, error) {

	errContext := "(entrypoint::get::builders::createBuildersStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create a builders store in build entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create a builders store in build entrypoint, configuration is required")
	}

	if conf.BuildersPath == "" {
		return nil, errors.New(errContext, "To create a builders store in build entrypoint, builders path must be provided in configuration")
	}

	buildersStore := buildersstore.NewStore()
	buildersConfiguration := buildersconfiguration.NewBuilders(e.fs, buildersStore)
	err := buildersConfiguration.LoadBuilders(conf.BuildersPath)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return buildersStore, nil
}

func (e *GetBuildersEntrypoint) createOutput() (repository.BuildersOutputter, error) {

	errContext := "(entrypoint::get::builders::createOutput)"

	if e.writer == nil {
		return nil, errors.New(errContext, "Get images entrypoint requires a writer to create the plain text output")
	}
	c := console.NewConsole(e.writer, nil)
	o := output.NewPlainOutput(
		output.WithWriter(c),
	)

	return o, nil
}
