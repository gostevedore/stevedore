package schedule

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
)

// Jobber interface defines a job element
type Jobber interface {
	Run(ctx context.Context)
}

// Worker
type Worker struct {
	context    context.Context
	WorkerPool chan chan Jobber
	JobChannel chan Jobber
	quit       chan bool
}

// NewWorker
func NewWorker(ctx context.Context, workerPool chan chan Jobber) (*Worker, error) {
	if workerPool == nil {
		return nil, errors.New("(schedule::NewWorker)", "workerPool is nil")
	}

	if ctx == nil {
		return nil, errors.New("(schedule::NewWorker)", "context is nil")
	}

	worker := &Worker{
		context:    ctx,
		WorkerPool: workerPool,
		JobChannel: make(chan Jobber),
		quit:       make(chan bool),
	}

	return worker, nil
}

// Start
func (w *Worker) Start() error {

	if w.context == nil {
		return errors.New("(schedule::Worker::Start)", "context is nil")
	}

	for {
		w.WorkerPool <- w.JobChannel

		select {
		case job := <-w.JobChannel:
			job.Run(w.context)
		case <-w.quit:
			return nil
		case <-w.context.Done():
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
