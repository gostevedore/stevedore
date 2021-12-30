package builders

import (
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/stretchr/testify/mock"
)

// MockBuilders is a mock of Builders
type MockBuilders struct {
	mock.Mock
}

// NewMockBuilders return a new MockBuilders
func NewMockBuilders() *MockBuilders {
	return &MockBuilders{}
}

// Find a builder by name
func (b *MockBuilders) Find(name string) (*builder.Builder, error) {
	args := b.Mock.Called(name)
	return args.Get(0).(*builder.Builder), args.Error(1)
}
