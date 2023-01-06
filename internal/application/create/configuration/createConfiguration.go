package configuration

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*CreateConfigurationApplication)

// CreateConfigurationApplication is an application service
type CreateConfigurationApplication struct {
	write configuration.ConfigurationWriter
}

// NewCreateConfigurationApplication creats a new application service
func NewCreateConfigurationApplication(options ...OptionsFunc) *CreateConfigurationApplication {

	app := &CreateConfigurationApplication{}
	app.Options(options...)

	return app
}

func WithWrite(w configuration.ConfigurationWriter) OptionsFunc {
	return func(a *CreateConfigurationApplication) {
		a.write = w
	}
}

// Options to configure the application
func (a *CreateConfigurationApplication) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Run method carries out the application tasks
func (a *CreateConfigurationApplication) Run(ctx context.Context, config *configuration.Configuration, optionsFunc ...OptionsFunc) error {
	var err error
	errContext := "(application::create::configuration::Run)"

	if config == nil {
		return errors.New(errContext, "Create configuration application requires a configuration")
	}

	err = config.ValidateConfiguration()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = a.write.Write(config)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
