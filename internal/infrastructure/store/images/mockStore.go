package images

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"

	"github.com/stretchr/testify/mock"
)

// MockStore is a mock implementation of the ImageStore interface
type MockStore struct {
	mock.Mock
}

// NewMockStore returns a new instance of the MockStore
func NewMockStore() *MockStore {
	return &MockStore{}
}

// Store is a mock implementation of the Store method
func (m *MockStore) Store(name string, version string, parent *image.Image) error {
	args := m.Called(name, version, parent)
	return args.Error(0)
}

// List is a mock implementation of the List method
func (m *MockStore) List() ([]*image.Image, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]*image.Image), args.Error(1)
	}
}

// FindByName is a mock implementation of the All method
func (m *MockStore) FindByName(name string) ([]*image.Image, error) {
	args := m.Called(name)
	return args.Get(0).([]*image.Image), args.Error(1)
}

// Find is a mock implementation of the Find method
func (m *MockStore) Find(name string, version string) ([]*image.Image, error) {
	args := m.Called(name, version)
	return args.Get(0).([]*image.Image), args.Error(1)
}

// FindGuaranteed is a mock implementation of the FindGuaranteed method
func (m *MockStore) FindGuaranteed(imageName, imageVersion string) ([]*image.Image, error) {
	args := m.Called(imageName, imageVersion)
	return args.Get(0).([]*image.Image), args.Error(1)
}
