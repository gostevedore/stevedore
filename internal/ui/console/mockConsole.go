package console

import (
	"github.com/stretchr/testify/mock"
)

type MockConsole struct {
	mock.Mock
}

// NewMockConsole returns a new mock console
func NewMockConsole() *MockConsole {
	return &MockConsole{}
}

// Write is a mock implementation of the Write method
func (c *MockConsole) Write(data []byte) (int, error) {
	args := c.Mock.Called(data)
	return args.Int(0), args.Error(1)
}

// PrintTable is a mock implementation of the PrintTable method
func (c *MockConsole) PrintTable(table [][]string) error {
	args := c.Mock.Called(table)
	return args.Error(0)
}
