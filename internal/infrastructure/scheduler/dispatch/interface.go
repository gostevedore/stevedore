package dispatch

import "github.com/gostevedore/stevedore/internal/infrastructure/scheduler"

// WorkerFactorier interface defines a worker factory
type WorkerFactorier interface {
	New(chan chan scheduler.Jobber) scheduler.Workerer
}
