package builders

import (
	"context"

	application "github.com/gostevedore/stevedore/internal/application/get/builders"
)

// Applicationer is the service for build commands
type Applicationer interface {
	Run(ctx context.Context, options *application.Options, optionsFunc ...application.OptionsFunc) error
}
