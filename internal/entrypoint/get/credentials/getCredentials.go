package credentials

import (
	"context"
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/get/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/credentials/compatibility"
	credentialsformatfactory "github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/factory"
	outputcredentials "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials"
	awsdefaultchain "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/AWSDefaultCredentialsChain"
	awsrolearn "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/AWSRoleARN"
	awsstaticcredentials "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/AWSStaticCredentials"
	sshagent "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/SSHAgent"
	privatekeyfile "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/privateKeyFile"
	usernamepassword "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/usernamePassword"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the application
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
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration) error {
	var err error
	var credentialsStore repository.CredentialsFilterer
	errContext := "(get::credentials::entrypoint::Execute)"

	if e.writer == nil {
		return errors.New(errContext, "To execute the entrypoint, a writer is required")
	}

	credentialsStore, err = e.createCredentialsFilter(conf)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	writer := console.NewConsole(e.writer, nil)
	output := outputcredentials.NewOutput(writer,
		usernamepassword.NewUsernamePasswordOutput(),
		awsstaticcredentials.NewAWSStaticCredentialsOutput(),
		awsrolearn.NewAWSRoleARNOutput(),
		awsdefaultchain.NewAWSDefaultCredentialsChainOutput(),
		privatekeyfile.NewPrivateKeyFileOutput(),
		sshagent.NewSSHAgentOutput(),
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

func (e *Entrypoint) createCredentialsLocalStore(conf *configuration.CredentialsConfiguration) (*credentialslocalstore.LocalStore, error) {

	errContext := "(get::credentials::entrypoint::createCredentialsLocalStore)"

	if conf == nil {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, credentials configuration is required")
	}

	if conf.Format == "" {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, credentials format must be specified")
	}

	if e.compatibility == nil {
		return nil, errors.New(errContext, "To create credentials local store in the entrypoint, compatibilitier is required")
	}

	credentialsCompatibility := credentialscompatibility.NewCredentialsCompatibility(e.compatibility)
	credentialsFormatFactory := credentialsformatfactory.NewFormatFactory()
	credentialsFormat, err := credentialsFormatFactory.Get(credentials.JSONFormat)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	store := credentialslocalstore.NewLocalStore(e.fs, conf.LocalStoragePath, credentialsFormat, credentialsCompatibility)

	return store, nil
}

func (e *Entrypoint) createCredentialsFilter(conf *configuration.Configuration) (repository.CredentialsFilterer, error) {
	errContext := "(get::credentials::entrypoint::createCredentialsFilter)"

	if e.fs == nil {
		return nil, errors.New(errContext, "To create the credentials filter in the entrypoint, a file system is required")
	}

	if conf == nil {
		return nil, errors.New(errContext, "To create the credentials filter in the entrypoint, configuration is required")
	}

	if conf.Credentials == nil {
		return nil, errors.New(errContext, "To create the credentials filter in the entrypoint, credentials configuration is required")
	}

	switch conf.Credentials.StorageType {
	case credentials.LocalStore:
		if conf.Credentials.LocalStoragePath == "" {
			return nil, errors.New(errContext, "To create credentials local store in the entrypoint, local storage path is required")
		}

		// create credentials store
		localstore, err := e.createCredentialsLocalStore(conf.Credentials)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		return localstore, nil
	default:
		return nil, errors.New(errContext, fmt.Sprintf("Unsupported credentials storage type '%s'", conf.Credentials.StorageType))
	}
}
