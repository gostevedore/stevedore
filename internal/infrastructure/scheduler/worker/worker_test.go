package worker

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"github.com/stretchr/testify/assert"
)

// TestNewWorker
func TestNewWorker(t *testing.T) {

	tests := []struct {
		desc       string
		workerPool chan chan scheduler.Jobber
	}{
		{
			desc:       "Testing create new worker",
			workerPool: make(chan chan scheduler.Jobber),
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		worker := NewWorker(test.workerPool)

		assert.NotNil(t, worker.JobChannel)
		assert.NotNil(t, worker.WorkerPool)
		assert.NotNil(t, worker.quit)
	}

}

// TestStart
func TestStartWorker(t *testing.T) {

	errContext := "(worker::Start)"

	tests := []struct {
		desc    string
		err     error
		worker  *Worker
		context context.Context
	}{
		{
			desc:    "Testing error starting a worker without context",
			err:     errors.New(errContext, "Worker requires a context to start"),
			context: nil,
			worker:  &Worker{},
		},
		{
			desc:    "Testing  error starting a worker without worker pool",
			err:     errors.New(errContext, "Worker requires a pool to register jobs"),
			context: context.TODO(),
			worker:  &Worker{},
		},
	}

	errChan := make(chan error)
	_ = errChan
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.worker.Start(test.context)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.worker.Stop()
			}

		})

	}
}

func TestRunJobOnAWorker(t *testing.T) {
	desc := "Testing run a job on a worker"

	t.Run(desc, func(t *testing.T) {
		t.Log(desc)

		worker := &Worker{
			WorkerPool: make(chan chan scheduler.Jobber),
			JobChannel: make(chan scheduler.Jobber),
			quit:       make(chan bool),
		}
		testJob := job.NewMockJob()
		testJob.Mock.On("Run", context.TODO()).Return(nil)

		go func() {
			// Pretask: it reads the worker jobchannel regsitration and voids to block the worker
			<-worker.WorkerPool
		}()

		go func() {
			worker.Start(context.TODO())
		}()
		defer worker.Stop()

		worker.JobChannel <- testJob

		assert.True(t, testJob.Mock.AssertExpectations(t))

	})
}
