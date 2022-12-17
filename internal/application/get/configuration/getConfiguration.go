package configuration

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*GetConfigurationApplication)

// GetConfigurationApplication is an application service
type GetConfigurationApplication struct {
	write configuration.ConfigurationWriter
}

// NewGetConfigurationApplication creats a new application service
func NewGetConfigurationApplication(options ...OptionsFunc) *GetConfigurationApplication {

	app := &GetConfigurationApplication{}
	app.Options(options...)

	return app
}

func WithWrite(w configuration.ConfigurationWriter) OptionsFunc {
	return func(a *GetConfigurationApplication) {
		a.write = w
	}
}

// Options to configure the application
func (a *GetConfigurationApplication) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Run method carries out the application tasks
func (a *GetConfigurationApplication) Run(ctx context.Context, options *Options, optionsFunc ...OptionsFunc) error {

	errContext := "(application::get::configuration::Run)"

	if options.Configuration == nil {
		return errors.New(errContext, "Get configuration application requires a configuration")
	}

	err := a.write.Write(options.Configuration)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
