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

// NewKeyFileAuthMethod creates a new KeyFileAuthMethod from the given credential
func NewKeyFileAuthMethod() *KeyFileAuthMethod {
	return &KeyFileAuthMethod{}
}

// AuthMethodConstructor creates a new KeyFileAuthMethod from the given credential
func (auth *KeyFileAuthMethod) AuthMethodConstructor(credential *credentials.Credential) (repository.AuthMethodReader, error) {

	if credential == nil {
		return nil, nil
	}

	if credential.PrivateKeyFile != "" {
		auth = &KeyFileAuthMethod{
			PrivateKeyFile: credential.PrivateKeyFile,
		}

		if credential.PrivateKeyPassword != "" {
			auth.PrivateKeyPassword = credential.PrivateKeyPassword
		}

		if credential.GitSSHUser != "" {
			auth.GitSSHUser = credential.GitSSHUser
		}

		return auth, nil
	} else {
		return nil, nil
	}
}

// Name returns the name of the authentication method
func (a *KeyFileAuthMethod) Name() string {
	return credentials.KeyFileAuthMethod
}
