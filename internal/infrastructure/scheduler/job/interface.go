package job

import "context"

// Commander interface defines the command to be executed
type Commander interface {
	Execute(context.Context) error
}
