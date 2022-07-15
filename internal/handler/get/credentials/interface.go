package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/application/get/credentials"
)

// GetCredentialsApplication is the service for build commands
type GetCredentialsApplication interface {
	Run(ctx context.Context, optionsFunc ...credentials.OptionsFunc) error
}
