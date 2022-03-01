package render

import (
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/stretchr/testify/mock"
)

// MockImageRender is a mock implementation of the ImageRender interface
type MockImageRender struct {
	mock.Mock
}

// NewMockImageRender creates a new mock image render
func NewMockImageRender() *MockImageRender {
	return &MockImageRender{}
}

// Render is a mock implementation of the Render method
func (m *MockImageRender) Render(name, version string, i *image.Image) error {
	args := m.Called(name, version, i)
	return args.Error(0)
}
