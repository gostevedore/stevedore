package sshagent

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthMethod(t *testing.T) {
	tests := []struct {
		desc       string
		method     *SSHAgentAuthMethod
		credential *credentials.Credential
		res        repository.AuthMethodReader
		err        error
	}{
		{
			desc:       "Testing get auth method with nil credential",
			method:     NewSSHAgentAuthMethod(),
			credential: nil,
			res:        nil,
		},
		{
			desc:       "Testing get auth method with false value on allow use ssh agent on the credential",
			method:     NewSSHAgentAuthMethod(),
			credential: &credentials.Credential{AllowUseSSHAgent: false},
			res:        nil,
		},
		{
			desc:       "Testing get auth method with allow use ssh agent and a user defined on the credential",
			method:     NewSSHAgentAuthMethod(),
			credential: &credentials.Credential{AllowUseSSHAgent: true, GitSSHUser: "user"},
			res:        &SSHAgentAuthMethod{GitSSHUser: "user"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			auth, err := test.method.AuthMethodConstructor(test.credential)
			if err != nil {
				assert.Equal(t, test.res, err)
			} else {
				assert.Equal(t, test.res, auth)
			}

		})
	}
}

func TestName(t *testing.T) {
	method := NewSSHAgentAuthMethod()
	assert.Equal(t, credentials.SSHAgentAuthMethod, method.Name())
}
