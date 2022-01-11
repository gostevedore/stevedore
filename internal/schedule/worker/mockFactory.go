package worker

import (
	"github.com/gostevedore/stevedore/internal/schedule"
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
func (f *MockWorkerFactory) New(workerPool chan chan schedule.Jobber) schedule.Workerer {
	args := f.Mock.Called(workerPool)

	return args.Get(0).(schedule.Workerer)
}
