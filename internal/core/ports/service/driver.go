package service

import (
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// DriverFactorier interface defines the factory to create a build driver
type DriverFactorier interface {
	Get(id string) (repository.BuildDriverer, error)
	Register(id string, driver repository.BuildDriverer) error
}
