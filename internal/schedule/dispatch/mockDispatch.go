package dispatch

import (
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/stretchr/testify/mock"
)

// MockDispatch is a mock of Dispatch interface
type MockDispatch struct {
	mock.Mock
}

// NewMockDispatch provides a mock of Dispatch interface
func NewMockDispatch() *MockDispatch {
	return &MockDispatch{}
}

// Enqueue provides a mock function with given fields: _a0, _a1
func (m *MockDispatch) Enqueue(job schedule.Jobber) {
	m.Called(job)
}
