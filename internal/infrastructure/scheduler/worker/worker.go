package worker

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
)

// Worker defines a worker
type Worker struct {
	WorkerPool chan chan scheduler.Jobber
	JobChannel chan scheduler.Jobber
	quit       chan bool
}

// NewWorker creates a new worker
func NewWorker(workerPool chan chan scheduler.Jobber) *Worker {

	worker := &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan scheduler.Jobber),
		quit:       make(chan bool),
	}

	return worker
}

// Start initiates the worker routine
func (w *Worker) Start(ctx context.Context) error {

	errContext := "(worker::Start)"

	if ctx == nil {
		return errors.New(errContext, "Worker requires a context to start")
	}

	if w.WorkerPool == nil {
		return errors.New(errContext, "Worker requires a pool to register jobs")
	}

	for {
		w.WorkerPool <- w.JobChannel

		select {
		case job := <-w.JobChannel:
			job.Run(ctx)
		case <-w.quit:
			return nil
		case <-ctx.Done():
			w.Stop()
		}
	}
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
