package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockCreateCredentialsEntrypoint is a mock of Entrypoint interface
type MockCreateCredentialsEntrypoint struct {
	mock.Mock
}

// NewMockCreateCredentialsEntrypoint provides an implementation Entrypoint interface
func NewMockCreateCredentialsEntrypoint() *MockCreateCredentialsEntrypoint {
	return &MockCreateCredentialsEntrypoint{}
}

// Execute provides a mock function
func (e *MockCreateCredentialsEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration) error {
	res := e.Called(ctx, args, conf)
	return res.Error(0)
}
