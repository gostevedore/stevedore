package plan

import "github.com/stretchr/testify/mock"

// MockPlan is a mock of Plan interface
type MockPlan struct {
	mock.Mock
}

// NewMockPlan returns a new MockPlan
func NewMockPlan() *MockPlan {
	return &MockPlan{}
}

// Plan mock
func (p *MockPlan) Plan(name string, version []string) ([]*Step, error) {
	args := p.Called(name, version)
	return args.Get(0).([]*Step), args.Error(1)
}
