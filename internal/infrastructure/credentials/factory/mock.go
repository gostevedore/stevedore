package factory

import (
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/mock"
)

// MockCredentialsFactory is a mock of CredentialsFactory interface
type MockCredentialsFactory struct {
	mock.Mock
}

// NewMockCredentialsFactory creates a new auth provider factory
func NewMockCredentialsFactory() *MockCredentialsFactory {
	return &MockCredentialsFactory{}
}

// Get provides a mock function with given fields: id
func (f *MockCredentialsFactory) Get(id string) (repository.AuthMethodReader, error) {
	args := f.Called(id)
	return args.Get(0).(repository.AuthMethodReader), args.Error(1)
}
