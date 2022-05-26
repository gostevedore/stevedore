package worker

import (
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/stretchr/testify/mock"
)

// MockWorkerFactory is a factory for creating workers
type MockWorkerFactory struct {
	mock.Mock
}

// NewMockWorkerFactory returns a new worker factory
func NewMockWorkerFactory() *MockWorkerFactory {
	return &MockWorkerFactory{}
}

// New returns a new worker constructor
func (f *MockWorkerFactory) New(workerPool chan chan scheduler.Jobber) scheduler.Workerer {
	args := f.Mock.Called(workerPool)

	return args.Get(0).(scheduler.Workerer)
}
