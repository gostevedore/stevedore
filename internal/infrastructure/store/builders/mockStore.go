package builders

import (
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/stretchr/testify/mock"
)

// MockStore is a mock of Builders
type MockStore struct {
	mock.Mock
}

// NewMockStore return a new MockStore
func NewMockStore() *MockStore {
	return &MockStore{}
}

// AddBuilder add a builder
func (b *MockStore) Store(builder *builder.Builder) error {
	args := b.Called(builder)
	return args.Error(0)
}

// Find a builder by name
func (b *MockStore) Find(name string) (*builder.Builder, error) {
	args := b.Mock.Called(name)
	return args.Get(0).(*builder.Builder), args.Error(1)
}
