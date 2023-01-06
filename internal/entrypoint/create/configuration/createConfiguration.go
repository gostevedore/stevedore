package configuration

import (
	"context"
	"path/filepath"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/create/configuration"
	handler "github.com/gostevedore/stevedore/internal/handler/create/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	output "github.com/gostevedore/stevedore/internal/infrastructure/configuration/output/file"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *CreateConfigurationEntrypoint)

// CreateConfigurationEntrypoint defines the entrypoint for the application
type CreateConfigurationEntrypoint struct {
	fs afero.Fs
}

// NewCreateConfigurationEntrypoint returns a new entrypoint
func NewCreateConfigurationEntrypoint(opts ...OptionsFunc) *CreateConfigurationEntrypoint {
	e := &CreateConfigurationEntrypoint{}
	e.Options(opts...)

	return e
}

// Options provides the options for the entrypoint
func (e *CreateConfigurationEntrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// WithFileSystem sets the writer for the entrypoint
func WithFileSystem(fs afero.Fs) OptionsFunc {
	return func(e *CreateConfigurationEntrypoint) {
		e.fs = fs
	}
}

// Execute is a pseudo-main method for the command
func (e *CreateConfigurationEntrypoint) Execute(ctx context.Context, options *Options) error {
	var err error
	var handlerOptions *handler.Options
	var createConfigurationHandler *handler.CreateConfigurationHandler
	var createConfigurationApplication *application.CreateConfigurationApplication
	var writer configuration.ConfigurationWriter

	errContext := "(entrypoint::create::configuration::Execute)"

	handlerOptions, err = e.prepareHandlerOptions(options)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	writer, err = e.createOutputWriter(options)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	createConfigurationApplication = application.NewCreateConfigurationApplication(
		application.WithWrite(writer),
	)

	createConfigurationHandler = handler.NewCreateConfigurationHandler(
		handler.WithApplication(createConfigurationApplication),
	)
	err = createConfigurationHandler.Handler(ctx, handlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (e *CreateConfigurationEntrypoint) prepareHandlerOptions(options *Options) (*handler.Options, error) {

	errContext := "(entrypoint::create::configuration::prepareHandlerOptions)"

	if options == nil {
		return nil, errors.New(errContext, "Create configuration entrypoint requires options to prepare handler options")
	}

	handlerOptions := &handler.Options{}

	handlerOptions.BuildersPath = options.BuildersPath
	handlerOptions.Concurrency = options.Concurrency
	handlerOptions.CredentialsFormat = options.CredentialsFormat
	handlerOptions.CredentialsLocalStoragePath = options.CredentialsLocalStoragePath
	handlerOptions.CredentialsStorageType = options.CredentialsStorageType
	handlerOptions.EnableSemanticVersionTags = options.EnableSemanticVersionTags
	handlerOptions.ImagesPath = options.ImagesPath
	handlerOptions.LogPathFile = options.LogPathFile
	handlerOptions.PushImages = options.PushImages
	handlerOptions.SemanticVersionTagsTemplates = append([]string{}, options.SemanticVersionTagsTemplates...)

	return handlerOptions, nil
}

func (e *CreateConfigurationEntrypoint) getConfigurationFileName(options *Options) (string, error) {

	errContext := "(entrypoint::create::configuration::getConfigurationFileName)"

	if options == nil {
		return "", errors.New(errContext, "Create configuration entrypoint requires options to get configuration file name")
	}

	fileName := filepath.Join(
		configuration.DefaultConfigFolder,
		strings.Join([]string{
			configuration.DefaultConfigFile,
			configuration.DefaultConfigFileExtention,
		}, "."))

	if len(options.ConfigurationFilePath) > 0 {
		fileName = options.ConfigurationFilePath
	}

	return fileName, nil
}

func (e *CreateConfigurationEntrypoint) createOutputWriter(options *Options) (configuration.ConfigurationWriter, error) {
	var writer configuration.ConfigurationWriter
	var fileName string
	var err error

	errContext := "(entrypoint::create::configuration::createOutputWriter)"

	if e.fs == nil {
		return nil, errors.New(errContext, "Create configuration entrypoint requires a filesystem to create the output writer")
	}

	if options == nil {
		return nil, errors.New(errContext, "Create configuration entrypoint requires options to create the output writer")
	}

	fileName, err = e.getConfigurationFileName(options)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if options.Force {
		writer = output.NewConfigurationFilePersist(
			output.WithFilePath(fileName),
			output.WithFileSystem(e.fs),
		)
	} else {
		writer = output.NewConfigurationFileSafePersist(
			output.WithFilePath(fileName),
			output.WithFileSystem(e.fs),
		)
	}

	return writer, nil
}
