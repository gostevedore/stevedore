package build

import (
	"context"

	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// Entrypointer is the interface that wraps the main build function
type Entrypointer interface {
	Execute(ctx context.Context, args []string, conf *configuration.Configuration, compatibility entrypoint.Compatibilitier, inputEntrypointOptions *entrypoint.Options, inputHandlerOptions *handler.Options) error
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
