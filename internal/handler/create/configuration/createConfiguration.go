package configuration

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*CreateConfigurationHandler)

// CreateConfigurationHandler is a handler for command
type CreateConfigurationHandler struct {
	app Applicationer
}

// NewCreateConfigurationHandler creates a new handler for command
func NewCreateConfigurationHandler(options ...OptionsFunc) *CreateConfigurationHandler {
	handler := &CreateConfigurationHandler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *CreateConfigurationHandler) {
		h.app = app
	}
}

// Options configure the service
func (h *CreateConfigurationHandler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// CreateConfigurationHandler handles build commands
func (h *CreateConfigurationHandler) Handler(ctx context.Context, options *Options) error {
	var err error

	errContext := "(handler::create::configuration::Handler)"

	if options == nil {
		return errors.New(errContext, "Create configuration handler requires the options parameter")
	}

	config := configuration.DefaultConfig()

	if len(options.BuildersPath) > 0 {
		config.BuildersPath = options.BuildersPath
	}

	if options.Concurrency > 0 {
		config.Concurrency = options.Concurrency
	}

	if len(options.CredentialsFormat) > 0 {
		config.Credentials.Format = options.CredentialsFormat
	}

	if len(options.CredentialsLocalStoragePath) > 0 {
		config.Credentials.LocalStoragePath = options.CredentialsLocalStoragePath
	}

	if len(options.CredentialsStorageType) > 0 {
		config.Credentials.StorageType = options.CredentialsStorageType
	}

	config.EnableSemanticVersionTags = options.EnableSemanticVersionTags

	if len(options.ImagesPath) > 0 {
		config.ImagesPath = options.ImagesPath
	}

	if len(options.LogPathFile) > 0 {
		config.LogPathFile = options.LogPathFile
	}

	config.PushImages = options.PushImages

	if len(options.SemanticVersionTagsTemplates) > 0 {
		config.SemanticVersionTagsTemplates = append([]string{}, options.SemanticVersionTagsTemplates...)
	}

	err = h.app.Run(ctx, config)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
