package promote

import (
	"context"
	"fmt"
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
	authfactory "github.com/gostevedore/stevedore/internal/infrastructure/auth/factory"
	authmethodbasic "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
	authmethodkeyfile "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/keyfile"
	authmethodsshagent "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/sshagent"
	authproviderawsecr "github.com/gostevedore/stevedore/internal/infrastructure/auth/provider/awsecr"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/provider/awsecr/token"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/provider/awsecr/token/awscredprovider"
	authproviderstore "github.com/gostevedore/stevedore/internal/infrastructure/auth/provider/store"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/compatibility/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	credentialsformatfactory "github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker/godockerbuilder"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/factory"
	defaultreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/default"
	dockerreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	credentialsstoreencryption "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
	credentialsenvvarsstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars"
	credentialsenvvarsstorebackend "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars/backend"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the build application
type Entrypoint struct {
	fs            afero.Fs
	writer        ConsoleWriter
	compatibility Compatibilitier
}

// NewEntrypoint returns a new entrypoint
func NewEntrypoint(opts ...OptionsFunc) *Entrypoint {
	e := &Entrypoint{}
	e.Options(opts...)

	return e
}

// WithWriter sets the writer for the entrypoint
func WithWriter(w ConsoleWriter) OptionsFunc {
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
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, entrypointOptions *Options, handlerOptions *handler.Options) error {
	var err error
	var promoteRepoFactory factory.PromoteFactory
	var credentialsFactory repository.AuthFactorier
	var semverGenerator *semver.SemVerGenerator
	var options *handler.Options
	var referenceName repository.ImageReferenceNamer

	errContext := "(promote::entrypoint::Execute)"

	options, err = e.prepareHandlerOptions(args, conf, handlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	promoteRepoFactory, err = e.createPromoteFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	credentialsFactory, err = e.createAuthFactory(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	semverGenerator, err = e.createSemanticVersionFactory()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	referenceName, err = e.createReferenceName(entrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	promoteService := application.NewApplication(
		application.WithPromoteFactory(promoteRepoFactory),
		application.WithCredentials(credentialsFactory),
		application.WithSemver(semverGenerator),
		application.WithReferenceNamer(referenceName),
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
		e.writer.Warn(fmt.Sprintf("Ignoring extra arguments: %v\n", args[1:]))
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

func (e *Entrypoint) createCredentialsFormater(conf *configuration.CredentialsConfiguration) (repository.Formater, error) {
	errContext := "(promote::entrypoint::createCredentialsFormater)"

	if conf.Format == "" {
		return nil, errors.New(errContext, "To create credentials store in the entrypoint, credentials format must be specified")
	}

	credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
	credentialsFormat, err := credentialsFormatFactory.Get(conf.Format)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credentialsFormat, nil
}

func (e *Entrypoint) createCredentialsStore(conf *configuration.CredentialsConfiguration) (repository.CredentialsStorer, error) {

	var store repository.CredentialsStorer
	var err error
	var format repository.Formater

	errContext := "(promote::entrypoint::createCredentialsStore)"

	if conf == nil {
		return nil, errors.New(errContext, "To create credentials store in the promote entrypoint, credentials configuration is required")
	}

	format, err = e.createCredentialsFormater(conf)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	encryption := credentialsstoreencryption.NewEncryption(
		credentialsstoreencryption.WithKey(conf.EncryptionKey),
	)

	switch conf.StorageType {
	case credentials.LocalStore:

		if e.compatibility == nil {
			return nil, errors.New(errContext, "To create credentials store in the promote entrypoint, compatibility is required")
		}

		if conf.LocalStoragePath == "" {
			return nil, errors.New(errContext, "To create credentials store in the promote entrypoint, local storage path is required")
		}

		credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)

		localStoreOpts := []credentialslocalstore.OptionsFunc{
			credentialslocalstore.WithFilesystem(e.fs),
			credentialslocalstore.WithCompatibility(credentialsCompatibility),
			credentialslocalstore.WithPath(conf.LocalStoragePath),
			credentialslocalstore.WithFormater(format),
		}

		if conf.EncryptionKey != "" {
			localStoreOpts = append(localStoreOpts, credentialslocalstore.WithEncryption(encryption))
		}

		store = credentialslocalstore.NewLocalStore(localStoreOpts...)

	case credentials.EnvvarsStore:
		store = credentialsenvvarsstore.NewEnvvarsStore(
			credentialsenvvarsstore.WithConsole(e.writer),
			credentialsenvvarsstore.WithBackend(credentialsenvvarsstorebackend.NewOSEnvvarsBackend()),
			credentialsenvvarsstore.WithFormater(format),
			credentialsenvvarsstore.WithEncryption(encryption),
		)

	default:
		return nil, errors.New(errContext, fmt.Sprintf("Unsupported credentials storage type '%s'", conf.StorageType))
	}
	return store, nil
}

func (e *Entrypoint) createAuthFactory(conf *configuration.Configuration) (repository.AuthFactorier, error) {
	errContext := "(promote::entrypoint::createAuthFactory)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials store in the promote entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials store in the promote entrypoint, configuration is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials store in the promote entrypoint, credentials configuration is required")
	}

	// create credentials store
	store, err := e.createCredentialsStore(conf.Credentials)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	// create auth methods
	basic := authmethodbasic.NewBasicAuthMethod()
	keyfile := authmethodkeyfile.NewKeyFileAuthMethod()
	sshagent := authmethodsshagent.NewSSHAgentAuthMethod()

	// create auth providers
	badge := authproviderstore.NewStoreAuthProvider(basic, keyfile, sshagent)

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

	awsecr := authproviderawsecr.NewAWSECRAuthProvider(tokenProvider)

	// create credentials factory
	factory := authfactory.NewAuthFactory(store, badge, awsecr)

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

func (e *Entrypoint) createReferenceName(options *Options) (repository.ImageReferenceNamer, error) {
	if options.UseDockerNormalizedName {
		return dockerreferencename.NewDockerNormalizedReferenceName(), nil
	}

	return defaultreferencename.NewDefaultReferenceName(), nil
}
