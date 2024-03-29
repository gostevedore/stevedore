package credentials

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*CreateCredentialsApplication)

// CreateCredentialsApplication is an application service
type CreateCredentialsApplication struct {
	store CredentialsStorer
}

// NewCreateCredentialsApplication creats a new application service
func NewCreateCredentialsApplication(options ...OptionsFunc) *CreateCredentialsApplication {

	service := &CreateCredentialsApplication{}
	service.Options(options...)

	return service
}

// WithCredentialsStore provides a function to configure the credentials store
func WithCredentialsStore(store CredentialsStorer) OptionsFunc {
	return func(app *CreateCredentialsApplication) {
		app.store = store
	}
}

// Options configure the service
func (a *CreateCredentialsApplication) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Run method carries out the application tasks
func (a *CreateCredentialsApplication) Run(ctx context.Context, id string, credential *credentials.Credential, optionsFunc ...OptionsFunc) error {
	var err error

	errContext := "(application::create::credentials::Run)"

	if a.store == nil {
		return errors.New(errContext, "To run the create credentials application, a credentials storer must be provided")
	}

	if credential == nil {
		return errors.New(errContext, "To run the create credentials application, a credential must be provided")
	}

	if id == "" {
		return errors.New(errContext, "To run the create credentials application, a id for credentials must be provided")
	}

	_, err = credential.IsValid()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = a.store.Store(id, credential)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error storing '%s' credentials", id), err)
	}

	return nil
}
