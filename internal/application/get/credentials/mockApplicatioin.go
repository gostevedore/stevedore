package credentials

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockApplication is a mock of build application
type MockApplication struct {
	mock.Mock
}

// NewMockApplication return a mock of build application
func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

// Build provides a mock function with given fields: ctx, buildPlan, name, version, options, optionsFunc
func (m *MockApplication) Build(ctx context.Context, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, optionsFunc)
	return args.Error(0)
}
