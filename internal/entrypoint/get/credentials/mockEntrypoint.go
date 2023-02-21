package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

// MockEntrypoint is a mock of Entrypoint interface
type MockEntrypoint struct {
	mock.Mock
}

// NewMockEntrypoint provides an implementation Entrypoint interface
func NewMockEntrypoint() *MockEntrypoint {
	return &MockEntrypoint{}
}

// Execute provides a mock function
func (e *MockEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration, inputEntrypointOptions *Options) error {
	res := e.Called(ctx, args, conf)
	return res.Error(0)
}
