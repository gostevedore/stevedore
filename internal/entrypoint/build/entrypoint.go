package build

import (
	"context"
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	godockerbuild "github.com/apenella/go-docker-builder/pkg/build"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	dockerclient "github.com/docker/docker/client"
	application "github.com/gostevedore/stevedore/internal/application/build"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	buildhandler "github.com/gostevedore/stevedore/internal/handler/build"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	buildersconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/builders"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/credentials/compatibility"
	credentialsfactory "github.com/gostevedore/stevedore/internal/infrastructure/credentials/factory"
	credentialsformatfactory "github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/factory"
	authmethodbasic "github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/basic"
	authmethodkeyfile "github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/keyfile"
	authmethodsshagent "github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/sshagent"
	authproviderawsecr "github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/awsecr"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/awsecr/token"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/awsecr/token/awscredprovider"
	authproviderbadge "github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/badge"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/ansible"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/ansible/goansible"
	driverdefault "github.com/gostevedore/stevedore/internal/infrastructure/driver/default"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/docker/godockerbuilder"
	dockercontext "github.com/gostevedore/stevedore/internal/infrastructure/driver/docker/godockerbuilder/context"
	gitauth "github.com/gostevedore/stevedore/internal/infrastructure/driver/docker/godockerbuilder/context/git/auth"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/dispatch"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/worker"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	credentialsstorefactory "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/factory"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the build application
type Entrypoint struct {
	fs            afero.Fs
	writer        io.Writer
	compatibility Compatibilitier
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

// WithCompatibility sets the compatibility for the entrypoint
func WithCompatibility(c Compatibilitier) OptionsFunc {
	return func(e *Entrypoint) {
		e.compatibility = c
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
	inputEntrypointOptions *Options,
	inputHandlerOptions *handler.Options) error {

	var buildDriverFactory factory.BuildDriverFactory
	var buildersStore *builders.Store
	var buildHandler *buildhandler.Handler
	var buildService *application.Application
	var commandFactory *command.BuildCommandFactory
	var credentialsFactory repository.CredentialsFactorier
	var dispatcher *dispatch.Dispatch
	var entrypointOptions *Options
	var err error
	var handlerOptions *handler.Options
	var imageName string
	var imageRender *render.ImageRender
	var imagesGraphTemplatesStore *imagesgraphtemplate.ImagesGraphTemplate
	var imagesStore *images.Store
	var jobFactory *job.JobFactory
	var planFactory *plan.PlanFactory
	var semVerFactory *semver.SemVerGenerator
	var graphTemplateFactory *graph.GraphTemplateFactory

	errContext := "(build::entrypoint::Execute)"

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

	credentialsFactory, err = e.createCredentialsFactory(conf)
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

	buildDriverFactory, err = e.createBuildDriverFactory(credentialsFactory, entrypointOptions)
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

	buildService = application.NewApplication(
		application.WithBuilders(buildersStore),
		application.WithCommandFactory(commandFactory),
		application.WithDriverFactory(buildDriverFactory),
		application.WithJobFactory(jobFactory),
		application.WithDispatch(dispatcher),
		application.WithSemver(semVerFactory),
		application.WithCredentials(credentialsFactory),
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

	imagesStore, err = e.createImagesStore(conf, imageRender, imagesGraphTemplatesStore)
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

	errContext := "(build::entrypoint::prepareEntrypointOptions)"

	options := &Options{}

	if conf == nil {
		return nil, errors.New(errContext, "To prepare build entrypoint options, configuration is required")
	}

	if inputEntrypointOptions == nil {
		return nil, errors.New(errContext, "To prepare build entrypoint options, entrypoint options are required")
	}

	options.Concurrency = inputEntrypointOptions.Concurrency
	if conf.Concurrency > 0 && options.Concurrency < 1 {
		options.Concurrency = conf.Concurrency
	}
	options.Debug = inputEntrypointOptions.Debug

	return options, nil
}

func (e *Entrypoint) prepareImageName(args []string) (string, error) {

	errContext := "(build::entrypoint::prepareImageName)"

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

	errContext := "(build::entrypoint::prepareHandlerOptions)"
	options := &handler.Options{}

	if conf == nil {
		return nil, errors.New(errContext, "To prepare handler options in build entrypoint, configuration is required")
	}

	if inputHandlerOptions == nil {
		return nil, errors.New(errContext, "To prepare handler options in build entrypoint, handler options are required")
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
	options.PersistentLabels = append([]string{}, inputHandlerOptions.PersistentLabels...)
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

func (e *Entrypoint) createCredentialsLocalStore(conf *configuration.CredentialsConfiguration) (*credentialslocalstore.LocalStore, error) {

	errContext := "(build::entrypoint::createCredentialsStore)"

	if conf == nil {
		return nil, errors.New(errContext, "To create credentials store in build entrypoint, credentials configuration is required")
	}

	if conf.Format == "" {
		return nil, errors.New(errContext, "To create credentials store in build entrypoint, credentials format must be specified")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create credentials store in build entrypoint, compatibility is required")
	}

	switch conf.StorageType {
	case credentials.LocalStore:
		if conf.LocalStoragePath == "" {
			return nil, errors.New(errContext, "To create credentials store in build entrypoint, local storage path is required")
		}

		credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)

		credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
		credentialsFormat, err := credentialsFormatFactory.Get(conf.Format)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		store := credentialslocalstore.NewLocalStore(e.fs, conf.LocalStoragePath, credentialsFormat, credentialsCompatibility)

		return store, nil
	default:
		return nil, errors.New(errContext, fmt.Sprintf("Unsupported credentials storage type '%s'", conf.StorageType))
	}
}

func (e *Entrypoint) createCredentialsFactory(conf *configuration.Configuration) (repository.CredentialsFactorier, error) {
	errContext := "(build::entrypoint::createCredentialsFactory)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials store in build entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials store in build entrypoint, configuration is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials store in build entrypoint, credentials configuration is required")
	}

	// create credentials store
	localstore, err := e.createCredentialsLocalStore(conf.Credentials)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	storefactory := credentialsstorefactory.NewCredentialsStoreFactory()
	storefactory.Register(credentials.LocalStore, localstore)
	// since there is only one store, we can use it directly
	store, err := storefactory.Get(credentials.LocalStore)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	// create authorization methods
	basic := authmethodbasic.NewBasicAuthMethod()
	keyfile := authmethodkeyfile.NewKeyFileAuthMethod()
	sshagent := authmethodsshagent.NewSSHAgentAuthMethod()

	// create authorization badge provider
	badge := authproviderbadge.NewBadgeCredentialsProvider(basic, keyfile, sshagent)

	// create authorization aws ecr provider
	tokenProvider := token.NewAWSECRToken(
		token.WithStaticCredentialsProvider(awscredprovider.NewStaticCredentialsProvider()),
		token.WithAssumeRoleARNProvider(awscredprovider.NewAssumerRoleARNProvider()),
		token.WithECRClientFactory(
			token.NewECRClientFactory(
				func(cfg aws.Config) token.ECRClienter {
					c := ecr.NewFromConfig(cfg)

					return c
				},
			),
		),
	)

	awsecr := authproviderawsecr.NewAWSECRCredentialsProvider(tokenProvider)

	// create credentials factory
	factory := credentialsfactory.NewCredentialsFactory(store, badge, awsecr)

	return factory, nil
}

func (e *Entrypoint) createBuildersStore(conf *configuration.Configuration) (*builders.Store, error) {

	errContext := "(build::entrypoint::createBuildersStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create a builders store in build entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create a builders store in build entrypoint, configuration is required")
	}

	if conf.BuildersPath == "" {
		return nil, errors.New(errContext, "To create a builders store in build entrypoint, builders path must be provided in configuration")
	}

	buildersStore := builders.NewStore()
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
	errContext := "(build::entrypoint::createImageRender)"

	if now == nil {
		return nil, errors.New(errContext, "To create an image render in build entrypoint, a nower is required")
	}

	return render.NewImageRender(now), nil
}

func (e *Entrypoint) createImagesStore(conf *configuration.Configuration, render repository.Renderer, graph imagesconfiguration.ImagesGraphTemplatesStorer) (*images.Store, error) {

	errContext := "(build::entrypoint::createImagesStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create an images store in build entrypoint, a filesystem is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create an images store in build entrypoint, configuration is required")
	}

	if render == nil {
		return nil, errors.New(errContext, "To create an images store in build entrypoint, image render is required")
	}

	if graph == nil {
		return nil, errors.New(errContext, "To create an images store in build entrypoint, images graph templates storer is required")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create an images store in build entrypoint, compatibility is required")
	}

	if conf.ImagesPath == "" {
		return nil, errors.New(errContext, "To create an images store in build entrypoint, images path must be provided in configuration")
	}

	store := images.NewStore(render)
	imagesConfiguration := imagesconfiguration.NewImagesConfiguration(e.fs, graph, store, render, e.compatibility)
	err := imagesConfiguration.LoadImagesToStore(conf.ImagesPath)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return store, nil
}

func (e *Entrypoint) createImagesGraphTemplatesStorer(factory *graph.GraphTemplateFactory) (*imagesgraphtemplate.ImagesGraphTemplate, error) {
	errContext := "(build::entrypoint::createImagesGraphTemplatesStorer)"

	if factory == nil {
		return nil, errors.New(errContext, "To create an images graph templates storer in build entrypoint, a graph template factory is required")
	}

	graphtemplate := imagesgraphtemplate.NewImagesGraphTemplate(factory)

	return graphtemplate, nil
}

func (e *Entrypoint) createGraphTemplateFactory() (*graph.GraphTemplateFactory, error) {
	return graph.NewGraphTemplateFactory(false), nil
}

func (e *Entrypoint) createBuildDriverFactory(credentialsFactory repository.CredentialsFactorier, options *Options) (factory.BuildDriverFactory, error) {

	var ansiblePlaybookDriver factory.BuildDriverFactoryFunc
	var defaultDriver factory.BuildDriverFactoryFunc
	var dockerDriver factory.BuildDriverFactoryFunc
	var dryRunDriver factory.BuildDriverFactoryFunc
	var err error

	errContext := "(entrypoint::createBuildDriverFactory)"

	if credentialsFactory == nil {
		return nil, errors.New(errContext, "Register drivers requires a credentials store in build entrypoint")
	}

	if options == nil {
		return nil, errors.New(errContext, "Register drivers requires options in build entrypoint")
	}

	if e.writer == nil {
		return nil, errors.New(errContext, "Register drivers requires a writer in build entrypoint")
	}

	factory := factory.NewBuildDriverFactory()

	ansiblePlaybookDriver, err = e.createAnsibleDriver(options)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	dockerDriver, err = e.createDockerDriver(credentialsFactory, options)
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

func (e *Entrypoint) createDefaultDriver() (factory.BuildDriverFactoryFunc, error) {

	f := func() (repository.BuildDriverer, error) {
		return driverdefault.NewDefaultDriver(e.writer), nil
	}

	return f, nil
}

func (e *Entrypoint) createDryRunDriver() (factory.BuildDriverFactoryFunc, error) {

	f := func() (repository.BuildDriverer, error) {
		return dryrun.NewDryRunDriver(e.writer), nil
	}

	return f, nil
}

func (e *Entrypoint) createAnsibleDriver(options *Options) (factory.BuildDriverFactoryFunc, error) {

	errContext := "(entrypoint::createAnsibleDriver)"

	if options == nil {
		return nil, errors.New(errContext, "Build entrypoint options are required to create ansible driver")
	}

	f := func() (repository.BuildDriverer, error) {

		errContext := "(entrypoint::build::createAnsibleDriver::BuildDriverFactoryFunc)"

		ansiblePlaybookDriver, err := ansible.NewAnsiblePlaybookDriver(goansible.NewGoAnsibleDriver(), e.writer)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		return ansiblePlaybookDriver, nil
	}

	return f, nil
}

func (e *Entrypoint) createDockerDriver(credentialsFactory repository.CredentialsFactorier, options *Options) (factory.BuildDriverFactoryFunc, error) {
	var dockerClient *dockerclient.Client
	var dockerDriver *docker.DockerDriver
	var dockerDriverBuldContext *dockercontext.DockerBuildContextFactory
	var err error
	var gitAuth *gitauth.GitAuthFactory
	var goDockerBuildDriver *godockerbuilder.GoDockerBuildDriver

	errContext := "(entrypoint::createDockerDriver)"

	if credentialsFactory == nil {
		return nil, errors.New(errContext, "Docker driver requires a credentials store in build entrypoint")
	}

	if options == nil {
		return nil, errors.New(errContext, "Build entrypoint options are required to create docker driver")
	}

	f := func() (repository.BuildDriverer, error) {

		errContext := "(entrypoint::build::createDockerDriver::BuildDriverFactoryFunc)"

		dockerClient, err = dockerclient.NewClientWithOpts(dockerclient.FromEnv)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		goDockerBuild := godockerbuild.NewDockerBuildCmd(dockerClient)
		gitAuth = gitauth.NewGitAuthFactory(credentialsFactory)
		dockerDriverBuldContext = dockercontext.NewDockerBuildContextFactory(gitAuth)
		goDockerBuildDriver = godockerbuilder.NewGoDockerBuildDriver(goDockerBuild, dockerDriverBuldContext)
		dockerDriver, err = docker.NewDockerDriver(goDockerBuildDriver, e.writer)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		return dockerDriver, nil
	}

	return f, nil
}

func (e *Entrypoint) createDispatcher(options *Options) (*dispatch.Dispatch, error) {
	dispatchWorker := worker.NewWorkerFactory()
	d := dispatch.NewDispatch(dispatchWorker, dispatch.WithNumWorkers(options.Concurrency))

	return d, nil
}

func (e *Entrypoint) createPlanFactory(store *images.Store, options *Options) (*plan.PlanFactory, error) {
	factory := plan.NewPlanFactory(store)

	return factory, nil
}
