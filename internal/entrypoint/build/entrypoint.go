package build

import (
	"context"
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	godockerbuild "github.com/apenella/go-docker-builder/pkg/build"
	dockerclient "github.com/docker/docker/client"
	buildersstore "github.com/gostevedore/stevedore/internal/builders/store"
	"github.com/gostevedore/stevedore/internal/configuration"
	buildersconfiguration "github.com/gostevedore/stevedore/internal/configuration/builders"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/credentials"
	driveransible "github.com/gostevedore/stevedore/internal/driver/ansible"
	"github.com/gostevedore/stevedore/internal/driver/ansible/goansible"
	driverdefault "github.com/gostevedore/stevedore/internal/driver/default"
	driverdocker "github.com/gostevedore/stevedore/internal/driver/docker"
	"github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder"
	driverdockercontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context"
	gitauth "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git/auth"
	driverdryrun "github.com/gostevedore/stevedore/internal/driver/dryrun"
	driverfactory "github.com/gostevedore/stevedore/internal/driver/factory"
	buildhandler "github.com/gostevedore/stevedore/internal/handler/build"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/images/graph"
	"github.com/gostevedore/stevedore/internal/images/render"
	"github.com/gostevedore/stevedore/internal/images/render/now"
	"github.com/gostevedore/stevedore/internal/images/store"
	imagesstore "github.com/gostevedore/stevedore/internal/images/store"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
	"github.com/gostevedore/stevedore/internal/schedule/job"
	"github.com/gostevedore/stevedore/internal/schedule/worker"
	"github.com/gostevedore/stevedore/internal/semver"
	buildservice "github.com/gostevedore/stevedore/internal/service/build"
	"github.com/gostevedore/stevedore/internal/service/build/command"
	"github.com/gostevedore/stevedore/internal/service/build/plan"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the build application
type Entrypoint struct {
	fs     afero.Fs
	writer io.Writer
}

// NewEntrypoint returns a new entrypoint
func NewEntrypoint(opts ...OptionsFunc) *Entrypoint {
	e := &Entrypoint{}
	e.Options(opts...)

	return e
}

// WithWriter sets the writer for the entrypoint
func WithWriter(w io.Writer) OptionsFunc {
	return func(e *Entrypoint) {
		e.writer = w
	}
}

// WithFileSystem sets the writer for the entrypoint
func WithFileSystem(fs afero.Fs) OptionsFunc {
	return func(e *Entrypoint) {
		e.fs = fs
	}
}

// Options provides the options for the entrypoint
func (e *Entrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute executes the entrypoint
func (e *Entrypoint) Execute(
	ctx context.Context,
	args []string,
	conf *configuration.Configuration,
	compatibility Compatibilitier,
	inputEntrypointOptions *Options,
	inputHandlerOptions *handler.Options) error {

	var buildDriverFactory driverfactory.BuildDriverFactory
	var buildersStore *buildersstore.BuildersStore
	var buildHandler *buildhandler.Handler
	var buildService *buildservice.Service
	var commandFactory *command.BuildCommandFactory
	var credentialsStore *credentials.CredentialsStore
	var dispatcher *dispatch.Dispatch
	var entrypointOptions *Options
	var err error
	var handlerOptions *handler.Options
	var imageName string
	var imageRender *render.ImageRender
	var imagesGraphTemplatesStore *imagesgraphtemplate.ImagesGraphTemplate
	var imagesStore *imagesstore.ImageStore
	var jobFactory *job.JobFactory
	var planFactory *plan.PlanFactory
	var semVerFactory *semver.SemVerGenerator
	var graphTemplateFactory *graph.GraphTemplateFactory

	errContext := "(Entrypoint::Execute)"

	imageName, err = e.prepareImageName(args)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	entrypointOptions, err = e.prepareEntrypointOptions(conf, inputEntrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	handlerOptions, err = e.prepareHandlerOptions(conf, inputHandlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	credentialsStore, err = e.createCredentialsStore(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildersStore, err = e.createBuildersStore(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	commandFactory, err = e.createCommandFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	jobFactory, err = e.createJobFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	semVerFactory, err = e.createSemVerFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildDriverFactory, err = e.createBuildDriverFactory(credentialsStore, entrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	dispatcher, err = e.createDispatcher(entrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = dispatcher.Start(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildService = buildservice.NewService(
		buildservice.WithBuilders(buildersStore),
		buildservice.WithCommandFactory(commandFactory),
		buildservice.WithDriverFactory(buildDriverFactory),
		buildservice.WithJobFactory(jobFactory),
		buildservice.WithDispatch(dispatcher),
		buildservice.WithSemver(semVerFactory),
		buildservice.WithCredentials(credentialsStore),
	)

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

	imagesStore, err = e.createImagesStore(conf, imageRender, imagesGraphTemplatesStore, compatibility)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	planFactory, err = e.createPlanFactory(imagesStore, entrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	buildHandler = buildhandler.NewHandler(planFactory, buildService)
	err = buildHandler.Handler(ctx, imageName, handlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (e *Entrypoint) prepareEntrypointOptions(conf *configuration.Configuration, inputEntrypointOptions *Options) (*Options, error) {

	errContext := "(Entrypoint::prepareEntrypointOptions)"

	options := &Options{}

	if conf == nil {
		return nil, errors.New(errContext, "To prepare entrypoint options, configuration is required")
	}

	if inputEntrypointOptions == nil {
		return nil, errors.New(errContext, "To prepare entrypoint options, entrypoint options are required")
	}

	options.Concurrency = inputEntrypointOptions.Concurrency
	if conf.Concurrency > 0 && options.Concurrency < 1 {
		options.Concurrency = conf.Concurrency
	}
	options.Debug = inputEntrypointOptions.Debug

	return options, nil
}

func (e *Entrypoint) prepareImageName(args []string) (string, error) {

	errContext := "(Entrypoint::prepareImageName)"

	if len(args) < 1 || args == nil {
		return "", errors.New(errContext, "To execute the build entrypoint, arguments are required")
	}

	imageName := args[0]
	if len(args) > 1 {
		fmt.Fprintf(e.writer, "Ignoring extra arguments: %v\n", args[1:])
	}

	return imageName, nil
}

func (e *Entrypoint) prepareHandlerOptions(conf *configuration.Configuration, inputHandlerOptions *handler.Options) (*handler.Options, error) {

	errContext := "(Entrypoint::prepareHandlerOptions)"
	options := &handler.Options{}

	if conf == nil {
		return nil, errors.New(errContext, "To prepare handler options, configuration is required")
	}

	if inputHandlerOptions == nil {
		return nil, errors.New(errContext, "To prepare handler options, handler options are required")
	}

	options.AnsibleConnectionLocal = inputHandlerOptions.AnsibleConnectionLocal
	options.AnsibleIntermediateContainerName = inputHandlerOptions.AnsibleIntermediateContainerName
	options.AnsibleInventoryPath = inputHandlerOptions.AnsibleInventoryPath
	options.AnsibleLimit = inputHandlerOptions.AnsibleLimit
	options.BuildOnCascade = inputHandlerOptions.BuildOnCascade
	options.CascadeDepth = inputHandlerOptions.CascadeDepth
	options.DryRun = inputHandlerOptions.DryRun
	options.EnableSemanticVersionTags = conf.EnableSemanticVersionTags || inputHandlerOptions.EnableSemanticVersionTags
	options.ImageFromName = inputHandlerOptions.ImageFromName
	options.ImageFromRegistryHost = inputHandlerOptions.ImageFromRegistryHost
	options.ImageFromRegistryNamespace = inputHandlerOptions.ImageFromRegistryNamespace
	options.ImageFromVersion = inputHandlerOptions.ImageFromVersion
	options.ImageName = inputHandlerOptions.ImageName
	options.ImageRegistryHost = inputHandlerOptions.ImageRegistryHost
	options.ImageRegistryNamespace = inputHandlerOptions.ImageRegistryNamespace
	options.Labels = append([]string{}, inputHandlerOptions.Labels...)
	options.PersistentVars = append([]string{}, inputHandlerOptions.PersistentVars...)
	options.PullParentImage = inputHandlerOptions.PullParentImage
	options.PushImagesAfterBuild = conf.PushImages || inputHandlerOptions.PushImagesAfterBuild
	options.RemoveImagesAfterPush = inputHandlerOptions.RemoveImagesAfterPush

	options.SemanticVersionTagsTemplates = append([]string{}, inputHandlerOptions.SemanticVersionTagsTemplates...)
	if inputHandlerOptions.EnableSemanticVersionTags && len(conf.SemanticVersionTagsTemplates) > 0 && len(options.SemanticVersionTagsTemplates) == 0 {
		options.SemanticVersionTagsTemplates = append([]string{}, conf.SemanticVersionTagsTemplates...)
	}
	options.Tags = append([]string{}, inputHandlerOptions.Tags...)
	options.Vars = append([]string{}, inputHandlerOptions.Vars...)
	options.Versions = inputHandlerOptions.Versions

	return options, nil
}

func (e *Entrypoint) createCredentialsStore(conf *configuration.Configuration) (*credentials.CredentialsStore, error) {
	errContext := "(Entrypoint::createCredentialsStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials store, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials store, configuration is required")
	}

	if conf.DockerCredentialsDir == "" {
		return nil, errors.New(errContext, "To create the credentials store, credentials path must be provided in configuration")
	}

	credentialsStore := credentials.NewCredentialsStore(e.fs)
	err := credentialsStore.LoadCredentials(conf.DockerCredentialsDir)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credentialsStore, nil
}

func (e *Entrypoint) createBuildersStore(conf *configuration.Configuration) (*buildersstore.BuildersStore, error) {

	errContext := "(Entrypoint::createBuildersStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create a builders store, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create a builders store, configuration is required")
	}

	if conf.BuildersPath == "" {
		return nil, errors.New(errContext, "To create a builders store, builders path must be provided in configuration")
	}

	buildersStore := buildersstore.NewBuildersStore()
	buildersConfiguration := buildersconfiguration.NewBuilders(e.fs, buildersStore)
	err := buildersConfiguration.LoadBuilders(conf.BuildersPath)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return buildersStore, nil
}

func (e *Entrypoint) createCommandFactory() (*command.BuildCommandFactory, error) {
	return command.NewBuildCommandFactory(), nil
}

func (e *Entrypoint) createJobFactory() (*job.JobFactory, error) {
	return job.NewJobFactory(), nil
}

func (e *Entrypoint) createSemVerFactory() (*semver.SemVerGenerator, error) {
	return semver.NewSemVerGenerator(), nil
}

func (e *Entrypoint) createImageRender(now render.Nower) (*render.ImageRender, error) {
	errContext := "(Entrypoint::createImageRender)"

	if now == nil {
		return nil, errors.New(errContext, "To create an image render, a nower is required")
	}

	return render.NewImageRender(now), nil
}

func (e *Entrypoint) createImagesStore(conf *configuration.Configuration, render imagesstore.ImageRenderer, graph imagesconfiguration.ImagesGraphTemplatesStorer, compatibility Compatibilitier) (*imagesstore.ImageStore, error) {

	errContext := "(Entrypoint::createImagesStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create an images store, a filesystem is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create an images store, configuration is required")
	}

	if render == nil {
		return nil, errors.New(errContext, "To create an images store, image render is required")
	}

	if graph == nil {
		return nil, errors.New(errContext, "To create an images store, images graph templates storer is required")
	}

	if compatibility == nil {
		return nil, errors.New(errContext, "To create an images store, compatibility is required")
	}

	if conf.ImagesPath == "" {
		return nil, errors.New(errContext, "To create an images store, images path must be provided in configuration")
	}

	store := imagesstore.NewImageStore(render)
	imagesConfiguration := imagesconfiguration.NewImagesConfiguration(e.fs, graph, store, compatibility)
	err := imagesConfiguration.LoadImagesToStore(conf.ImagesPath)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return store, nil
}

func (e *Entrypoint) createImagesGraphTemplatesStorer(factory *graph.GraphTemplateFactory) (*imagesgraphtemplate.ImagesGraphTemplate, error) {
	errContext := "(Entrypoint::createImagesGraphTemplatesStorer)"

	if factory == nil {
		return nil, errors.New(errContext, "To create an images graph templates storer, a graph template factory is required")
	}

	graphtemplate := imagesgraphtemplate.NewImagesGraphTemplate(factory)

	return graphtemplate, nil
}

func (e *Entrypoint) createGraphTemplateFactory() (*graph.GraphTemplateFactory, error) {
	return graph.NewGraphTemplateFactory(false), nil
}

func (e *Entrypoint) createBuildDriverFactory(credentialsStore *credentials.CredentialsStore, options *Options) (driverfactory.BuildDriverFactory, error) {

	var ansiblePlaybookDriver repository.BuildDriverer
	var defaultDriver repository.BuildDriverer
	var dockerDriver repository.BuildDriverer
	var dryRunDriver repository.BuildDriverer
	var err error

	errContext := "(entrypoint::createBuildDriverFactory)"

	if credentialsStore == nil {
		return nil, errors.New(errContext, "Register drivers requires a credentials store")
	}

	if options == nil {
		return nil, errors.New(errContext, "Register drivers requires options")
	}

	if e.writer == nil {
		return nil, errors.New(errContext, "Register drivers requires a writer")
	}

	factory := driverfactory.NewBuildDriverFactory()

	ansiblePlaybookDriver, err = e.createAnsibleDriver(options)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	dockerDriver, err = e.createDockerDriver(credentialsStore, options)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	defaultDriver, err = e.createDefaultDriver()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	dryRunDriver, err = e.createDryRunDriver()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	factory.Register(image.AnsiblePlaybookDriverName, ansiblePlaybookDriver)
	factory.Register(image.DockerDriverName, dockerDriver)
	factory.Register(image.DefaultDriverName, defaultDriver)
	factory.Register(image.DryRunDriverName, dryRunDriver)

	return factory, nil
}

func (e *Entrypoint) createDefaultDriver() (repository.BuildDriverer, error) {
	return driverdefault.NewDefaultDriver(e.writer), nil
}

func (e *Entrypoint) createDryRunDriver() (repository.BuildDriverer, error) {
	return driverdryrun.NewDryRunDriver(e.writer), nil
}

func (e *Entrypoint) createAnsibleDriver(options *Options) (repository.BuildDriverer, error) {

	errContext := "(entrypoint::createAnsibleDriver)"

	if options == nil {
		return nil, errors.New(errContext, "Entrypoint options are required to create ansible driver")
	}

	ansiblePlaybookDriver, err := driveransible.NewAnsiblePlaybookDriver(goansible.NewGoAnsibleDriver(), e.writer)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return ansiblePlaybookDriver, nil
}

func (e *Entrypoint) createDockerDriver(credentialsStore *credentials.CredentialsStore, options *Options) (repository.BuildDriverer, error) {
	var dockerClient *dockerclient.Client
	var dockerDriver *driverdocker.DockerDriver
	var dockerDriverBuldContext *driverdockercontext.DockerBuildContextFactory
	var err error
	var gitAuth *gitauth.GitAuthFactory
	var goDockerBuildDriver *godockerbuilder.GoDockerBuildDriver

	errContext := "(entrypoint::createDockerDriver)"

	if credentialsStore == nil {
		return nil, errors.New(errContext, "Docker driver requires a credentials store")
	}

	if options == nil {
		return nil, errors.New(errContext, "Entrypoint options are required to create docker driver")
	}

	dockerClient, err = dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	goDockerBuild := godockerbuild.NewDockerBuildCmd(dockerClient)
	gitAuth = gitauth.NewGitAuthFactory(credentialsStore)
	dockerDriverBuldContext = driverdockercontext.NewDockerBuildContextFactory(gitAuth)
	goDockerBuildDriver = godockerbuilder.NewGoDockerBuildDriver(goDockerBuild, dockerDriverBuldContext)
	dockerDriver, err = driverdocker.NewDockerDriver(goDockerBuildDriver, e.writer)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return dockerDriver, nil
}

func (e *Entrypoint) createDispatcher(options *Options) (*dispatch.Dispatch, error) {
	dispatchWorker := worker.NewWorkerFactory()
	d := dispatch.NewDispatch(dispatchWorker, dispatch.WithNumWorkers(options.Concurrency))

	return d, nil
}

func (e *Entrypoint) createPlanFactory(store *store.ImageStore, options *Options) (*plan.PlanFactory, error) {
	factory := plan.NewPlanFactory(store)

	return factory, nil
}
