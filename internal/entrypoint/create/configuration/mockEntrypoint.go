package configuration

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockCreateConfigurationEntrypoint is a mock of Entrypoint interface
type MockCreateConfigurationEntrypoint struct {
	mock.Mock
}

// NewMockCreateConfigurationEntrypoint provides an implementation Entrypoint interface
func NewMockCreateConfigurationEntrypoint() *MockCreateConfigurationEntrypoint {
	return &MockCreateConfigurationEntrypoint{}
}

// Execute provides a mock function
func (e *MockCreateConfigurationEntrypoint) Execute(ctx context.Context, options *Options) error {
	res := e.Called(ctx, options)
	return res.Error(0)
}
