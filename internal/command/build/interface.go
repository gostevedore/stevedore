package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/configuration"
	buildentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	buildhandler "github.com/gostevedore/stevedore/internal/handler/build"
)

// Entrypointer is the interface that wraps the main build function
type Entrypointer interface {
	Execute(ctx context.Context, args []string, conf *configuration.Configuration, entrypointOptions *buildentrypoint.EntrypointOptions, handlerOptions *buildhandler.HandlerOptions) error
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
