package builders

import (
	"context"

	handler "github.com/gostevedore/stevedore/internal/handler/get/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockGetBuildersEntrypoint is a mock of Entrypoint interface
type MockGetBuildersEntrypoint struct {
	mock.Mock
}

// NewMockGetBuildersEntrypoint provides an implementation Entrypoint interface
func NewMockGetBuildersEntrypoint() *MockGetBuildersEntrypoint {
	return &MockGetBuildersEntrypoint{}
}

// Execute provides a mock function
func (e *MockGetBuildersEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, options *handler.Options) error {
	res := e.Called(ctx, args, conf, options)
	return res.Error(0)
}
