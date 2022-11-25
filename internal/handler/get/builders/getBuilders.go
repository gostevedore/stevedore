package builders

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/builders"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*GetBuildersHandler)

// GetBuildersHandler is a handler for command
type GetBuildersHandler struct {
	app Applicationer
}

// NewGetBuildersHandler creates a new handler for command
func NewGetBuildersHandler(options ...OptionsFunc) *GetBuildersHandler {
	handler := &GetBuildersHandler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *GetBuildersHandler) {
		h.app = app
	}
}

// Options configure the service
func (h *GetBuildersHandler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// GetBuildersHandler handles build commands
func (h *GetBuildersHandler) Handler(ctx context.Context, options *Options) error {
	var err error

	errContext := "(handler::get::builders::Handler)"

	if options == nil {
		return errors.New(errContext, "Get builders handler requires handler options")
	}

	if h.app == nil {
		return errors.New(errContext, "Get builders handler requires an application")
	}

	applicationOptions := &application.Options{}

	applicationOptions.Filter = append([]string{}, options.Filter...)

	err = h.app.Run(ctx, applicationOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
