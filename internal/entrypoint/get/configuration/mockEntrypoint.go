package configuration

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockGetConfigurationEntrypoint is a mock of Entrypoint interface
type MockGetConfigurationEntrypoint struct {
	mock.Mock
}

// NewMockGetConfigurationEntrypoint provides an implementation Entrypoint interface
func NewMockGetConfigurationEntrypoint() *MockGetConfigurationEntrypoint {
	return &MockGetConfigurationEntrypoint{}
}

// Execute provides a mock function
func (e *MockGetConfigurationEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration) error {
	res := e.Called(ctx, args, conf)
	return res.Error(0)
}
