package images

import (
	"context"

	handler "github.com/gostevedore/stevedore/internal/handler/get/images"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockGetImagesEntrypoint is a mock of Entrypoint interface
type MockGetImagesEntrypoint struct {
	mock.Mock
}

// NewMockGetImagesEntrypoint provides an implementation Entrypoint interface
func NewMockGetImagesEntrypoint() *MockGetImagesEntrypoint {
	return &MockGetImagesEntrypoint{}
}

// Execute provides a mock function
func (e *MockGetImagesEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, inputEntrypointOptions *Options, inputHandlerOptions *handler.Options) error {
	res := e.Called(ctx, args, conf, inputEntrypointOptions, inputHandlerOptions)
	return res.Error(0)
}
