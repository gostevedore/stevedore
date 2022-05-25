package dispatch

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
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

func (m *MockDispatch) Start(ctx context.Context, opts ...OptionsFunc) error {
	args := m.Called(ctx, opts)
	return args.Error(0)
}

// Enqueue provides a mock function with given fields: _a0, _a1
func (m *MockDispatch) Enqueue(job scheduler.Jobber) {
	m.Called(job)
}
