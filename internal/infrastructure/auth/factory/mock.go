package factory

import (
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/mock"
)

// MockAuthFactory is a mock of AuthFactory interface
type MockAuthFactory struct {
	mock.Mock
}

// NewMockAuthFactory creates a new auth provider factory
func NewMockAuthFactory() *MockAuthFactory {
	return &MockAuthFactory{}
}

// Get provides a mock function with given fields: id
func (f *MockAuthFactory) Get(id string) (repository.AuthMethodReader, error) {
	args := f.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(repository.AuthMethodReader), args.Error(1)
}
