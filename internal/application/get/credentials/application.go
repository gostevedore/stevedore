package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*Application)

// Application is an application service
type Application struct {
	credentials repository.CredentialsFilterer
	output      repository.CredentialsPrinter
}

// NewApplication creats a new application service
func NewApplication(options ...OptionsFunc) *Application {

	service := &Application{}
	service.Options(options...)

	return service
}

func WithCredentials(credentials repository.CredentialsFilterer) OptionsFunc {
	return func(a *Application) {
		a.credentials = credentials
	}
}

func WithOutput(output repository.CredentialsPrinter) OptionsFunc {
	return func(a *Application) {
		a.output = output
	}
}

// Options configure the service
func (a *Application) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Run method carries out the application tasks
func (a *Application) Run(ctx context.Context, optionsFunc ...OptionsFunc) error {

	errContext := "(application::get::credentials::Run)"

	a.Options(optionsFunc...)

	credentialsList := a.credentials.All()
	err := a.output.Print(credentialsList)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
