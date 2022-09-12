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

// Info prints a info message
func (c *MockConsole) Info(msg ...interface{}) {
	c.Mock.Called(msg)
}

// Warn prints a warning message
func (c *MockConsole) Warn(msg ...interface{}) {
	c.Mock.Called(msg)
}

// Error prints an error message
func (c *MockConsole) Error(msg ...interface{}) {
	c.Mock.Called(msg)
}

// Debug prints a debug message
func (c *MockConsole) Debug(msg ...interface{}) {
	c.Mock.Called(msg)
}

// Read read a line from console reader
func (c *MockConsole) Read() string {
	args := c.Mock.Called()
	return args.String(0)
}

// ReadPassword read password from console reader
func (c *MockConsole) ReadPassword(prompt string) (string, error) {
	args := c.Mock.Called(prompt)
	return args.String(0), args.Error(1)
}
