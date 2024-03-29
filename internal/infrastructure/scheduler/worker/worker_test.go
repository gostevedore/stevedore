package worker

import (
	"context"
	"testing"
	"time"

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

		assert.NotNil(t, worker)
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

		workerPool := make(chan chan scheduler.Jobber, 1)
		worker := NewWorker(workerPool)
		testJob := job.NewMockJob()
		testJob.Mock.On("Run", context.TODO()).Return(nil)

		defer worker.Stop()
		go func() {
			err := worker.Start(context.TODO())
			assert.Nil(t, err)
		}()

		workerPool <- worker.JobChannel
		jobChannel := <-workerPool
		jobChannel <- testJob

		time.Sleep(200 * time.Millisecond)

		assert.True(t, testJob.Mock.AssertExpectations(t))

	})
}
