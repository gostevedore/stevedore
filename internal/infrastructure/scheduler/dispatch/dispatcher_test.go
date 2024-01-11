package dispatch

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/worker"
	"github.com/stretchr/testify/assert"
)

func TestNewDispatch(t *testing.T) {
	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	tests := []struct {
		desc          string
		context       context.Context
		workerFactory WorkerFactorier
		numWorkers    int
		resWorkers    int
	}{
		{
			desc:       "Testing a new dispatcher creation with an invalid number of workers",
			context:    cancelContext,
			numWorkers: 0,
			resWorkers: 1,
		},
		{
			desc:       "Testing a new dispatcher creation",
			context:    cancelContext,
			numWorkers: 1,
			resWorkers: 1,
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		d := NewDispatch(test.workerFactory, WithNumWorkers(test.numWorkers))
		assert.NotNil(t, d.WorkerPool)
		assert.NotNil(t, d.inputJobQueue)
		assert.Equal(t, test.resWorkers, d.NumWorkers)

	}

}

func TestStart(t *testing.T) {

	errContext := "(dispatch::Start)"

	tests := []struct {
		desc              string
		err               error
		dispatch          *Dispatch
		numWorkers        int
		context           context.Context
		prepareAssertFunc func(*Dispatch, chan chan scheduler.Jobber)
	}{
		{
			desc:     "Testing error when starting a dispatcher without a context",
			dispatch: &Dispatch{},
			context:  nil,
			err:      errors.New(errContext, "Dispatch requires a context to start"),
		},
		{
			desc:     "Testing error when starting a dispatcher without a worker factory",
			dispatch: &Dispatch{},
			context:  context.TODO(),
			err:      errors.New(errContext, "Dispatch requires a worker factory"),
		},
		{
			desc: "Testing error when starting a dispatcher without a worker pool",
			dispatch: &Dispatch{
				workerFactory: worker.NewMockWorkerFactory(),
			},
			context: context.TODO(),
			err:     errors.New(errContext, "Dispatch requires a worker pool"),
		},
		// Test is commented because the workers stated is not being checked outside the goroutine
		// {
		// 	desc:    "Testing error when starting a new dispatcher that workers failed to start",
		// 	err:     &errors.Error{},
		// 	context: context.TODO(),
		// 	dispatch: &Dispatch{
		// 		WorkerPool:    make(chan chan scheduler.Jobber, 1),
		// 		inputJobQueue: make(chan scheduler.Jobber),
		// 		NumWorkers:    1,
		// 		workerFactory: worker.NewMockWorkerFactory(),
		// 	},
		// 	numWorkers: 1,
		// 	prepareAssertFunc: func(d *Dispatch, pool chan chan scheduler.Jobber) {
		// 		w := worker.NewMockWorker()
		// 		w.On("Start", context.TODO()).Return(errors.New(errContext, "error"))

		// 		d.workerFactory.(*worker.MockWorkerFactory).On("New", pool).Return(w)
		// 	},
		// },
		{
			desc:    "Testing start a new dispatcher with 5 workers",
			err:     &errors.Error{},
			context: context.TODO(),
			dispatch: &Dispatch{
				WorkerPool:    make(chan chan scheduler.Jobber, 5),
				inputJobQueue: make(chan scheduler.Jobber),
				NumWorkers:    5,
				workerFactory: worker.NewMockWorkerFactory(),
			},
			numWorkers: 5,
			prepareAssertFunc: func(d *Dispatch, pool chan chan scheduler.Jobber) {
				d.workerFactory.(*worker.MockWorkerFactory).On("New", pool).Return(worker.NewWorker(pool))
			},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.dispatch, test.dispatch.WorkerPool)
			}

			err := test.dispatch.Start(test.context)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.dispatch.workerFactory.(*worker.MockWorkerFactory).AssertNumberOfCalls(t, "New", test.numWorkers)
			}
		})

	}
}

// func TestEnqueue(t *testing.T) {
// 	backgroundContext := context.Background()
// 	cancelContext, cancel := context.WithCancel(backgroundContext)
// 	defer cancel()

// 	dispatch := &Dispatch{
// 		context:       cancelContext,
// 		WorkerPool:    make(chan chan Jobber, 1),
// 		inputJobQueue: make(chan Jobber),
// 		NumWorkers:    1,
// 	}

// 	go dispatch.Start()

// 	job := &MockJobber{}
// 	dispatch.Enqueue(job)
// 	// sleep to wait a grace time for processing job
// 	time.Sleep(1 * time.Millisecond)
// 	assert.True(t, job.run)

// }
