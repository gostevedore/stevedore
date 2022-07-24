package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*CreateCredentialsHandler)

// CreateCredentialsHandler is a handler for get credentials commands
type CreateCredentialsHandler struct {
	app Applicationer
}

// NewCreateCredentialsHandler creates a new handler for build commands
func NewCreateCredentialsHandler(options ...OptionsFunc) *CreateCredentialsHandler {
	handler := &CreateCredentialsHandler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *CreateCredentialsHandler) {
		h.app = app
	}
}

// Options configure the service
func (h *CreateCredentialsHandler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// CreateCredentialsHandler handles build commands
func (h *CreateCredentialsHandler) Handler(ctx context.Context) error {
	var err error

	errContext := "(create/credentials::Handler)"

	err = h.app.Run(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
