package factory

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// CredentialsFactory is a factory for auth providers
type CredentialsFactory struct {
	store                repository.CredentialsStorer
	credentialsProviders []repository.CredentialsProviderer
}

// NewCredentialsFactory creates a new auth provider factory
func NewCredentialsFactory(store repository.CredentialsStorer, auth ...repository.CredentialsProviderer) *CredentialsFactory {

	factory := &CredentialsFactory{
		store: store,
	}

	factory.credentialsProviders = append([]repository.CredentialsProviderer{}, auth...)

	return factory
}

// Get returns a new auth provider
func (f *CredentialsFactory) Get(id string) (repository.AuthMethodReader, error) {

	var err error
	var badge *credentials.Badge
	var method repository.AuthMethodReader
	errContext := "(credentials::factory::CredentialsFactory::Get)"

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
