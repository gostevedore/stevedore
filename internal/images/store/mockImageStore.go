package store

import (
	"github.com/gostevedore/stevedore/internal/images/image"

	"github.com/stretchr/testify/mock"
)

// MockImageStore is a mock implementation of the ImageStore interface
type MockImageStore struct {
	mock.Mock
}

// NewMockImageStore returns a new instance of the MockImageStore
func NewMockImageStore() *MockImageStore {
	return &MockImageStore{}
}

// Store is a mock implementation of the Store method
func (m *MockImageStore) Store(name string, version string, parent *image.Image) error {
	args := m.Called(name, version, parent)
	return args.Error(0)
}

// List is a mock implementation of the List method
func (m *MockImageStore) List() ([]*image.Image, error) {
	args := m.Called()
	return args.Get(0).([]*image.Image), args.Error(1)
}

// FindByName is a mock implementation of the All method
func (m *MockImageStore) FindByName(name string) ([]*image.Image, error) {
	args := m.Called(name)
	return args.Get(0).([]*image.Image), args.Error(1)
}

// Find is a mock implementation of the Find method
func (m *MockImageStore) Find(name string, version string) (*image.Image, error) {
	args := m.Called(name, version)
	return args.Get(0).(*image.Image), args.Error(1)
}

// FindGuaranteed is a mock implementation of the FindGuaranteed method
func (m *MockImageStore) FindGuaranteed(findName, findVersion, imageName, imageVersion string) (*image.Image, error) {
	args := m.Called(findName, findVersion, imageName, imageVersion)
	return args.Get(0).(*image.Image), args.Error(1)
}
