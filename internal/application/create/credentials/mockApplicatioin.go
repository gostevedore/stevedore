package credentials

import (
	"context"

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
func (m *MockCreateCredentialsApplication) Build(ctx context.Context, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, optionsFunc)
	return args.Error(0)
}
