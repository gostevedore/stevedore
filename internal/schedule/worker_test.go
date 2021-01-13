package schedule

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

// TestNewWorker
func TestNewWorker(t *testing.T) {

	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	tests := []struct {
		desc       string
		err        error
		workerPool chan chan Jobber
		context    context.Context
	}{
		{
			desc:       "Testing create new worker with nil workerPool",
			err:        errors.New("(schedule::NewWorker)", "workerPool is nil"),
			workerPool: nil,
			context:    context.Background(),
		},
		{
			desc:       "Testing create new worker with nil context",
			err:        errors.New("(schedule::NewWorker)", "context is nil"),
			workerPool: make(chan chan Jobber),
			context:    nil,
		},
		{
			desc:       "Testing create new worker",
			err:        nil,
			workerPool: make(chan chan Jobber),
			context:    cancelContext,
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		worker, err := NewWorker(test.context, test.workerPool)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.NotNil(t, worker.JobChannel)
			assert.NotNil(t, worker.WorkerPool)
			assert.NotNil(t, worker.quit)
		}
	}

}

// TestStart
func TestStartWorker(t *testing.T) {

	tests := []struct {
		desc     string
		err      error
		worker   *Worker
		preFunc  func(w *Worker)
		postFunc func(w *Worker)
	}{
		{
			desc: "Testing start a worker without context",
			err:  errors.New("(schedule::Worker::Start)", "context is nil"),
			worker: &Worker{
				WorkerPool: make(chan chan Jobber),
				JobChannel: make(chan Jobber),
				quit:       make(chan bool),
				context:    nil,
			},
			preFunc:  func(w *Worker) {},
			postFunc: func(w *Worker) {},
		},
	}

	errChan := make(chan error)
	for _, test := range tests {

		t.Log(test.desc)

		if test.preFunc != nil {
			go test.preFunc(test.worker)
		}

		go func() {
			err := test.worker.Start()
			errChan <- err
		}()

		err := <-errChan
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}

		if test.preFunc != nil {
			go test.postFunc(test.worker)
		}
	}
	close(errChan)
}

func TestRunJobOnAWorker(t *testing.T) {
	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	t.Log("Testing run job on a worker")

	worker := &Worker{
		WorkerPool: make(chan chan Jobber),
		JobChannel: make(chan Jobber),
		quit:       make(chan bool),
		context:    cancelContext,
	}
	job := &MockJobber{}

	// Pretask
	// it will let to register the worker on the worker pool without blocking it
	go func() {
		<-worker.WorkerPool
	}()

	go func() {
		worker.Start()
	}()

	worker.JobChannel <- job
	go func() {
		worker.Stop()
	}()

	assert.True(t, job.run)

}

// TestStop
func TestStopWorker(t *testing.T) {
	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	tests := []struct {
		desc    string
		err     error
		worker  *Worker
		preFunc func(w *Worker)
	}{
		{
			desc: "Testing stop worker",
			err:  nil,
			worker: &Worker{
				WorkerPool: make(chan chan Jobber),
				JobChannel: make(chan Jobber),
				quit:       make(chan bool),
				context:    cancelContext,
			},
			preFunc: func(w *Worker) {
				w.Start()
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		go test.preFunc(test.worker)
		go test.worker.Stop()
	}

}
