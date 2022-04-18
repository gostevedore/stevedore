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
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/driver"
	ansibledriver "github.com/gostevedore/stevedore/internal/driver/ansible"
	"github.com/gostevedore/stevedore/internal/driver/ansible/goansible"
	defaultdriver "github.com/gostevedore/stevedore/internal/driver/default"
	dockerdriver "github.com/gostevedore/stevedore/internal/driver/docker"
	"github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder"
	dockerdrivercontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context"
	gitauth "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git/auth"
	dryrundriver "github.com/gostevedore/stevedore/internal/driver/dryrun"
	buildservice "github.com/gostevedore/stevedore/internal/engine/build"
	"github.com/gostevedore/stevedore/internal/engine/build/command"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
	build "github.com/gostevedore/stevedore/internal/handler/build"
	buildhandler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/images/image/render"
	"github.com/gostevedore/stevedore/internal/images/image/render/now"
	"github.com/gostevedore/stevedore/internal/images/store"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
	"github.com/gostevedore/stevedore/internal/schedule/job"
	"github.com/gostevedore/stevedore/internal/schedule/worker"
	"github.com/gostevedore/stevedore/internal/semver"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the build application
type Entrypoint struct {
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

// Options provides the options for the entrypoint
func (e *Entrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute executes the entrypoint
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, entrypointOptions *EntrypointOptions, handlerOptions *build.HandlerOptions) error {
	var err error
	var buildHandler *buildhandler.Handler
	var buildDriverFactory driver.BuildDriverFactory
	var dispatcher *dispatch.Dispatch
	var planFactory *plan.PlanFactory
	var buildService *buildservice.Service

	errContext := "(Entrypoint::Execute)"

	if conf == nil {
		return errors.New(errContext, "To execute the build entrypoint, configuration is required")
	}

	if len(args) < 1 {
		return errors.New(errContext, "To execute the build entrypoint, arguments are required")
	}

	if entrypointOptions == nil {
		return errors.New(errContext, "To execute the build entrypoint, entrypoint options are required")
	}

	if handlerOptions == nil {
		return errors.New(errContext, "To execute the build entrypoint, handler options are required")
	}

	imageName := args[0]
	if len(args) > 1 {
		fmt.Fprintf(e.writer, "Ignoring extra arguments: %v\n", args[1:])
	}

	if conf.Concurrency > 0 && entrypointOptions.Concurrency < 1 {
		entrypointOptions.Concurrency = conf.Concurrency
	}
	handlerOptions.PushImagesAfterBuild = conf.PushImages || handlerOptions.PushImagesAfterBuild
	handlerOptions.EnableSemanticVersionTags = conf.EnableSemanticVersionTags || handlerOptions.EnableSemanticVersionTags
	if handlerOptions.EnableSemanticVersionTags && len(conf.SemanticVersionTagsTemplates) > 0 && len(handlerOptions.SemanticVersionTagsTemplates) == 0 {
		handlerOptions.SemanticVersionTagsTemplates = append([]string{}, conf.SemanticVersionTagsTemplates...)
	}

	credentialsStore := credentials.NewCredentialsStore()
	buildersStore := buildersstore.NewBuildersStore()
	commandFactory := command.NewBuildCommandFactory()
	jobFactory := job.NewJobFactory()
	semverser := semver.NewSemVerGenerator()

	buildDriverFactory, err = e.createBuildDriverFactory(credentialsStore, entrypointOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	dispatcher, err = e.createDispatcher(entrypointOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = dispatcher.Start(ctx)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	buildService = buildservice.NewService(
		buildservice.WithBuilders(buildersStore),
		buildservice.WithCommandFactory(commandFactory),
		buildservice.WithDriverFactory(buildDriverFactory),
		buildservice.WithJobFactory(jobFactory),
		buildservice.WithDispatch(dispatcher),
		buildservice.WithSemver(semverser),
		buildservice.WithCredentials(credentialsStore),
	)

	renderImages := render.NewImageRender(now.NewNow())
	imagesStore := store.NewImageStore(renderImages)

	planFactory, err = e.createPlanFactory(imagesStore, entrypointOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	buildHandler = buildhandler.NewHandler(planFactory, buildService)
	err = buildHandler.Handler(ctx, imageName, handlerOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

func (e *Entrypoint) createBuildDriverFactory(credentialsStore *credentials.CredentialsStore, options *EntrypointOptions) (driver.BuildDriverFactory, error) {

	var ansiblePlaybookDriver driver.BuildDriverer
	var defaultDriver driver.BuildDriverer
	var dockerDriver driver.BuildDriverer
	var dryRunDriver driver.BuildDriverer
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

	factory := driver.NewBuildDriverFactory()

	ansiblePlaybookDriver, err = e.createAnsibleDriver(options)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}
	dockerDriver, err = e.createDockerDriver(credentialsStore, options)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}
	defaultDriver, err = e.createDefaultDriver()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	dryRunDriver, err = e.createDryRunDriver()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	factory.Register("ansible-playbook", ansiblePlaybookDriver)
	factory.Register("docker", dockerDriver)
	factory.Register("default", defaultDriver)
	factory.Register("dry-run", dryRunDriver)

	return factory, nil
}

func (e *Entrypoint) createDefaultDriver() (driver.BuildDriverer, error) {
	return defaultdriver.NewDefaultDriver(e.writer), nil
}

func (e *Entrypoint) createDryRunDriver() (driver.BuildDriverer, error) {
	return dryrundriver.NewDryRunDriver(e.writer), nil
}

func (e *Entrypoint) createAnsibleDriver(options *EntrypointOptions) (driver.BuildDriverer, error) {

	errContext := "(entrypoint::createAnsibleDriver)"

	if options == nil {
		return nil, errors.New(errContext, "Entrypoint options are required to create ansible driver")
	}

	ansiblePlaybookDriver, err := ansibledriver.NewAnsiblePlaybookDriver(goansible.NewGoAnsibleDriver(), e.writer)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return ansiblePlaybookDriver, nil
}

func (e *Entrypoint) createDockerDriver(credentialsStore *credentials.CredentialsStore, options *EntrypointOptions) (driver.BuildDriverer, error) {
	var dockerClient *dockerclient.Client
	var dockerDriver *dockerdriver.DockerDriver
	var dockerDriverBuldContext *dockerdrivercontext.DockerBuildContextFactory
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
		return nil, errors.New(errContext, err.Error())
	}

	goDockerBuild := godockerbuild.NewDockerBuildCmd(dockerClient)
	gitAuth = gitauth.NewGitAuthFactory(credentialsStore)
	dockerDriverBuldContext = dockerdrivercontext.NewDockerBuildContextFactory(gitAuth)
	goDockerBuildDriver = godockerbuilder.NewGoDockerBuildDriver(goDockerBuild, dockerDriverBuldContext)
	dockerDriver, err = dockerdriver.NewDockerDriver(goDockerBuildDriver, e.writer)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return dockerDriver, nil
}

func (e *Entrypoint) createDispatcher(options *EntrypointOptions) (*dispatch.Dispatch, error) {
	dispatchWorker := worker.NewWorkerFactory()
	d := dispatch.NewDispatch(dispatchWorker, dispatch.WithNumWorkers(options.Concurrency))

	return d, nil
}

func (e *Entrypoint) createPlanFactory(store *store.ImageStore, options *EntrypointOptions) (*plan.PlanFactory, error) {
	factory := plan.NewPlanFactory(store)

	return factory, nil
}
