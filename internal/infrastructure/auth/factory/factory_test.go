package factory

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/sshagent"
	credential "github.com/gostevedore/stevedore/internal/infrastructure/auth/provider/store"
	mockstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/mock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	errContext := "(credentials::factory::AuthFactory::Get)"
	tests := []struct {
		desc              string
		factory           *AuthFactory
		id                string
		prepareAssertFunc func(*AuthFactory)
		res               repository.AuthMethodReader
		err               error
	}{
		{
			desc: "Testing error on credentials factory getting credentials with empty id",
			factory: NewAuthFactory(
				mockstore.NewMockStore(),
				credential.NewStoreAuthProvider(
					keyfile.NewKeyFileAuthMethod(),
					basic.NewBasicAuthMethod(),
					sshagent.NewSSHAgentAuthMethod(),
				),
			),
			id:  "",
			err: errors.New(errContext, "To get credentials, you must provide an id"),
		},
		{
			desc: "Testing get credentials for basic auth method",
			factory: NewAuthFactory(
				mockstore.NewMockStore(),
				credential.NewStoreAuthProvider(
					keyfile.NewKeyFileAuthMethod(),
					basic.NewBasicAuthMethod(),
					sshagent.NewSSHAgentAuthMethod(),
				),
			),
			id: "credentials",
			prepareAssertFunc: func(f *AuthFactory) {
				f.store.(*mockstore.MockStore).On("Get", "credentials").Return(
					&credentials.Credential{
						Username: "username",
						Password: "password",
					},
					nil,
				)
			},
			res: &basic.BasicAuthMethod{
				Username: "username",
				Password: "password",
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get credentials for key file auth method",
			factory: NewAuthFactory(
				mockstore.NewMockStore(),
				credential.NewStoreAuthProvider(
					keyfile.NewKeyFileAuthMethod(),
					basic.NewBasicAuthMethod(),
					sshagent.NewSSHAgentAuthMethod(),
				),
			),
			id: "credentials",
			prepareAssertFunc: func(f *AuthFactory) {
				f.store.(*mockstore.MockStore).On("Get", "credentials").Return(
					&credentials.Credential{
						PrivateKeyFile:     "private_key_file",
						PrivateKeyPassword: "private_key_password",
					},
					nil,
				)
			},
			res: &keyfile.KeyFileAuthMethod{
				PrivateKeyFile:     "private_key_file",
				PrivateKeyPassword: "private_key_password",
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get credentials for sshagent auth method",
			factory: NewAuthFactory(
				mockstore.NewMockStore(),
				credential.NewStoreAuthProvider(
					keyfile.NewKeyFileAuthMethod(),
					basic.NewBasicAuthMethod(),
					sshagent.NewSSHAgentAuthMethod(),
				),
			),
			id: "credentials",
			prepareAssertFunc: func(f *AuthFactory) {
				f.store.(*mockstore.MockStore).On("Get", "credentials").Return(
					&credentials.Credential{
						AllowUseSSHAgent: true,
					},
					nil,
				)
			},
			res: &sshagent.SSHAgentAuthMethod{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.factory)
			}

			res, err := test.factory.Get(test.id)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}

}
