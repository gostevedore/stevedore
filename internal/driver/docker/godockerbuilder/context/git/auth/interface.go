package gitauth

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/gostevedore/stevedore/internal/credentials"
)

// GitAuther is an interface for git authentication
type GitAuther interface {
	Auth() (transport.AuthMethod, error)
}

// CredentialsStorer
type CredentialsStorer interface {
	GetCredentials(registy string) (*credentials.RegistryUserPassAuth, error)
}
