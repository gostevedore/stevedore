package gitauth

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// GitAuther is an interface for git authentication
type GitAuther interface {
	Auth() (transport.AuthMethod, error)
}
