package schedule

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
)

type Dispatch struct {
	context       context.Context
	WorkerPool    chan chan Jobber
	inputJobQueue chan Jobber
	NumWorkers    int
}

func NewDispatch(ctx context.Context, numWorkers int) (*Dispatch, error) {

	if numWorkers < 1 {
		return nil, errors.New("(schedule::NewDispatch)", "Invalid value for number of workers, it must be greater than zero")
	}

	dispatch := &Dispatch{
		context:       ctx,
		WorkerPool:    make(chan chan Jobber, numWorkers),
		inputJobQueue: make(chan Jobber),
		NumWorkers:    numWorkers,
	}

	return dispatch, nil
}

// Start
func (d *Dispatch) Start() error {

	for i := 0; i < d.NumWorkers; i++ {
		worker, err := NewWorker(d.context, d.WorkerPool)
		if err != nil {
			return errors.New("(schedule::Dispatch::Start)", "Worker could not be created", err)
		}

		go worker.Start()
	}

	go d.dispatch()

	return nil
}

func (d *Dispatch) dispatch() {
	for {
		select {
		case job := <-d.inputJobQueue:
			go func(job Jobber) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		case <-d.context.Done():
			// TODO: run a graceful stop to all
			return
		}
	}
}

// Queue
func (d *Dispatch) Enqueue(job Jobber) {
	d.inputJobQueue <- job
}
