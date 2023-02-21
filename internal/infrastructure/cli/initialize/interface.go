package initialize

import (
	"context"

	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/configuration"
)

// Entrypointer is the interface that wraps the main function
type Entrypointer interface {
	Execute(ctx context.Context, options *entrypoint.Options) error
}
