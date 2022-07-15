package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
)

// Handler is a handler for get credentials commands
type Handler struct {
	app GetCredentialsApplication
}

// NewHandler creates a new handler for build commands
func NewHandler(a GetCredentialsApplication) *Handler {
	return &Handler{
		app: a,
	}
}

// Handler handles build commands
func (h *Handler) Handler(ctx context.Context) error {
	var err error

	errContext := "(get::credentials::Handler)"

	err = h.app.Run(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
