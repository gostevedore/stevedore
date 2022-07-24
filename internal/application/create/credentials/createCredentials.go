package credentials

import "context"

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*CreateCredentialsApplication)

// CreateCredentialsApplication is an application service
type CreateCredentialsApplication struct {
}

// NewCreateCredentialsApplication creats a new application service
func NewCreateCredentialsApplication(options ...OptionsFunc) *CreateCredentialsApplication {

	service := &CreateCredentialsApplication{}
	service.Options(options...)

	return service
}

// Options configure the service
func (a *CreateCredentialsApplication) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Build starts the building process
func (a *CreateCredentialsApplication) Run(ctx context.Context, optionsFunc ...OptionsFunc) error {
	return nil
}
