package dispatch

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/schedule"
)

const DefaultNumWorkers = 1

type Dispatch struct {
	WorkerPool    chan chan schedule.Jobber
	inputJobQueue chan schedule.Jobber
	NumWorkers    int
	workerFactory WorkerFactorier
}

func NewDispatch(numWorkers int, workerFactory WorkerFactorier) *Dispatch {

	if numWorkers < 1 {
		numWorkers = DefaultNumWorkers
	}

	dispatch := &Dispatch{
		WorkerPool:    make(chan chan schedule.Jobber, numWorkers),
		inputJobQueue: make(chan schedule.Jobber),
		NumWorkers:    numWorkers,
		workerFactory: workerFactory,
	}

	return dispatch
}

// Start prepares dispatcher to start workers and dispatch jobs
func (d *Dispatch) Start(ctx context.Context) error {

	errContext := "(dispatch::Start)"

	if ctx == nil {
		return errors.New(errContext, "Dispatch requires a context to start")
	}

	if d.workerFactory == nil {
		return errors.New(errContext, "Dispatch requires a worker factory")
	}

	if d.WorkerPool == nil {
		return errors.New(errContext, "Dispatch requires a worker pool")
	}

	for i := 0; i < d.NumWorkers; i++ {
		worker := d.workerFactory.New(d.WorkerPool)

		go worker.Start(ctx)
	}

	go d.dispatch()

	return nil
}

func (d *Dispatch) dispatch() {
	for {
		j := <-d.inputJobQueue
		go func(j schedule.Jobber) {
			jobChannel := <-d.WorkerPool
			jobChannel <- j
		}(j)
	}
}

// Queue
func (d *Dispatch) Enqueue(job schedule.Jobber) {
	d.inputJobQueue <- job
}
