package factory

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
)

type CredentialsStoreFactory struct {
	configuration *configuration.Configuration
	backend       map[string]repository.CredentialsStorer
}

func NewCredentialsStoreFactory(configuration *configuration.Configuration) *CredentialsStoreFactory {
	return &CredentialsStoreFactory{
		configuration: configuration,
	}
}

func (f *CredentialsStoreFactory) Register(id string, store repository.CredentialsStorer) error {

	errContext := "(factory::Register)"

	if f.backend == nil {
		f.backend = make(map[string]repository.CredentialsStorer)
	}

	if _, exists := f.backend[id]; exists {
		return errors.New(errContext, fmt.Sprintf("Credentials store with id '%s' already exists", id))
	}

	f.backend[id] = store

	return nil
}

func (f *CredentialsStoreFactory) Get() (repository.CredentialsStorer, error) {

	errContext := "(factory::Get)"

	if f.configuration.DockerCredentialsDir == "" {
		return nil, errors.New(errContext, "Docker credentials directory is not defined on configuration")
	}

	backend := f.backend[credentials.LocalStore]
	err := backend.(*local.LocalStore).LoadCredentials(f.configuration.DockerCredentialsDir)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return backend, nil
}
