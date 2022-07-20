package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*Handler)

// Handler is a handler for get credentials commands
type Handler struct {
	app Applicationer
}

// NewHandler creates a new handler for build commands
func NewHandler(options ...OptionsFunc) *Handler {
	handler := &Handler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *Handler) {
		h.app = app
	}
}

// Options configure the service
func (h *Handler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// Handler handles build commands
func (h *Handler) Handler(ctx context.Context) error {
	var err error

	errContext := "(get::credentials::Handler)"

	if h.app == nil {
		return errors.New(errContext, "Handler application is not configured")
	}

	err = h.app.Run(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
