package configuration

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGetConfigurationApplication is a mock of build application
type MockGetConfigurationApplication struct {
	mock.Mock
}

// NewMockGetConfigurationApplication return a mock of get configuration application
func NewMockGetConfigurationApplication() *MockGetConfigurationApplication {
	return &MockGetConfigurationApplication{}
}

// Run provides a mock function to carry out the application tasks
func (m *MockGetConfigurationApplication) Run(ctx context.Context, options *Options, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, options, optionsFunc)
	return args.Error(0)
}
