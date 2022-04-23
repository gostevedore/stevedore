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
	"github.com/gostevedore/stevedore/internal/promote"
	repodocker "github.com/gostevedore/stevedore/internal/promote/docker"
	repodockercopy "github.com/gostevedore/stevedore/internal/promote/docker/promoter"
	repodryrun "github.com/gostevedore/stevedore/internal/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/semver"
	service "github.com/gostevedore/stevedore/internal/service/promote"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the build application
type Entrypoint struct {
	writer io.Writer
	fs     afero.Fs
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

// WitFileSystem creates a new file system
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
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, options *handler.Options) error {
	var err error
	var promoteRepoFactory promote.PromoteFactory
	var credentialsStore *credentials.CredentialsStore
	var semverGenerator *semver.SemVerGenerator

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

	options.EnableSemanticVersionTags = conf.EnableSemanticVersionTags || options.EnableSemanticVersionTags
	if options.EnableSemanticVersionTags && len(conf.SemanticVersionTagsTemplates) > 0 && len(options.SemanticVersionTagsTemplates) == 0 {
		options.SemanticVersionTagsTemplates = append([]string{}, conf.SemanticVersionTagsTemplates...)
	}

	promoteRepoFactory, err = e.createPromoteFactory()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if conf.DockerCredentialsDir == "" {
		return errors.New(errContext, "Docker credentials path must be provided in the configuration")
	}

	credentialsStore, err = e.createCredentialsStore(conf.DockerCredentialsDir)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	semverGenerator, err = e.createSemanticVersionFactory()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	promoteService := service.NewService(
		service.WithPromoteFactory(promoteRepoFactory),
		service.WithCredentials(credentialsStore),
		service.WithSemver(semverGenerator),
	)

	promoteHandler := handler.NewHandler(promoteService)
	err = promoteHandler.Handler(ctx, options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

func (e *Entrypoint) createCredentialsStore(path string) (*credentials.CredentialsStore, error) {
	errContext := "(Entrypoint::createPromoteRepoFactory)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials store, a file system is required")
	}

	credentialsStore := credentials.NewCredentialsStore(e.fs)
	err := credentialsStore.LoadCredentials(path)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return credentialsStore, nil
}

func (e *Entrypoint) createPromoteFactory() (promote.PromoteFactory, error) {

	errContext := "(Entrypoint::createPromoteFactory)"

	dockerClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	copyCmd := copy.NewDockerImageCopyCmd(dockerClient)
	copyCmdFacade := repodockercopy.NewDockerCopy(copyCmd)
	promoteRepoDocker := repodocker.NewDockerPromote(copyCmdFacade, os.Stdout)
	promoteRepoDryRun := repodryrun.NewDryRunPromote(copyCmdFacade, os.Stdout)
	promoteRepoFactory := promote.NewPromoteFactory()
	err = promoteRepoFactory.Register(promote.DockerPromoterName, promoteRepoDocker)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}
	err = promoteRepoFactory.Register(promote.DryRunPromoterName, promoteRepoDryRun)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return promoteRepoFactory, nil
}

func (e *Entrypoint) createSemanticVersionFactory() (*semver.SemVerGenerator, error) {
	return semver.NewSemVerGenerator(), nil
}
