package entrypoint

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/copy"
	dockerclient "github.com/docker/docker/client"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	repofactory "github.com/gostevedore/stevedore/internal/promote"
	repodocker "github.com/gostevedore/stevedore/internal/promote/docker"
	repodockercopy "github.com/gostevedore/stevedore/internal/promote/docker/promoter"
	repodryrun "github.com/gostevedore/stevedore/internal/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/semver"
	service "github.com/gostevedore/stevedore/internal/service/promote"
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
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, options *handler.Options) error {

	errContext := "(Entrypoint::Execute)"

	if conf == nil {
		return errors.New(errContext, "To execute the promote entrypoint, configuration is required")
	}

	if len(args) < 1 {
		return errors.New(errContext, "To execute the promote entrypoint, arguments are required")
	}

	if options == nil {
		return errors.New(errContext, "To execute the promote entrypoint, handler options are required")
	}

	options.SourceImageName = args[0]
	if len(args) > 1 {
		fmt.Fprintf(e.writer, "Ignoring extra arguments: %v\n", args[1:])
	}

	dockerClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	copyCmd := copy.NewDockerImageCopyCmd(dockerClient)
	copyCmdFacade := repodockercopy.NewDockerCopy(copyCmd)
	promoteRepoDocker := repodocker.NewDockerPromote(copyCmdFacade, os.Stdout)
	promoteRepoDryRun := repodryrun.NewDryRunPromote(copyCmdFacade, os.Stdout)
	promoteRepoFactory := repofactory.NewPromoteFactory()
	err = promoteRepoFactory.Register("docker", promoteRepoDocker)
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	err = promoteRepoFactory.Register("dry-run", promoteRepoDryRun)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	credentialsStore := credentials.NewCredentialsStore()
	err = credentialsStore.LoadCredentials(conf.DockerCredentialsDir)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	semverGenerator := semver.NewSemVerGenerator()

	promoteService := service.NewService(promoteRepoFactory, conf, credentialsStore, semverGenerator)

	promoteHandler := handler.NewHandler(promoteService)
	err = promoteHandler.Handler(ctx, options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
