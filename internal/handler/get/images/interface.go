package images

import (
	"context"

	application "github.com/gostevedore/stevedore/internal/application/get/images"
)

// Applicationer is the service for build commands
type Applicationer interface {
	Run(ctx context.Context, options *application.Options, optionsFunc ...application.OptionsFunc) error
}
