package dispatch

import (
	"github.com/gostevedore/stevedore/internal/schedule"
)

// WorkerFactorier interface defines a worker factory
type WorkerFactorier interface {
	New(chan chan schedule.Jobber) schedule.Workerer
}
