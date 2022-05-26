package graph

import (
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/image"
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

// Iterate is a mock implementation of the Iterate method
func (m *MockImagesGraphTemplate) Iterate() <-chan GraphNoder {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(<-chan GraphNoder)
}
