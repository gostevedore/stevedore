package build

import (
	"context"

	handler "github.com/gostevedore/stevedore/internal/command/build/handler"
)

// Handlerer is a handler for build commands
type Handlerer interface {
	Handler(ctx context.Context, options *handler.HandlerOptions) error
}