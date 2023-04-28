package credentials

import (
	"context"

	application "github.com/gostevedore/stevedore/internal/application/create/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

// Applicationer is the service for build commands
type Applicationer interface {
	Run(ctx context.Context, id string, credential *credentials.Credential, optionsFunc ...application.OptionsFunc) error
}
