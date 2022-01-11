package worker

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockWorker is a mock of Worker
type MockWorker struct {
	mock.Mock
}

// NewMockWorker provides a mock function with given fields: workerPool
func NewMockWorker() *MockWorker {
	return &MockWorker{}
}

// Start provides a mock function to be called
func (m *MockWorker) Start(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

// Stop provides a mock function to be called
func (m *MockWorker) Stop() {
	m.Called()
}
