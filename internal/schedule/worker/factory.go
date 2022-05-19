package worker

import "github.com/gostevedore/stevedore/internal/schedule"

// WorkerFactory is a factory for creating workers
type WorkerFactory struct{}

// NewWorkerFactory returns a new worker factory
func NewWorkerFactory() *WorkerFactory {
	return &WorkerFactory{}
}

// New returns a new worker constructor
func (f *WorkerFactory) New(workerPool chan chan schedule.Jobber) schedule.Workerer {
	return NewWorker(workerPool)
}
