package credentials

import (
	"context"

	handler "github.com/gostevedore/stevedore/internal/handler/create/credentials"
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
func (e *MockCreateCredentialsEntrypoint) Execute(ctx context.Context,
	args []string,
	conf *configuration.Configuration,
	inputEntrypointOptions *Options,
	inputHandlerOptions *handler.Options) error {
	res := e.Called(ctx, args, conf, inputEntrypointOptions, inputHandlerOptions)
	return res.Error(0)
}
