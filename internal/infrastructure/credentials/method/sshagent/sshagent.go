package sshagent

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

type SSHAgentAuthMethod struct {
	GitSSHUser string `json:"git_ssh_user"`
}

// NewSSHAgentAuthMethod creates a new SSHAgentAuthMethod from the given badge
func NewSSHAgentAuthMethod() *SSHAgentAuthMethod {
	return &SSHAgentAuthMethod{}
}

func (a *SSHAgentAuthMethod) AuthMethod(badge *credentials.Badge) (repository.AuthMethodReader, error) {

	if badge == nil {
		return nil, nil
	}

	if badge.GitSSHUser != "" {
		a = &SSHAgentAuthMethod{
			GitSSHUser: badge.GitSSHUser,
		}
	}

	return a, nil
}

// Name returns the name of the authentication method
func (a *SSHAgentAuthMethod) Name() string {
	return credentials.SSHAgentAuthMethod
}
