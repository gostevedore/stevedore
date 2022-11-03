package promote

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/copy"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	dockerclient "github.com/docker/docker/client"
	application "github.com/gostevedore/stevedore/internal/application/promote"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
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
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker/godockerbuilder"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	credentialsstorefactory "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/factory"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
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
func WithCompatibility(compatibility Compatibilitier) OptionsFunc {
	return func(e *Entrypoint) {
		e.compatibility = compatibility
	}
}

// Options provides the options for the entrypoint
func (e *Entrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute is a pseudo-main method for the command
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, handlerOptions *handler.Options) error {
	var err error
	var promoteRepoFactory factory.PromoteFactory
	var credentialsFactory repository.CredentialsFactorier
	var semverGenerator *semver.SemVerGenerator
	var options *handler.Options

	errContext := "(promote::entrypoint::Execute)"

	options, err = e.prepareHandlerOptions(args, conf, handlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	promoteRepoFactory, err = e.createPromoteFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	credentialsFactory, err = e.createCredentialsFactory(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	semverGenerator, err = e.createSemanticVersionFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	promoteService := application.NewApplication(
		application.WithPromoteFactory(promoteRepoFactory),
		application.WithCredentials(credentialsFactory),
		application.WithSemver(semverGenerator),
	)

	promoteHandler := handler.NewHandler(promoteService)
	err = promoteHandler.Handler(ctx, options)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (e *Entrypoint) prepareHandlerOptions(args []string, conf *configuration.Configuration, inputOptions *handler.Options) (*handler.Options, error) {
	errContext := "(promote::entrypoint::prepareHandlerOptions)"

	if len(args) < 1 || args == nil {
		return nil, errors.New(errContext, "To execute the promote entrypoint, promote image argument is required")
	}

	if inputOptions == nil {
		return nil, errors.New(errContext, "To execute the promote entrypoint, handler options are required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To execute the promote entrypoint, configuration is required")
	}

	if len(args) > 1 {
		fmt.Fprintf(e.writer, "Ignoring extra arguments: %v\n", args[1:])
	}

	options := &handler.Options{}

	options.DryRun = inputOptions.DryRun
	options.EnableSemanticVersionTags = conf.EnableSemanticVersionTags || inputOptions.EnableSemanticVersionTags
	options.TargetImageName = inputOptions.TargetImageName
	options.TargetImageRegistryNamespace = inputOptions.TargetImageRegistryNamespace
	options.TargetImageRegistryHost = inputOptions.TargetImageRegistryHost
	options.TargetImageTags = append([]string{}, inputOptions.TargetImageTags...)
	options.RemoveTargetImageTags = inputOptions.RemoveTargetImageTags
	options.SemanticVersionTagsTemplates = append([]string{}, inputOptions.SemanticVersionTagsTemplates...)
	if inputOptions.EnableSemanticVersionTags && len(conf.SemanticVersionTagsTemplates) > 0 && len(inputOptions.SemanticVersionTagsTemplates) == 0 {
		options.SemanticVersionTagsTemplates = append([]string{}, conf.SemanticVersionTagsTemplates...)
	}
	options.SourceImageName = args[0]
	options.PromoteSourceImageTag = inputOptions.PromoteSourceImageTag
	options.RemoteSourceImage = inputOptions.RemoteSourceImage

	return options, nil
}

func (e *Entrypoint) createCredentialsLocalStore(conf *configuration.CredentialsConfiguration) (*credentialslocalstore.LocalStore, error) {

	errContext := "(build::entrypoint::createCredentialsStore)"

	if conf == nil {
		return nil, errors.New(errContext, "To create credentials store in promote entrypoint, credentials configuration is required")
	}

	if conf.Format == "" {
		return nil, errors.New(errContext, "To create credentials store in promote entrypoint, credentials format must be specified")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create credentials store in promote entrypoint, compatibility is required")
	}

	switch conf.StorageType {
	case credentials.LocalStore:
		if conf.LocalStoragePath == "" {
			return nil, errors.New(errContext, "To create credentials store in promote entrypoint, local storage path is required")
		}

		credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)

		credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
		credentialsFormat, err := credentialsFormatFactory.Get(credentials.JSONFormat)
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
	errContext := "(promote::entrypoint::createCredentialsFactory)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials store in promote entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials store in promote entrypoint, configuration is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials store in promote entrypoint, credentials configuration is required")
	}

	// create credentials store
	localstore, err := e.createCredentialsLocalStore(conf.Credentials)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	storefactory := credentialsstorefactory.NewCredentialsStoreFactory()
	storefactory.Register(credentials.LocalStore, localstore)
	// since there is only one store, we can just use it directly
	store, err := storefactory.Get(credentials.LocalStore)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	// create auth methods
	basic := authmethodbasic.NewBasicAuthMethod()
	keyfile := authmethodkeyfile.NewKeyFileAuthMethod()
	sshagent := authmethodsshagent.NewSSHAgentAuthMethod()

	// create auth providers
	badge := authproviderbadge.NewBadgeCredentialsProvider(basic, keyfile, sshagent)

	// create authorization aws ecr provider
	tokenProvider := token.NewAWSECRToken(
		token.WithAssumeRoleARNProvider(awscredprovider.NewAssumerRoleARNProvider()),
		token.WithStaticCredentialsProvider(awscredprovider.NewStaticCredentialsProvider()),
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

func (e *Entrypoint) createPromoteFactory() (factory.PromoteFactory, error) {

	errContext := "(promote::entrypoint::createPromoteFactory)"

	dockerClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	copyCmd := copy.NewDockerImageCopyCmd(dockerClient)
	copyCmdFacade := godockerbuilder.NewDockerCopy(copyCmd)
	promoteRepoDocker := docker.NewDockerPromote(copyCmdFacade, os.Stdout)
	promoteRepoDryRun := dryrun.NewDryRunPromote(os.Stdout)
	promoteRepoFactory := factory.NewPromoteFactory()
	err = promoteRepoFactory.Register(image.DockerPromoterName, promoteRepoDocker)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	err = promoteRepoFactory.Register(image.DryRunPromoterName, promoteRepoDryRun)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return promoteRepoFactory, nil
}

func (e *Entrypoint) createSemanticVersionFactory() (*semver.SemVerGenerator, error) {
	return semver.NewSemVerGenerator(), nil
}
