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

func (o *SSHAgentOutput) Output(credential *credentials.Credential) (string, string, error) {

	errContext := "(output::credentials::types::SSHAgentOutput::Output)"

	if credential == nil {
		return "", "", errors.New(errContext, "To show credential output, credential must be provided")
	}

	if credential.AllowUseSSHAgent {
		detail := "Use SSH agent"

		if credential.GitSSHUser != "" {
			detail = fmt.Sprintf("%s, with git user '%s'", detail, credential.GitSSHUser)
		}

		return SSHAgentType, detail, nil
	} else {
		return "", "", nil
	}
}
