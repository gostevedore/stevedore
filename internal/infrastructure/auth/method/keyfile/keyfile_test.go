package keyfile

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthMethod(t *testing.T) {
	tests := []struct {
		desc       string
		method     *KeyFileAuthMethod
		credential *credentials.Credential
		res        repository.AuthMethodReader
		err        error
	}{
		{
			desc:       "Testing get auth method with nil credential",
			method:     NewKeyFileAuthMethod(),
			credential: nil,
			res:        nil,
		},
		{
			desc:       "Testing get auth method private key, password and user defined on the credential",
			method:     NewKeyFileAuthMethod(),
			credential: &credentials.Credential{PrivateKeyFile: "private key", PrivateKeyPassword: "password", GitSSHUser: "user"},
			res:        &KeyFileAuthMethod{PrivateKeyFile: "private key", PrivateKeyPassword: "password", GitSSHUser: "user"},
		},
		{
			desc:       "Testing get auth method private key not defined on the credential",
			method:     NewKeyFileAuthMethod(),
			credential: &credentials.Credential{},
			res:        nil,
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
	method := NewKeyFileAuthMethod()
	assert.Equal(t, credentials.KeyFileAuthMethod, method.Name())
}
