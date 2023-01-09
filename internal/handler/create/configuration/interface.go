package configuration

import (
	"context"

	application "github.com/gostevedore/stevedore/internal/application/create/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// Applicationer is the service for createConfiguration commands
type Applicationer interface {
	Run(ctx context.Context, config *configuration.Configuration, optionsFunc ...application.OptionsFunc) error
}
