package promote

import (
	"context"

	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// HandlerPromoter is the interface that wraps the handler promote
type HandlerPromoter interface {
	Handler(ctx context.Context, options *handler.Options) error
}

// Entrypointer is the interface that wraps the main build function
type Entrypointer interface {
	Execute(ctx context.Context, args []string, conf *configuration.Configuration, entrypointOptions *entrypoint.Options, handlerOptions *handler.Options) error
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
