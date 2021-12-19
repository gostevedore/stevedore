package job

import "context"

// BuildCommander interface defines the command to build a docker image
type BuildCommander interface {
	Execute(context.Context) error
}
