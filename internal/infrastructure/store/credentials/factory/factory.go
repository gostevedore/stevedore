package factory

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

type CredentialsStoreFactory struct {
	store map[string]repository.CredentialsStorer
}

func NewCredentialsStoreFactory() *CredentialsStoreFactory {
	return &CredentialsStoreFactory{
		store: make(map[string]repository.CredentialsStorer),
	}
}

func (f *CredentialsStoreFactory) Register(id string, store repository.CredentialsStorer) error {

	errContext := "(store::credentials::factory::Register)"

	if f.store == nil {
		f.store = make(map[string]repository.CredentialsStorer)
	}

	if _, exists := f.store[id]; exists {
		return errors.New(errContext, fmt.Sprintf("Credentials store with id '%s' already exists", id))
	}

	f.store[id] = store

	return nil
}

func (f *CredentialsStoreFactory) Get(id string) (repository.CredentialsStorer, error) {

	errContext := "(store::credentials::factory::Get)"

	if f.store == nil {
		return nil, errors.New(errContext, "Credentials factory store is not initialized")
	}

	store, exist := f.store[id]
	if !exist {
		return nil, errors.New(errContext, fmt.Sprintf("Credentials store with id '%s' does not exist", id))
	}

	return store, nil
}
