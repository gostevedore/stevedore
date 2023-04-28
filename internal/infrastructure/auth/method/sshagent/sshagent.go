package sshagent

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

type SSHAgentAuthMethod struct {
	GitSSHUser string `json:"git_ssh_user"`
}

// NewSSHAgentAuthMethod creates a new SSHAgentAuthMethod from the given credential
func NewSSHAgentAuthMethod() *SSHAgentAuthMethod {
	return &SSHAgentAuthMethod{}
}

// AuthMethodConstructor return a SSHAgentAuthMethod
func (a *SSHAgentAuthMethod) AuthMethodConstructor(credential *credentials.Credential) (repository.AuthMethodReader, error) {

	if credential == nil {
		return nil, nil
	}

	if !credential.AllowUseSSHAgent {
		return nil, nil
	}

	if credential.GitSSHUser != "" {
		a = &SSHAgentAuthMethod{
			GitSSHUser: credential.GitSSHUser,
		}

	}

	return a, nil
}

// Name returns the name of the authentication method
func (a *SSHAgentAuthMethod) Name() string {
	return credentials.SSHAgentAuthMethod
}
