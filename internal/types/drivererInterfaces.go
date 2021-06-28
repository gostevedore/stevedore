package types

import "context"

// Driverer element
type Driverer interface {
	Run(context.Context) error
}
