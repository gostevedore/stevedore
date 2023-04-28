package sshagent

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	errContext := "(output::credentials::types::SSHAgentOutput::Output)"

	tests := []struct {
		desc            string
		output          *SSHAgentOutput
		credential      *credentials.Credential
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for SSHAgentOutput and credential is nil",
			output:          NewSSHAgentOutput(),
			credential:      nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "To show credential output, credential must be provided"),
		},
		{
			desc:   "Testing generate output for AWSDefaultCredentialsChain",
			output: NewSSHAgentOutput(),
			credential: &credentials.Credential{
				AllowUseSSHAgent: true,
				GitSSHUser:       "user",
			},
			detail:          "Use SSH agent, with git user 'user'",
			credentialsType: SSHAgentType,
			err:             &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			credentialsType, detail, err := test.output.Output(test.credential)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.credentialsType, credentialsType)
				assert.Equal(t, test.detail, detail)
			}
		})
	}
}
