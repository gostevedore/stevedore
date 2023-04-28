package mock

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockStore is a mocked store for credentials
type MockStore struct {
	mock.Mock
}

// NewMockStore creates a new mocked store for credentials
func NewMockStore() *MockStore {
	return &MockStore{}
}

// Store stores a credential
func (m *MockStore) Store(id string, credential *credentials.Credential) error {
	args := m.Mock.Called(id, credential)
	return args.Error(0)
}

// Get returns a auth for the credential id
func (m *MockStore) Get(id string) (*credentials.Credential, error) {
	args := m.Mock.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*credentials.Credential), args.Error(1)
	}
}

// All returns all credentials
func (m *MockStore) All() ([]*credentials.Credential, error) {
	args := m.Mock.Called()
	return args.Get(0).([]*credentials.Credential), args.Error(1)
}
