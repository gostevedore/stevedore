package dispatch

import (
	"context"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/schedule"
)

// DefaultNumWorkers is the default number of workers
const DefaultNumWorkers = 1

// OptionsFunc is a function used to configure the dispatcher
type OptionsFunc func(*Dispatch)

// Dispatch is a dispatcher that executes jobs
type Dispatch struct {
	WorkerPool    chan chan schedule.Jobber
	inputJobQueue chan schedule.Jobber
	NumWorkers    int
	workerFactory WorkerFactorier
	once          sync.Once
}

// New creates a new dispatcher
func NewDispatch(workerFactory WorkerFactorier, options ...OptionsFunc) *Dispatch {

	dispatch := &Dispatch{
		WorkerPool:    make(chan chan schedule.Jobber, DefaultNumWorkers),
		inputJobQueue: make(chan schedule.Jobber),
		workerFactory: workerFactory,
	}

	dispatch.Options(options...)

	return dispatch
}

// Options configure the stevedore command
func (d *Dispatch) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(d)
	}
}

func WithNumWorkers(n int) OptionsFunc {
	return func(d *Dispatch) {
		if n < 1 {
			d.NumWorkers = DefaultNumWorkers
		} else {
			d.NumWorkers = n
		}
		d.WorkerPool = make(chan chan schedule.Jobber, d.NumWorkers)
	}
}

// Start prepares dispatcher to start workers and dispatch jobs
func (d *Dispatch) Start(ctx context.Context, opts ...OptionsFunc) error {

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

	d.Options(opts...)

	if d.NumWorkers < 1 {
		d.NumWorkers = DefaultNumWorkers
	}

	d.once.Do(func() {
		for i := 0; i < d.NumWorkers; i++ {
			worker := d.workerFactory.New(d.WorkerPool)

			go worker.Start(ctx)
		}

		go d.dispatch()
	})

	return nil
}

// dispatch is the main loop of the dispatcher
func (d *Dispatch) dispatch() {

	for {
		j := <-d.inputJobQueue
		go func(j schedule.Jobber) {
			jobChannel := <-d.WorkerPool
			jobChannel <- j
		}(j)
	}
}

// Enqueue enqueues a job to be executed by a worker
func (d *Dispatch) Enqueue(job schedule.Jobber) {
	d.inputJobQueue <- job
}
