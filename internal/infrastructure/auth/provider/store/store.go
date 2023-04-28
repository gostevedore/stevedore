package store

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// StoreAuthProvider return auth method from credential
type StoreAuthProvider struct {
	methods []repository.AuthMethodConstructor
}

// NewStoreAuthProvider return new instance of StoreAuthProvider
func NewStoreAuthProvider(methods ...repository.AuthMethodConstructor) *StoreAuthProvider {
	return &StoreAuthProvider{
		methods: methods,
	}
}

// Get return user password auth for docker registry
func (a *StoreAuthProvider) Get(credential *credentials.Credential) (repository.AuthMethodReader, error) {
	var err error
	var method repository.AuthMethodReader

	if credential == nil {
		return nil, nil
	}

	errContext := "(credentials::provider::StoreAuthProvider::Get)"

	for _, m := range a.methods {
		method, err = m.AuthMethodConstructor(credential)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		if method != nil {
			return method, nil
		}
	}

	return nil, nil
}
