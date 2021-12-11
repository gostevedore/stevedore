package driver

import "context"

// BuildDriverer
type BuildDriverer interface {
	//Run(context.Context) error
	Build(context.Context, BuildDriverOptions) error
}
