package images

import (
	"context"

	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/images"
	handler "github.com/gostevedore/stevedore/internal/handler/get/images"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// Entrypointer is the interface that wraps the main function
type Entrypointer interface {
	Execute(ctx context.Context, args []string, conf *configuration.Configuration, entrypointOptions *entrypoint.Options, handlerOptions *handler.Options) error
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
