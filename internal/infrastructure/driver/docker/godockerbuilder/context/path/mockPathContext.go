package pathcontext

import (
	"github.com/stretchr/testify/mock"
)

// MockPathBuildContext defines a docker build mock path context
type MockPathBuildContext struct {
	mock.Mock
}

// NewMockPathBuildContext provides a mock function with given fields
func NewMockPathBuildContext(path string) *MockPathBuildContext {
	return &MockPathBuildContext{}
}

// WithPath provides a mock function to set the path value
func (c *MockPathBuildContext) WithPath(path string) {
	c.Called(path)
}
