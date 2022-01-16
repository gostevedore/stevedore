package graph

import (
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	"github.com/stretchr/testify/mock"
)

// MockImagesGraphTemplate is a graph template for images
type MockImagesGraphTemplate struct {
	mock.Mock
}

func NewMockImagesGraphTemplate() *MockImagesGraphTemplate {
	return &MockImagesGraphTemplate{}
}

// AddImage is a mock implementation of the AddImage method
func (m *MockImagesGraphTemplate) AddImage(name, version string, image *image.Image) error {
	args := m.Called(name, version, image)
	return args.Error(0)
}
