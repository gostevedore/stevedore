package credentials

import (
	"context"

	application "github.com/gostevedore/stevedore/internal/application/get/credentials"
)

// Applicationer is the service for build commands
type Applicationer interface {
	Run(ctx context.Context, optionsFunc ...application.OptionsFunc) error
}
