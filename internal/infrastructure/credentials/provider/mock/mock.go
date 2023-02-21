package mock

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/mock"
)

// MockAuthProvider return user password auth for docker registry
type MockAuthProvider struct {
	mock.Mock
}

// NewMockAuthProvider return new instance of MockAuthProvider
func NewMockAuthProvider() *MockAuthProvider {
	return &MockAuthProvider{}
}

// Get return user password auth for docker registry
func (m *MockAuthProvider) Get(badge *credentials.Badge) (repository.AuthMethodReader, error) {
	args := m.Called(badge)
	return args.Get(0).(repository.AuthMethodReader), args.Error(1)
}
