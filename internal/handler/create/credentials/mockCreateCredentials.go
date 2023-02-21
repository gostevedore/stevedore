package credentials

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockCreateCredentialsHandler is a handler for create credentials commands
type MockCreateCredentialsHandler struct {
	mock.Mock
}

// NewMockCreateCredentialsHandler creates a new handler for create credentials commands
func NewMockCreateCredentialsHandler() *MockCreateCredentialsHandler {
	handler := &MockCreateCredentialsHandler{}

	return handler
}

// MockCreateCredentialsHandler handles create credentials commands
func (h *MockCreateCredentialsHandler) Handler(ctx context.Context, id string, options *Options) error {
	args := h.Called(ctx, id, options)
	return args.Error(0)
}
