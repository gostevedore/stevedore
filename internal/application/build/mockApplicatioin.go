package build

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
func (m *MockApplication) Build(ctx context.Context, buildPlan Planner, name string, version []string, options *Options, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, buildPlan, name, version, options, optionsFunc)
	return args.Error(0)
}
