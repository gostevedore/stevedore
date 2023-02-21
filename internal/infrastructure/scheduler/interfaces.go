package scheduler

import "context"

// Workerer interface defines a worker
type Workerer interface {
	Start(ctx context.Context) error
	Stop()
}

// Jobber interface defines a job element
type Jobber interface {
	Run(context.Context)
	Wait() error
	Done() <-chan struct{}
	Err() <-chan error
	Close()
}
