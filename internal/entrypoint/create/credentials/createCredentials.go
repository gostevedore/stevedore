package credentials

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/create/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/create/credentials"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/compatibility/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	credentialsformatfactory "github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/factory"
	credentialsstoreencryption "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
	credentialsenvvarsstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars"
	credentialsenvvarsstorebackend "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars/backend"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
)

const (
	getPasswordInputMessage           = "Password: "
	getAWSSecretAccessKeyInputMessage = "AWS Secret Access Key: "
	getPrivateKeyPasswordInputMessage = "Password: "
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *CreateCredentialsEntrypoint)

// CreateCredentialsEntrypoint defines the entrypoint for the application
type CreateCredentialsEntrypoint struct {
	console       Consoler
	compatibility Compatibilitier
	fs            afero.Fs
}

// NewCreateCredentialsEntrypoint returns a new entrypoint
func NewCreateCredentialsEntrypoint(opts ...OptionsFunc) *CreateCredentialsEntrypoint {
	e := &CreateCredentialsEntrypoint{}
	e.Options(opts...)

	return e
}

// Options provides the options for the entrypoint
func (e *CreateCredentialsEntrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// WithConsole sets the console writer/reader for the entrypoint
func WithConsole(console Consoler) OptionsFunc {
	return func(e *CreateCredentialsEntrypoint) {
		e.console = console
	}
}

// WithFileSystem sets the writer for the entrypoint
func WithFileSystem(fs afero.Fs) OptionsFunc {
	return func(e *CreateCredentialsEntrypoint) {
		e.fs = fs
	}
}

// WithCompatibility sets the compatibility for the entrypoint
func WithCompatibility(c Compatibilitier) OptionsFunc {
	return func(e *CreateCredentialsEntrypoint) {
		e.compatibility = c
	}
}

// Execute provides a mock function
func (e *CreateCredentialsEntrypoint) Execute(
	ctx context.Context,
	args []string,
	conf *configuration.Configuration,
	inputEntrypointOptions *Options,
	inputHandlerOptions *handler.Options,
) error {
	var err error
	var handlerOptions *handler.Options
	var id string
	var credentialsStore application.CredentialsStorer

	errContext := "(create::credentials::entrypoint::Execute)"

	id, err = e.prepareCredentialsId(args, inputEntrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	handlerOptions, err = e.prepareHandlerOptions(inputEntrypointOptions, inputHandlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	conf, err = e.prepareConfiguration(conf, inputEntrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	credentialsStore, err = e.createCredentialsStore(conf, inputEntrypointOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	app := application.NewCreateCredentialsApplication(
		application.WithCredentialsStore(credentialsStore),
	)

	h := handler.NewCreateCredentialsHandler(
		handler.WithApplication(app),
	)
	err = h.Handler(ctx, id, handlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if e.console != nil {
		e.console.Info(fmt.Sprintf("Credentials '%s' successfully created", id))
	}

	return nil
}

func (e *CreateCredentialsEntrypoint) prepareCredentialsId(args []string, options *Options) (string, error) {

	errContext := "(create::credentials::entrypoint:::prepareCredentialsId)"

	if options != nil && options.DEPRECATEDRegistryHost != "" {
		return options.DEPRECATEDRegistryHost, nil
	}

	if len(args) < 1 || args == nil {
		return "", errors.New(errContext, "To execute the create credentials entrypoint, an argument with credential id is required")
	}

	id := args[0]
	if len(args) > 1 {
		e.console.Warn(fmt.Sprintf("Ignoring extra arguments: %v\n", args[1:]))
	}

	return id, nil
}

func (e *CreateCredentialsEntrypoint) prepareConfiguration(conf *configuration.Configuration, options *Options) (*configuration.Configuration, error) {

	errContext := "(create::credentials::entrypoint::prepareConfiguration)"

	if options == nil {
		return nil, errors.New(errContext, "Entrypoint options must be provided to prepare configuration")
	}

	if conf == nil {
		return nil, errors.New(errContext, "Configuration must be provided to prepare configuration")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "Configuration credentials must be provided to prepare configuration")
	}

	if conf.Credentials.StorageType == "" {
		return nil, errors.New(errContext, "Credentials storage type must be provided to prepare configuration")
	}

	switch conf.Credentials.StorageType {
	case credentials.LocalStore:
		if options.LocalStoragePath != "" {
			conf.Credentials.LocalStoragePath = options.LocalStoragePath
		}
	}

	return conf, nil
}

// getPassword ask for password
func (e *CreateCredentialsEntrypoint) getPassword() (string, error) {

	errContext := "(create::credentials::entrypoint::getPassword)"

	if e.console == nil {
		return "", errors.New(errContext, "Console must be provided to execute create credentials entrypoint")
	}

	password, err := e.console.ReadPassword(getPasswordInputMessage)
	if err != nil {
		return "", errors.New(errContext, "Error reading password", err)
	}
	fmt.Fprintln(e.console)

	return password, nil

}

// getAWSSecretAccessKey ask for aws secret access key
func (e *CreateCredentialsEntrypoint) getAWSSecretAccessKey() (string, error) {

	errContext := "(create::credentials::entrypoint::getAWSSecretAccessKey)"

	if e.console == nil {
		return "", errors.New(errContext, "Console must be provided to execute create credentials entrypoint")
	}

	awsSecretAccessKey, err := e.console.ReadPassword(getAWSSecretAccessKeyInputMessage)
	if err != nil {
		return "", errors.New(errContext, "Error reading AWS secret access key", err)
	}
	fmt.Fprintln(e.console)

	return awsSecretAccessKey, nil
}

// getPrivateKeyPassword ask for private key password
func (e *CreateCredentialsEntrypoint) getPrivateKeyPassword() (string, error) {

	errContext := "(create::credentials::entrypoint::getPrivateKeyPassword)"

	if e.console == nil {
		return "", errors.New(errContext, "Console must be provided to execute create credentials entrypoint")
	}

	privateKeyPassword, err := e.console.ReadPassword(getPrivateKeyPasswordInputMessage)
	if err != nil {
		return "", errors.New(errContext, "Error reading private key password", err)
	}
	fmt.Fprintln(e.console)

	return privateKeyPassword, nil
}

// prepareHandlerOptions set handler options before execute the handler
func (e *CreateCredentialsEntrypoint) prepareHandlerOptions(inputEntrypointOptions *Options, inputHandlerOptions *handler.Options) (*handler.Options, error) {
	var password, awsSecretAccessKey, privateKeyPassword string
	var err error

	errContext := "(create::credentials::entrypoint::prepareHandlerOptions)"

	if inputHandlerOptions == nil {
		return nil, errors.New(errContext, "Handler options must be provided to execute create credentials entrypoint")
	}

	if inputEntrypointOptions == nil {
		return nil, errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint")
	}

	options := &handler.Options{}
	options.AllowUseSSHAgent = inputHandlerOptions.AllowUseSSHAgent
	options.AWSAccessKeyID = inputHandlerOptions.AWSAccessKeyID
	options.AWSProfile = inputHandlerOptions.AWSProfile
	options.AWSRegion = inputHandlerOptions.AWSRegion
	options.AWSRoleARN = inputHandlerOptions.AWSRoleARN
	if len(inputHandlerOptions.AWSSharedConfigFiles) > 0 {
		options.AWSSharedConfigFiles = append([]string{}, inputHandlerOptions.AWSSharedConfigFiles...)
	}
	if len(inputHandlerOptions.AWSSharedCredentialsFiles) > 0 {
		options.AWSSharedCredentialsFiles = append([]string{}, inputHandlerOptions.AWSSharedCredentialsFiles...)
	}
	options.AWSUseDefaultCredentialsChain = inputHandlerOptions.AWSUseDefaultCredentialsChain
	options.GitSSHUser = inputHandlerOptions.GitSSHUser
	options.PrivateKeyFile = inputHandlerOptions.PrivateKeyFile
	options.PrivateKeyPassword = inputHandlerOptions.PrivateKeyPassword
	options.Username = inputHandlerOptions.Username

	if inputHandlerOptions.Username != "" {
		password, err = e.getPassword()
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		options.Password = password
	}

	if inputHandlerOptions.AWSAccessKeyID != "" {
		awsSecretAccessKey, err = e.getAWSSecretAccessKey()
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		options.AWSSecretAccessKey = awsSecretAccessKey
	}

	if inputEntrypointOptions.AskPrivateKeyPassword {
		privateKeyPassword, err = e.getPrivateKeyPassword()
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		options.PrivateKeyPassword = privateKeyPassword
	}

	return options, nil
}

func (e *CreateCredentialsEntrypoint) createCredentialsFormater(conf *configuration.CredentialsConfiguration) (repository.Formater, error) {
	errContext := "(create::credentials::entrypoint::createCredentialsFormater)"

	if conf.Format == "" {
		return nil, errors.New(errContext, "To create credentials store in the create credentials entrypoint, credentials format must be specified")
	}

	credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
	credentialsFormat, err := credentialsFormatFactory.Get(conf.Format)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credentialsFormat, nil
}

func (e *CreateCredentialsEntrypoint) createCredentialsEnvvarsStore(conf *configuration.CredentialsConfiguration) (*credentialsenvvarsstore.EnvvarsStore, error) {
	errContext := "(create::credentials::entrypoint::createCredentialsEnvvarsStore)"
	var err error
	var format repository.Formater

	format, err = e.createCredentialsFormater(conf)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	encryption := credentialsstoreencryption.NewEncryption(
		credentialsstoreencryption.WithKey(conf.EncryptionKey),
	)

	store := credentialsenvvarsstore.NewEnvvarsStore(
		credentialsenvvarsstore.WithBackend(credentialsenvvarsstorebackend.NewOSEnvvarsBackend()),
		credentialsenvvarsstore.WithConsole(e.console),
		credentialsenvvarsstore.WithEncryption(encryption),
		credentialsenvvarsstore.WithFormater(format),
	)

	return store, nil
}

func (e *CreateCredentialsEntrypoint) createCredentialsLocalStore(comp credentialslocalstore.CredentialsCompatibilier, conf *configuration.CredentialsConfiguration) (*local.LocalStore, error) {
	var err error
	var format repository.Formater

	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	if comp == nil {
		return nil, errors.New(errContext, "To create the credentials local store, credentials compatibilitier is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials local store, credentials configuration is required")
	}

	if conf.LocalStoragePath == "" {
		return nil, errors.New(errContext, "To create the credentials local store, local storage path is required")
	}

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials local store, filesystem is required")
	}

	format, err = e.createCredentialsFormater(conf)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	localStoreOpts := []credentialslocalstore.OptionsFunc{
		credentialslocalstore.WithCompatibility(comp),
		credentialslocalstore.WithFilesystem(e.fs),
		credentialslocalstore.WithFormater(format),
		credentialslocalstore.WithPath(conf.LocalStoragePath),
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

func (e *CreateCredentialsEntrypoint) createCredentialsStore(conf *configuration.Configuration, options *Options) (application.CredentialsStorer, error) {

	var store application.CredentialsStorer
	var err error
	var credentialsFormat repository.Formater
	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials store, configuration is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials store, credentials configuration is required")
	}

	if conf.Credentials.StorageType == "" {
		return nil, errors.New(errContext, "To create the credentials store, credentials storage type must be defined")
	}

	if options == nil {
		return nil, errors.New(errContext, "To create the credentials store, options are required")
	}

	switch conf.Credentials.StorageType {
	case credentials.LocalStore:

		if e.compatibility == nil {
			return nil, errors.New(errContext, "To create the credentials store, compatibilitier is required")
		}

		if conf.Credentials.Format == "" {
			return nil, errors.New(errContext, "To create the credentials store, credentials format must be defined")
		}

		credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)
		credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
		credentialsFormat, err = credentialsFormatFactory.Get(conf.Credentials.Format)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		if options.ForceCreate {
			store, err = e.createCredentialsLocalStore(credentialsCompatibility, conf.Credentials)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}
		} else {
			store, err = e.createCredentialsLocalStoreWithSafeStore(credentialsCompatibility, conf.Credentials, credentialsFormat)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}
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

func (e *CreateCredentialsEntrypoint) createCredentialsLocalStoreWithSafeStore(comp credentialslocalstore.CredentialsCompatibilier, conf *configuration.CredentialsConfiguration, format repository.Formater) (*local.LocalStoreWithSafeStore, error) {

	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	if comp == nil {
		return nil, errors.New(errContext, "To create the credentials local store, credentials compatibilitier is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials local store, credentials configuration is required")
	}

	if conf.LocalStoragePath == "" {
		return nil, errors.New(errContext, "To create the credentials local store, local storage path is required")
	}

	if format == nil {
		return nil, errors.New(errContext, "To create the credentials local store, formater is required")
	}

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials local store, filesystem is required")
	}

	localStoreOpts := []credentialslocalstore.OptionsFunc{
		credentialslocalstore.WithFilesystem(e.fs),
		credentialslocalstore.WithCompatibility(comp),
		credentialslocalstore.WithPath(conf.LocalStoragePath),
		credentialslocalstore.WithFormater(format),
	}

	if conf.EncryptionKey != "" {
		encryption := credentialsstoreencryption.NewEncryption(
			credentialsstoreencryption.WithKey(conf.EncryptionKey),
		)

		localStoreOpts = append(localStoreOpts, credentialslocalstore.WithEncryption(encryption))
	}

	store := credentialslocalstore.NewLocalStoreWithSafeStore(localStoreOpts...)

	return store, nil
}
