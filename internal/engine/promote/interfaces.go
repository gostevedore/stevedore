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
	GetPromoter(string) (promote.Promoter, error)
}

// Semverser
type Semverser interface {
	// GenerateSemvVer(version string) error
	// GenerateVersionTree(tmpl []string) ([]string, error)
	GenerateSemverList(version []string, tmpls []string) ([]string, error)
}