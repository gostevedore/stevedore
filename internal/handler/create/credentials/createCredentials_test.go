package credentials

import (
	"context"
	"testing"

	application "github.com/gostevedore/stevedore/internal/application/create/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		desc              string
		handler           *CreateCredentialsHandler
		id                string
		options           *Options
		prepareAssertFunc func(*CreateCredentialsHandler)
		err               error
	}{
		{
			desc: "Testing run create credentials handler",
			handler: NewCreateCredentialsHandler(
				WithApplication(application.NewMockCreateCredentialsApplication()),
			),
			id: "id",
			options: &Options{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(h *CreateCredentialsHandler) {
				h.app.(*application.MockCreateCredentialsApplication).On(
					"Run",
					context.TODO(),
					"id",
					&credentials.Credential{
						Username:                  "username",
						Password:                  "password",
						AWSSharedConfigFiles:      []string{},
						AWSSharedCredentialsFiles: []string{},
					},
					mock.Anything,
				).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.handler)
			}

			err := test.handler.Handler(context.TODO(), test.id, test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.handler.app.(*application.MockCreateCredentialsApplication).AssertExpectations(t)
			}
		})
	}
}

func TestCreateCredentialFromOptions(t *testing.T) {
	tests := []struct {
		desc    string
		options *Options
		res     *credentials.Credential
	}{
		{
			desc: "Testing create credentials from options",
			options: &Options{
				AllowUseSSHAgent:              true,
				AWSAccessKeyID:                "AWSAccessKeyID",
				AWSProfile:                    "AWSProfile",
				AWSRegion:                     "AWSRegion",
				AWSRoleARN:                    "AWSRoleARN",
				AWSSecretAccessKey:            "AWSSecretAccessKey",
				AWSSharedConfigFiles:          []string{"AWSSharedConfigFiles"},
				AWSSharedCredentialsFiles:     []string{"AWSSharedCredentialsFiles"},
				AWSUseDefaultCredentialsChain: true,
				GitSSHUser:                    "GitSSHUser",
				Password:                      "Password",
				PrivateKeyFile:                "PrivateKeyFile",
				PrivateKeyPassword:            "PrivateKeyPassword",
				Username:                      "Username",
			},
			res: &credentials.Credential{
				AWSAccessKeyID:                "AWSAccessKeyID",
				AWSRegion:                     "AWSRegion",
				AWSRoleARN:                    "AWSRoleARN",
				AWSSecretAccessKey:            "AWSSecretAccessKey",
				AWSProfile:                    "AWSProfile",
				AWSSharedCredentialsFiles:     []string{"AWSSharedCredentialsFiles"},
				AWSSharedConfigFiles:          []string{"AWSSharedConfigFiles"},
				AWSUseDefaultCredentialsChain: true,
				Password:                      "Password",
				Username:                      "Username",
				PrivateKeyFile:                "PrivateKeyFile",
				PrivateKeyPassword:            "PrivateKeyPassword",
				GitSSHUser:                    "GitSSHUser",
				AllowUseSSHAgent:              true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			credential := createCredentialFromOptions(test.options)
			assert.Equal(t, test.res, credential)
		})
	}
}
