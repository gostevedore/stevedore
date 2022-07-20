package sshagent

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

const (
	// SSHAgentType is the name of the basic authentication method
	SSHAgentType = "SSH agent"
)

type SSHAgentOutput struct{}

func NewSSHAgentOutput() *SSHAgentOutput {
	return &SSHAgentOutput{}
}

func (o *SSHAgentOutput) Output(badge *credentials.Badge) (string, string, error) {

	errContext := "(output::credentials::types::SSHAgentOutput::Output)"

	if badge == nil {
		return "", "", errors.New(errContext, "To show badge output, badge must be provided")
	}

	if badge.AllowUseSSHAgent {
		detail := "Use SSH agent"

		if badge.GitSSHUser != "" {
			detail = fmt.Sprintf("%s, with git user '%s'", detail, badge.GitSSHUser)
		}

		return SSHAgentType, detail, nil
	} else {
		return "", "", nil
	}
}
