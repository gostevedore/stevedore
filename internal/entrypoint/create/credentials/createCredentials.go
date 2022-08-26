package credentials

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/create/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/create/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/credentials/compatibility"
	credentialsformatfactory "github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
)

const (
	getPasswordInputMessage           = "Password: "
	getAWSSecretAccessKeyInputMessage = "AWS Secret Access Key: "
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

	id, err = e.prepareCredentialsId(args)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	handlerOptions, err = e.prepareHandlerOptions(inputEntrypointOptions, inputHandlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	conf, err = e.prepareConfiguration(inputEntrypointOptions, conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	credentialsStore, err = e.createCredentialsStore(conf)
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

	return nil
}

func (e *CreateCredentialsEntrypoint) prepareCredentialsId(args []string) (string, error) {

	errContext := "(create::credentials::entrypoint:::prepareCredentialsId)"

	if len(args) < 1 || args == nil {
		return "", errors.New(errContext, "To execute the create credentials entrypoint, an argument with credential id is required")
	}

	id := args[0]
	if len(args) > 1 {
		e.console.Warn(fmt.Sprintf("Ignoring extra arguments: %v\n", args[1:]))
	}

	return id, nil
}

func (e *CreateCredentialsEntrypoint) prepareConfiguration(options *Options, conf *configuration.Configuration) (*configuration.Configuration, error) {
	if options.LocalStoragePath != "" {
		conf.Credentials.LocalStoragePath = options.LocalStoragePath
	}

	return conf, nil
}

// getPassword ask for password
func (e *CreateCredentialsEntrypoint) getPassword(options *Options) (string, error) {

	errContext := "(create::credentials::entrypoint::getPassword)"

	if options == nil {
		return "", errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint")
	}

	if e.console == nil {
		return "", errors.New(errContext, "Console must be provided to execute create credentials entrypoint")
	}

	password, err := e.console.ReadPassword(getPasswordInputMessage)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	return password, nil

}

// getAWSSecretAccessKey ask for aws secret access key
func (e *CreateCredentialsEntrypoint) getAWSSecretAccessKey(options *Options) (string, error) {

	errContext := "(create::credentials::entrypoint::getAWSSecretAccessKey)"

	if options == nil {
		return "", errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint")
	}

	if e.console == nil {
		return "", errors.New(errContext, "Console must be provided to execute create credentials entrypoint")
	}

	awsSecretAccessKey, err := e.console.ReadPassword(getAWSSecretAccessKeyInputMessage)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	return awsSecretAccessKey, nil

}

// prepareHandlerOptions set handler options before execute the handler
func (e *CreateCredentialsEntrypoint) prepareHandlerOptions(inputEntrypointOptions *Options, inputHandlerOptions *handler.Options) (*handler.Options, error) {
	var password, awsSecretAccessKey string
	var err error

	errContext := "(create::credentials::entrypoint::prepareHandlerOptions)"

	if inputEntrypointOptions == nil {
		return nil, errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint")
	}

	if inputHandlerOptions == nil {
		return nil, errors.New(errContext, "Handler options must be provided to execute create credentials entrypoint")
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

	if inputEntrypointOptions.AskPassword {
		password, err = e.getPassword(inputEntrypointOptions)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		options.Password = password
	}

	if inputEntrypointOptions.AskAWSSecretAccessKey {
		awsSecretAccessKey, err = e.getAWSSecretAccessKey(inputEntrypointOptions)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		options.AWSSecretAccessKey = awsSecretAccessKey
	}

	return options, nil
}

//
func (e *CreateCredentialsEntrypoint) createCredentialsStore(conf *configuration.Configuration) (application.CredentialsStorer, error) {

	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create the credentials store, compatibilitier is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials store, credentials configuration is required")
	}

	if conf.Credentials.Format == "" {
		return nil, errors.New(errContext, "To create the credentials store, credentials format must be defined")
	}

	if conf.Credentials.StorageType == "" {
		return nil, errors.New(errContext, "To create the credentials store, credentials storage type must be defined")
	}

	credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)
	credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
	credentialsFormat, err := credentialsFormatFactory.Get(conf.Credentials.Format)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	switch conf.Credentials.StorageType {
	case credentials.LocalStore:
		store, err := e.createCredentialsLocalStore(credentialsCompatibility, conf.Credentials, credentialsFormat)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		return store, nil
	default:
		return nil, errors.New(errContext, fmt.Sprintf("Unsupported credentials storage type '%s'", conf.Credentials.StorageType))
	}

}

func (e *CreateCredentialsEntrypoint) createCredentialsLocalStore(comp credentialslocalstore.CredentialsCompatibilier, conf *configuration.CredentialsConfiguration, format repository.Formater) (*local.LocalStore, error) {

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

	store := credentialslocalstore.NewLocalStore(e.fs, conf.LocalStoragePath, format, comp)

	return store, nil
}
