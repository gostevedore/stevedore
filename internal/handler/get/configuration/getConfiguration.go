package configuration

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/configuration"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*GetConfigurationHandler)

// GetConfigurationHandler is a handler for command
type GetConfigurationHandler struct {
	app Applicationer
}

// NewGetConfigurationHandler creates a new handler for command
func NewGetConfigurationHandler(options ...OptionsFunc) *GetConfigurationHandler {
	handler := &GetConfigurationHandler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *GetConfigurationHandler) {
		h.app = app
	}
}

// Options configure the service
func (h *GetConfigurationHandler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// GetConfigurationHandler handles build commands
func (h *GetConfigurationHandler) Handler(ctx context.Context, options *Options) error {
	var err error

	errContext := "(get/configuration::Handler)"

	appOptions := &application.Options{}
	appOptions.Configuration = options.Configuration

	err = h.app.Run(ctx, appOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
