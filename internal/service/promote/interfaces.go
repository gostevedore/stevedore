package promote

import (
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/promote"
)

// CredentialsStorer
type CredentialsStorer interface {
	GetCredentials(registy string) (*credentials.RegistryUserPassAuth, error)
}

// Outputter
type Outputter interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
}

// PromoteFactorier
type PromoteFactorier interface {
	Get(string) (promote.Promoter, error)
}

// Semverser
type Semverser interface {
	GenerateSemverList(version []string, tmpls []string) ([]string, error)
}
