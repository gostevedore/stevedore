package configuration

import (
	"context"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/configuration"
	handler "github.com/gostevedore/stevedore/internal/handler/get/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/output/console"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *GetConfigurationEntrypoint)

// GetConfigurationEntrypoint defines the entrypoint for the application
type GetConfigurationEntrypoint struct {
	writer io.Writer
}

// NewGetConfigurationEntrypoint returns a new entrypoint
func NewGetConfigurationEntrypoint(opts ...OptionsFunc) *GetConfigurationEntrypoint {
	e := &GetConfigurationEntrypoint{}
	e.Options(opts...)

	return e
}

// WithWriter sets the writer for the entrypoint
func WithWriter(w io.Writer) OptionsFunc {
	return func(e *GetConfigurationEntrypoint) {
		e.writer = w
	}
}

// Options provides the options for the entrypoint
func (e *GetConfigurationEntrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute is a pseudo-main method for the command
func (e *GetConfigurationEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration) error {

	var err error
	var getConfigurationApplication *application.GetConfigurationApplication
	var getConfigurationHandler *handler.GetConfigurationHandler

	errContext := "(get::configuration::entrypoint::Execute)"

	console := console.NewConfigurationConsoleOutput(e.writer)
	getConfigurationApplication = application.NewGetConfigurationApplication(
		application.WithWrite(console),
	)

	getConfigurationHandler = handler.NewGetConfigurationHandler(
		handler.WithApplication(getConfigurationApplication),
	)

	handlerOptions := handler.Options{
		Configuration: conf,
	}

	err = getConfigurationHandler.Handler(ctx, &handlerOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
