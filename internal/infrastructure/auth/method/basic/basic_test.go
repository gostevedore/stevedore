package basic

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthMethod(t *testing.T) {
	tests := []struct {
		desc       string
		method     *BasicAuthMethod
		credential *credentials.Credential
		res        repository.AuthMethodReader
		err        error
	}{
		{
			desc:       "Testing get auth method with nil credential",
			method:     NewBasicAuthMethod(),
			credential: nil,
			res:        nil,
		},
		{
			desc:       "Testing get auth method username and password defined on the credential",
			method:     NewBasicAuthMethod(),
			credential: &credentials.Credential{Username: "username", Password: "password"},
			res:        &BasicAuthMethod{Username: "username", Password: "password"},
		},
		{
			desc:       "Testing get auth method username and password not defined on the credential",
			method:     NewBasicAuthMethod(),
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
	method := NewBasicAuthMethod()
	assert.Equal(t, credentials.BasicAuthMethod, method.Name())
}
