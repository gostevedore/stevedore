package configuration

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockCreateConfigurationApplication is a mock of the application
type MockCreateConfigurationApplication struct {
	mock.Mock
}

// NewMockCreateConfigurationApplication return a mock of the application
func NewMockCreateConfigurationApplication() *MockCreateConfigurationApplication {
	return &MockCreateConfigurationApplication{}
}

// Run provides a mock function to carry out the application tasks
func (m *MockCreateConfigurationApplication) Run(ctx context.Context, config *configuration.Configuration, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, config, optionsFunc)
	return args.Error(0)
}
