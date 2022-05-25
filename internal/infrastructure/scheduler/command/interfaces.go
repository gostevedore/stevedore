package command

import "context"

// // BuildDriverer interface defines which methods are used to build a docker image
// type BuildDriverer interface {
// 	Build(context.Context, *driver.BuildDriverOptions) error
// }

// BuildCommander interface defines the command to build a docker image
type BuildCommander interface {
	Execute(context.Context) error
}
