package plan

import "github.com/stretchr/testify/mock"

// MockPlanFactory is a factory to create Planner
type MockPlanFactory struct {
	mock.Mock
}

// NewMockPlanFactory creates a new MockPlanFactory
func NewMockPlanFactory() *MockPlanFactory {
	return &MockPlanFactory{}
}

// NewPlan provides a mock function with given fields: id, parameters
func (f *MockPlanFactory) NewPlan(id string, parameters map[string]interface{}) (Planner, error) {
	args := f.Called(id, parameters)
	return args.Get(0).(Planner), args.Error(1)
}
