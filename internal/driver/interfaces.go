package driver

import "context"

// Driverer
type Driverer interface {
	Run(context.Context) error
}
