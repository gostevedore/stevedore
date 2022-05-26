package build

import (
	"context"

	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockEntrypoint is a mock of Entrypoint interface
type MockEntrypoint struct {
	mock.Mock
}

// NewMockEntrypoint provides a mock of Entrypoint interface
func NewMockEntrypoint() *MockEntrypoint {
	return &MockEntrypoint{}
}

// Execute provides a mock function
func (e *MockEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, compatibility Compatibilitier, entrypointOptions *Options, handlerOptions *handler.Options) error {
	res := e.Called(ctx, args, conf, compatibility, entrypointOptions, handlerOptions)
	return res.Error(0)
}
