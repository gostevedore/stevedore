package sshagent

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthMethod(t *testing.T) {
	tests := []struct {
		desc   string
		method *SSHAgentAuthMethod
		badge  *credentials.Badge
		res    repository.AuthMethodReader
		err    error
	}{
		{
			desc:   "Testing get auth method with nil badge",
			method: NewSSHAgentAuthMethod(),
			badge:  nil,
			res:    nil,
		},
		{
			desc:   "Testing get auth method with false value on allow use ssh agent on the badge",
			method: NewSSHAgentAuthMethod(),
			badge:  &credentials.Badge{AllowUseSSHAgent: false},
			res:    nil,
		},
		{
			desc:   "Testing get auth method with allow use ssh agent and a user defined on the badge",
			method: NewSSHAgentAuthMethod(),
			badge:  &credentials.Badge{AllowUseSSHAgent: true, GitSSHUser: "user"},
			res:    &SSHAgentAuthMethod{GitSSHUser: "user"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			auth, err := test.method.AuthMethodConstructor(test.badge)
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
