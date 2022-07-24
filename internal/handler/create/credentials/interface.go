package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/application/create/credentials"
)

// Applicationer is the service for build commands
type Applicationer interface {
	Run(ctx context.Context, optionsFunc ...credentials.OptionsFunc) error
}
