package credentials

import "context"

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*Application)

// Application is an application service
type Application struct {
}

// NewApplication creats a new application service
func NewApplication(options ...OptionsFunc) *Application {

	service := &Application{}
	service.Options(options...)

	return service
}

// Options configure the service
func (a *Application) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Build starts the building process
func (a *Application) Run(ctx context.Context, optionsFunc ...OptionsFunc) error {
	return nil
}
