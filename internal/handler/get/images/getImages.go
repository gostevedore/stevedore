package images

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/images"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*GetImagesHandler)

// GetImagesHandler is a handler for command
type GetImagesHandler struct {
	app Applicationer
}

// NewGetImagesHandler creates a new handler for command
func NewGetImagesHandler(options ...OptionsFunc) *GetImagesHandler {
	handler := &GetImagesHandler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *GetImagesHandler) {
		h.app = app
	}
}

// Options configure the service
func (h *GetImagesHandler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// GetImagesHandler handles build commands
func (h *GetImagesHandler) Handler(ctx context.Context, options *Options) error {
	var err error

	errContext := "(handler::get::images::Handler)"

	if options == nil {
		return errors.New(errContext, "Get images handler requires handler options")
	}

	if h.app == nil {
		return errors.New(errContext, "Get images handler requires an application")
	}

	applicationOptions := &application.Options{}

	applicationOptions.Filter = append([]string{}, options.Filter...)

	err = h.app.Run(ctx, applicationOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
