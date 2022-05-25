package worker

import (
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
)

// WorkerFactory is a factory for creating workers
type WorkerFactory struct{}

// NewWorkerFactory returns a new worker factory
func NewWorkerFactory() *WorkerFactory {
	return &WorkerFactory{}
}

// New returns a new worker constructor
func (f *WorkerFactory) New(workerPool chan chan scheduler.Jobber) scheduler.Workerer {
	return NewWorker(workerPool)
}
