package command

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockBuildCommand is a mock of Command
type MockBuildCommand struct {
	mock.Mock
}

// NewMockBuildCommand creates a MockBuildCommand
func NewMockBuildCommand() *MockBuildCommand {
	return &MockBuildCommand{}
}

// Execute performs the action
func (c *MockBuildCommand) Execute(ctx context.Context) error {
	args := c.Called(ctx)
	return args.Error(0)
}
