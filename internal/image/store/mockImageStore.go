package store

import (
	"github.com/gostevedore/stevedore/internal/image"

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

// All is a mock implementation of the All method
func (m *MockImageStore) All(name string) ([]*image.Image, error) {
	args := m.Called(name)
	return args.Get(0).([]*image.Image), args.Error(1)
}

// Find is a mock implementation of the Find method
func (m *MockImageStore) Find(name string, version string) (*image.Image, error) {
	args := m.Called(name, version)
	return args.Get(0).(*image.Image), args.Error(1)
}
