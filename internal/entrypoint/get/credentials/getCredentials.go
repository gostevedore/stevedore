package credentials

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/get/credentials"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/compatibility/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	credentialsformatfactory "github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/factory"
	outputcredentials "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials"
	awsdefaultchain "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/AWSDefaultCredentialsChain"
	awsrolearn "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/AWSRoleARN"
	awsstaticcredentials "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/AWSStaticCredentials"
	sshagent "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/SSHAgent"
	privatekeyfile "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/privateKeyFile"
	usernamepassword "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/usernamePassword"
	credentialsstoreencryption "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
	credentialsenvvarsstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars"
	credentialsenvvarsstorebackend "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars/backend"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the application
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

// WithCompatibility set the
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

// Execute is a pseudo-main method for the command
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, inputEntrypointOptions *Options) error {
	var err error
	var credentialsStore repository.CredentialsFilterer
	errContext := "(get::credentials::entrypoint::Execute)"

	var usernamepasswordoutput outputcredentials.Outputter
	var awsstaticcredentialsoutput outputcredentials.Outputter
	var awsrolearnoutput outputcredentials.Outputter
	var awsdefaultchainoutput outputcredentials.Outputter
	var privatekeyfileoutput outputcredentials.Outputter
	var sshagentoutput outputcredentials.Outputter

	if e.writer == nil {
		return errors.New(errContext, "To execute the entrypoint, a writer is required")
	}

	credentialsStore, err = e.createCredentialsFilter(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	writer := console.NewConsole(e.writer, nil)

	usernamepasswordoutput = usernamepassword.NewUsernamePasswordOutput()
	awsstaticcredentialsoutput = awsstaticcredentials.NewAWSStaticCredentialsOutput()
	awsrolearnoutput = awsrolearn.NewAWSRoleARNOutput()
	awsdefaultchainoutput = awsdefaultchain.NewAWSDefaultCredentialsChainOutput()
	privatekeyfileoutput = privatekeyfile.NewPrivateKeyFileOutput()
	sshagentoutput = sshagent.NewSSHAgentOutput()

	if inputEntrypointOptions.ShowSecrets {
		usernamepasswordoutput = usernamepassword.NewUsernamePasswordWithSecretsOutput(usernamepasswordoutput.(*usernamepassword.UsernamePasswordOutput))

		awsstaticcredentialsoutput = awsstaticcredentials.NewAWSStaticCredentialsWithSecretsOutput(awsstaticcredentialsoutput.(*awsstaticcredentials.AWSStaticCredentialsOutput))

		awsrolearnoutput = awsrolearn.NewAWSRoleARNWithSecretsOutput(awsrolearnoutput.(*awsrolearn.AWSRoleARNOutput))

		privatekeyfileoutput = privatekeyfile.NewPrivateKeyFileWithSecretsOutput(privatekeyfileoutput.(*privatekeyfile.PrivateKeyFileOutput))
	}

	output := outputcredentials.NewOutput(writer,
		usernamepasswordoutput,
		awsstaticcredentialsoutput,
		awsrolearnoutput,
		awsdefaultchainoutput,
		privatekeyfileoutput,
		sshagentoutput,
	)

	getCredentialsApplication := application.NewApplication(
		application.WithCredentials(credentialsStore),
		application.WithOutput(output),
	)

	getCredentialsHandler := handler.NewHandler(
		handler.WithApplication(getCredentialsApplication),
	)
	err = getCredentialsHandler.Handler(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (e Entrypoint) createCredentialsFormater(conf *configuration.CredentialsConfiguration) (repository.Formater, error) {
	errContext := "(get::credentials::entrypoint::createCredentialsFormater)"

	if conf.Format == "" {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, credentials format must be specified")
	}

	credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
	credentialsFormat, err := credentialsFormatFactory.Get(conf.Format)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credentialsFormat, nil
}

func (e *Entrypoint) createCredentialsLocalStore(conf *configuration.CredentialsConfiguration) (*credentialslocalstore.LocalStore, error) {
	var err error
	var format repository.Formater

	errContext := "(get::credentials::entrypoint::createCredentialsLocalStore)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, credentials configuration is required")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, compatibilitier is required")
	}

	format, err = e.createCredentialsFormater(conf)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)

	localStoreOpts := []credentialslocalstore.OptionsFunc{
		credentialslocalstore.WithFilesystem(e.fs),
		credentialslocalstore.WithCompatibility(credentialsCompatibility),
		credentialslocalstore.WithPath(conf.LocalStoragePath),
		credentialslocalstore.WithFormater(format),
	}

	if conf.EncryptionKey != "" {
		encryption := credentialsstoreencryption.NewEncryption(
			credentialsstoreencryption.WithKey(conf.EncryptionKey),
		)

		localStoreOpts = append(localStoreOpts, credentialslocalstore.WithEncryption(encryption))
	}

	store := credentialslocalstore.NewLocalStore(localStoreOpts...)

	return store, nil
}

func (e *Entrypoint) createCredentialsEnvvarsStore(conf *configuration.CredentialsConfiguration) (*credentialsenvvarsstore.EnvvarsStore, error) {

	var err error
	var format repository.Formater

	errContext := "(get::credentials::entrypoint::createCredentialsEnvvarsStore)"

	format, err = e.createCredentialsFormater(conf)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	encryption := credentialsstoreencryption.NewEncryption(
		credentialsstoreencryption.WithKey(conf.EncryptionKey),
	)

	store := credentialsenvvarsstore.NewEnvvarsStore(
		credentialsenvvarsstore.WithBackend(credentialsenvvarsstorebackend.NewOSEnvvarsBackend()),
		credentialsenvvarsstore.WithEncryption(encryption),
		credentialsenvvarsstore.WithFormater(format),
	)

	return store, nil
}

func (e *Entrypoint) createCredentialsFilter(conf *configuration.Configuration) (repository.CredentialsFilterer, error) {
	errContext := "(get::credentials::entrypoint::createCredentialsFilter)"
	var store repository.CredentialsFilterer
	var err error

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials filter in the entrypoint, configuration is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials filter in the entrypoint, credentials configuration is required")
	}

	if conf.Credentials.Format == "" {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, credentials format must be specified")
	}

	switch conf.Credentials.StorageType {
	case credentials.LocalStore:
		if conf.Credentials.LocalStoragePath == "" {
			return nil, errors.New(errContext, "To create credentials local store in the entrypoint, local storage path is required")
		}

		// create credentials store
		store, err = e.createCredentialsLocalStore(conf.Credentials)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

	case credentials.EnvvarsStore:
		store, err = e.createCredentialsEnvvarsStore(conf.Credentials)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

	default:
		return nil, errors.New(errContext, fmt.Sprintf("Unsupported credentials storage type '%s'", conf.Credentials.StorageType))
	}
	return store, nil
}
