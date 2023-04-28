package store

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/sshagent"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tests := []struct {
		desc     string
		badge    *credentials.Credential
		provider *StoreAuthProvider
		res      repository.AuthMethodReader
		err      error
	}{
		{
			desc:     "Testing get auth method from badge credentials provider when badge is nil",
			badge:    nil,
			provider: NewStoreAuthProvider(),
			res:      nil,
		},
		{
			desc:  "Testing get auth method from badge credentials provider when badge is not nil and no auth method is found",
			badge: &credentials.Credential{},
			provider: NewStoreAuthProvider(
				keyfile.NewKeyFileAuthMethod(),
				basic.NewBasicAuthMethod(),
				sshagent.NewSSHAgentAuthMethod(),
			),
			res: nil,
		},
		{
			desc: "Testing get auth method from badge credentials provider when badge is not nil and no auth method is found",
			badge: &credentials.Credential{
				Username: "username",
				Password: "password",
			},
			provider: NewStoreAuthProvider(
				keyfile.NewKeyFileAuthMethod(),
				basic.NewBasicAuthMethod(),
				sshagent.NewSSHAgentAuthMethod(),
			),
			res: &basic.BasicAuthMethod{
				Username: "username",
				Password: "password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			auth, err := test.provider.Get(test.badge)
			if err != nil {
				assert.Equal(t, test.res, err)
			} else {
				assert.Equal(t, test.res, auth)
			}
		})
	}
}
