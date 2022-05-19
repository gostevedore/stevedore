package build

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockService is a mock of build service
type MockService struct {
	mock.Mock
}

// NewMockService return a mock of build service
func NewMockService() *MockService {
	return &MockService{}
}

// Build provides a mock function with given fields: ctx, buildPlan, name, version, options, optionsFunc
func (m *MockService) Build(ctx context.Context, buildPlan Planner, name string, version []string, options *ServiceOptions, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, buildPlan, name, version, options, optionsFunc)
	return args.Error(0)
}
