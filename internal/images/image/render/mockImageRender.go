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
func (m *MockImageRender) Render(name, version string, parent *image.Image, img ImageSerializer) error {
	args := m.Called(name, version, parent, img)
	return args.Error(0)
}
