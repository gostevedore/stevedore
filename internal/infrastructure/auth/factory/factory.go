package factory

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// AuthFactory is a factory for auth providers
type AuthFactory struct {
	store                repository.CredentialsStorer
	credentialsProviders []repository.AuthProviderer
}

// NewAuthFactory creates a new auth provider factory
func NewAuthFactory(store repository.CredentialsStorer, auth ...repository.AuthProviderer) *AuthFactory {

	factory := &AuthFactory{
		store: store,
	}

	factory.credentialsProviders = append([]repository.AuthProviderer{}, auth...)

	return factory
}

// Get returns a new auth provider
func (f *AuthFactory) Get(id string) (repository.AuthMethodReader, error) {

	var err error
	var badge *credentials.Credential
	var method repository.AuthMethodReader
	errContext := "(credentials::factory::AuthFactory::Get)"

	if id == "" {
		return nil, errors.New(errContext, "To get credentials, you must provide an id")
	}

	badge, err = f.store.Get(id)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	for _, provider := range f.credentialsProviders {

		method, err = provider.Get(badge)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		if method != nil {
			return method, nil
		}
	}

	return nil, nil
}
