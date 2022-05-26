package promote

import (
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// Outputter
type Outputter interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
}

// PromoteFactorier
type PromoteFactorier interface {
	Get(string) (repository.Promoter, error)
	Register(id string, promoter repository.Promoter) error
}

// Semverser
type Semverser interface {
	GenerateSemverList(version []string, tmpls []string) ([]string, error)
}
