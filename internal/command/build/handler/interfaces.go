package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/engine/build"
)

// ServiceBuilder is the service for build commands
type ServiceBuilder interface {
	Build(ctx context.Context, name string, version []string, options *build.ServiceOptions) error
}
