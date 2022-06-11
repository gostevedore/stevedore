package keyfile

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// KeyFileAuthMethod is used for private key file authentication
type KeyFileAuthMethod struct {
	// PrivateKeyFile is the path to the private key file
	PrivateKeyFile string `json:"private_key_file"`
	// PrivateKeyPassword is the password for the private key file
	PrivateKeyPassword string `json:"private_key_password"`
	// GitSSHUser is the user to use for git ssh
	GitSSHUser string `json:"git_ssh_user"`
}

// NewKeyFileAuthMethod creates a new KeyFileAuthMethod from the given badge
func NewKeyFileAuthMethod() *KeyFileAuthMethod {
	return &KeyFileAuthMethod{}
}

// AuthMethod creates a new KeyFileAuthMethod from the given badge
func (a *KeyFileAuthMethod) AuthMethod(badge *credentials.Badge) (repository.AuthMethodReader, error) {
	if badge == nil {
		return nil, nil
	}

	if badge.PrivateKeyFile != "" {
		a = &KeyFileAuthMethod{
			PrivateKeyFile: badge.PrivateKeyFile,
		}

		if badge.PrivateKeyPassword == "" {
			a.PrivateKeyPassword = badge.PrivateKeyPassword
		}

		if badge.GitSSHUser == "" {
			a.GitSSHUser = badge.GitSSHUser
		}
	}

	return a, nil
}

// Name returns the name of the authentication method
func (a *KeyFileAuthMethod) Name() string {
	return credentials.KeyFileAuthMethod
}
