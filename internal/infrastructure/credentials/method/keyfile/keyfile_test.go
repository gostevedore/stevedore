package keyfile

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthMethod(t *testing.T) {
	tests := []struct {
		desc   string
		method *KeyFileAuthMethod
		badge  *credentials.Badge
		res    repository.AuthMethodReader
		err    error
	}{
		{
			desc:   "Testing get auth method with nil badge",
			method: NewKeyFileAuthMethod(),
			badge:  nil,
			res:    nil,
		},
		{
			desc:   "Testing get auth method private key, password and user defined on the badge",
			method: NewKeyFileAuthMethod(),
			badge:  &credentials.Badge{PrivateKeyFile: "private key", PrivateKeyPassword: "password", GitSSHUser: "user"},
			res:    &KeyFileAuthMethod{PrivateKeyFile: "private key", PrivateKeyPassword: "password", GitSSHUser: "user"},
		},
		{
			desc:   "Testing get auth method private key not defined on the badge",
			method: NewKeyFileAuthMethod(),
			badge:  &credentials.Badge{},
			res:    nil,
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
	method := NewKeyFileAuthMethod()
	assert.Equal(t, credentials.KeyFileAuthMethod, method.Name())
}
