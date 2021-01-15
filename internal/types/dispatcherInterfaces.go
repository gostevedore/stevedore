package types

import "github.com/gostevedore/stevedore/internal/schedule"

// Dispatcher
type Dispatcher interface {
	Start() error
	Enqueue(schedule.Jobber)
}
