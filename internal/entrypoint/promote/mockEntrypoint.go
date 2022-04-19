package entrypoint

import (
	"context"

	"github.com/gostevedore/stevedore/internal/configuration"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
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
func (e *MockEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, options *handler.Options) error {
	res := e.Called(ctx, args, conf, options)
	return res.Error(0)
}
