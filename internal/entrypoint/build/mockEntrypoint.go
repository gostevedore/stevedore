package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/configuration"
	build "github.com/gostevedore/stevedore/internal/handler/build"
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
func (e *MockEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, entrypointOptions *EntrypointOptions, handlerOptions *build.HandlerOptions) error {
	res := e.Called(ctx, args, conf, entrypointOptions, handlerOptions)
	return res.Error(0)
}
