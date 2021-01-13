package types

import (
	"stevedore/internal/schedule"
)

// Dispatcher
type Dispatcher interface {
	Start() error
	Enqueue(schedule.Jobber)
}
