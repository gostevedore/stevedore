package builders

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGetBuildersApplication is a mock of build application
type MockGetBuildersApplication struct {
	mock.Mock
}

// NewMockGetBuildersApplication return a mock of build application
func NewMockGetBuildersApplication() *MockGetBuildersApplication {
	return &MockGetBuildersApplication{}
}

// Run provides a mock function to carry out the application tasks
func (m *MockGetBuildersApplication) Run(ctx context.Context, options *Options, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, options, optionsFunc)
	return args.Error(0)
}
