package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockCreateCredentialsApplication is a mock of build application
type MockCreateCredentialsApplication struct {
	mock.Mock
}

// NewMockCreateCredentialsApplication return a mock of build application
func NewMockCreateCredentialsApplication() *MockCreateCredentialsApplication {
	return &MockCreateCredentialsApplication{}
}

// Build provides a mock function with given fields: ctx, buildPlan, name, version, options, optionsFunc
func (m *MockCreateCredentialsApplication) Run(ctx context.Context, id string, credential *credentials.Credential, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, id, credential, optionsFunc)
	return args.Error(0)
}
