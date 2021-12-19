package driver

import "context"

// BuildDriverer interface defines which methods are used to build a docker image
type BuildDriverer interface {
	//Run(context.Context) error
	Build(context.Context, BuildDriverOptions) error
}
