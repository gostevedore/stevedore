package factory

import (
	"fmt"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/mock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	errContext := "(store::credentials::factory::Register)"

	tests := []struct {
		desc    string
		factory *CredentialsStoreFactory
		id      string
		store   repository.CredentialsStorer
		res     repository.CredentialsStorer
		err     error
	}{
		{
			desc:    "Testing register an store to factory with an uninitialized factory store",
			factory: &CredentialsStoreFactory{},
			id:      credentials.LocalStore,
			store:   mock.NewMockStore(),
			res:     mock.NewMockStore(),
			err:     &errors.Error{},
		},
		{
			desc: "Testing error when registering an store to factory with an already registered id",
			factory: &CredentialsStoreFactory{
				store: map[string]repository.CredentialsStorer{
					credentials.LocalStore: mock.NewMockStore(),
				},
			},
			id:    credentials.LocalStore,
			store: mock.NewMockStore(),
			err:   errors.New(errContext, fmt.Sprintf("Credentials store with id '%s' already exists", credentials.LocalStore)),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.factory.Register(test.id, test.store)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.factory.store[test.id])
			}
		})
	}
}

func TestGet(t *testing.T) {

	errContext := "(store::credentials::factory::Get)"

	tests := []struct {
		desc    string
		factory *CredentialsStoreFactory
		id      string
		res     repository.CredentialsStorer
		err     error
	}{
		{
			desc:    "Testing error when getting an store from factory with an uninitialized factory store",
			factory: &CredentialsStoreFactory{},
			id:      credentials.LocalStore,
			err:     errors.New(errContext, "Credentials factory store is not initialized"),
		},
		{
			desc: "Testing error when getting an store from factory with an unexisting id",
			factory: &CredentialsStoreFactory{
				store: map[string]repository.CredentialsStorer{
					credentials.LocalStore: mock.NewMockStore(),
				},
			},
			id:  "unexisting",
			res: mock.NewMockStore(),
			err: errors.New(errContext, "Credentials store with id 'unexisting' does not exist"),
		},
		{
			desc: "Testing get a credentials store from factory",
			factory: &CredentialsStoreFactory{
				store: map[string]repository.CredentialsStorer{
					credentials.LocalStore: mock.NewMockStore(),
				},
			},
			id:  credentials.LocalStore,
			res: mock.NewMockStore(),
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.factory.Get(test.id)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, store)
			}
		})
	}
}
