package builders

import (
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/stretchr/testify/mock"
)

// MockBuildersStore is a mock of Builders
type MockBuildersStore struct {
	mock.Mock
}

// NewMockBuildersStore return a new MockBuildersStore
func NewMockBuildersStore() *MockBuildersStore {
	return &MockBuildersStore{}
}

// AddBuilder add a builder
func (b *MockBuildersStore) Store(builder *builder.Builder) error {
	args := b.Called(builder)
	return args.Error(0)
}

// Find a builder by name
func (b *MockBuildersStore) Find(name string) (*builder.Builder, error) {
	args := b.Mock.Called(name)
	return args.Get(0).(*builder.Builder), args.Error(1)
}
