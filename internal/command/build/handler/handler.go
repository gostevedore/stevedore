package build

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/engine/build"
)

// Handler is a handler for build commands
type Handler struct {
	service ServiceBuilder
}

// NewHandler creates a new handler for build commands
func NewHandler(p ServiceBuilder) *Handler {
	return &Handler{
		service: p,
	}
}

// Handler handles build commands
func (h *Handler) Handler(ctx context.Context, options *HandlerOptions) error {

	errContext := "(build::Handler)"
	buildServiceOptions := &build.ServiceOptions{}

	err := h.service.Build(ctx, options.ImageName, options.Versions, buildServiceOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
