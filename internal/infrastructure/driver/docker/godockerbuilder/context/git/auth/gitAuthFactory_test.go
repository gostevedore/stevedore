package gitauth

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gitcontextbasicauth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
	gitcontextkeyauth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	gitcontextsshagentauth "github.com/apenella/go-docker-builder/pkg/auth/git/sshagent"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	credentialsstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAuthMethod(t *testing.T) {
	errContext := "(GitAuthFactory::GenerateAuthMethod)"
	tests := []struct {
		desc    string
		options *builder.DockerDriverGitContextAuthOptions
		factory *GitAuthFactory
		res     GitAuther
		err     error
	}{
		{
			desc:    "Testing error when options is nil",
			options: nil,
			factory: nil,
			res:     nil,
			err:     errors.New(errContext, "Git context auth options is required to generate an an auth method"),
		},
		{
			desc: "Testing generate basic auth authorization method",
			options: &builder.DockerDriverGitContextAuthOptions{
				Username: "user",
				Password: "pass",
			},
			factory: &GitAuthFactory{},
			res: &gitcontextbasicauth.BasicAuth{
				Username: "user",
				Password: "pass",
			},
			err: errors.New(errContext, "Git context auth options is required to generate an an auth method"),
		},
		{
			desc: "Testing generate private key auth authorization method",
			options: &builder.DockerDriverGitContextAuthOptions{
				PrivateKeyFile:     "keyfile",
				PrivateKeyPassword: "keypass",
				GitSSHUser:         "user",
			},
			factory: &GitAuthFactory{},
			res: &gitcontextkeyauth.KeyAuth{
				PkFile:     "keyfile",
				PkPassword: "keypass",
				GitSSHUser: "user",
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing generate ssh-agent auth authorization method",
			options: &builder.DockerDriverGitContextAuthOptions{
				GitSSHUser: "user",
			},
			factory: &GitAuthFactory{},
			res: &gitcontextsshagentauth.SSHAgentAuth{
				GitSSHUser: "user",
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing generate authorization method from credentials id",
			options: &builder.DockerDriverGitContextAuthOptions{
				CredentialsID: "test-credentials",
			},
			factory: &GitAuthFactory{
				Credentials: &credentialsstore.CredentialsStore{
					Store: map[string]*credentials.UserPasswordAuth{
						"1c88d75d861f84fd80b43bb117b2fcde": {
							Username: "user",
							Password: "pass",
						},
					},
				},
			},
			res: &gitcontextbasicauth.BasicAuth{
				Username: "user",
				Password: "pass",
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error when credentials store is nil and is passed a credentials id",
			options: &builder.DockerDriverGitContextAuthOptions{
				CredentialsID: "test-credentials",
			},
			factory: &GitAuthFactory{
				Credentials: nil,
			},
			res: &gitcontextbasicauth.BasicAuth{
				Username: "user",
				Password: "pass",
			},
			err: errors.New(errContext, "Credentials store is expected when a credentials id is configured"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			auth, err := test.factory.GenerateAuthMethod(test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, auth, test.res)
			}
		})
	}

}
