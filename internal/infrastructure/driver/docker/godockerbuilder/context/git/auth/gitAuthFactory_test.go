package gitauth

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gitcontextbasicauth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
	gitcontextkeyauth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	gitcontextsshagentauth "github.com/apenella/go-docker-builder/pkg/auth/git/sshagent"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/sshagent"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAuthMethod(t *testing.T) {
	errContext := "(GitAuthFactory::GenerateAuthMethod)"
	tests := []struct {
		desc              string
		options           *builder.DockerDriverGitContextAuthOptions
		factory           *GitAuthFactory
		prepareAssertFunc func(*GitAuthFactory)
		res               GitAuther
		err               error
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
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
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
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
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
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
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
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
			res: &gitcontextbasicauth.BasicAuth{
				Username: "user",
				Password: "pass",
			},
			err: &errors.Error{},
			prepareAssertFunc: func(f *GitAuthFactory) {
				f.Credentials.(*factory.MockCredentialsFactory).On("Get", "test-credentials").Return(
					&basic.BasicAuthMethod{
						Username: "user",
						Password: "pass",
					},
					nil,
				)
			},
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
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.factory)
			}

			auth, err := test.factory.GenerateAuthMethod(test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, auth, test.res)
			}
		})
	}
}

func TestGenerateAuthMethodFromCredentials(t *testing.T) {
	tests := []struct {
		desc              string
		credentialsID     string
		factory           *GitAuthFactory
		res               GitAuther
		err               error
		prepareAssertFunc func(*GitAuthFactory)
	}{
		{
			desc: "Testing generate basic auth authorization method",
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
			credentialsID: "registry.test",
			res:           &gitcontextbasicauth.BasicAuth{},
			err:           &errors.Error{},
			prepareAssertFunc: func(f *GitAuthFactory) {
				f.Credentials.(*factory.MockCredentialsFactory).On("Get", "registry.test").Return(
					&basic.BasicAuthMethod{
						Username: "user",
						Password: "pass",
					},
					nil,
				)
			},
		},
		{
			desc: "Testing generate basic auth authorization method when either user nor password exists",
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
			credentialsID: "registry.test",
			res:           nil,
			err:           &errors.Error{},
			prepareAssertFunc: func(f *GitAuthFactory) {
				f.Credentials.(*factory.MockCredentialsFactory).On("Get", "registry.test").Return(
					&basic.BasicAuthMethod{},
					nil,
				)
			},
		},
		{
			desc: "Testing generate key file auth authorization method",
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
			credentialsID: "registry.test",
			res:           &gitcontextkeyauth.KeyAuth{},
			err:           &errors.Error{},
			prepareAssertFunc: func(f *GitAuthFactory) {
				f.Credentials.(*factory.MockCredentialsFactory).On("Get", "registry.test").Return(
					&keyfile.KeyFileAuthMethod{
						PrivateKeyFile: "keyfile",
					},
					nil,
				)
			},
		},
		{
			desc: "Testing generate key file auth authorization method without private key file",
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
			credentialsID: "registry.test",
			res:           nil,
			err:           &errors.Error{},
			prepareAssertFunc: func(f *GitAuthFactory) {
				f.Credentials.(*factory.MockCredentialsFactory).On("Get", "registry.test").Return(
					&keyfile.KeyFileAuthMethod{},
					nil,
				)
			},
		},
		{
			desc: "Testing generate ssh-agent (default) auth authorization method",
			factory: NewGitAuthFactory(
				factory.NewMockCredentialsFactory(),
			),
			credentialsID: "registry.test",
			res:           &gitcontextsshagentauth.SSHAgentAuth{},
			err:           &errors.Error{},
			prepareAssertFunc: func(f *GitAuthFactory) {
				f.Credentials.(*factory.MockCredentialsFactory).On("Get", "registry.test").Return(
					&sshagent.SSHAgentAuthMethod{},
					nil,
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.factory)
			}

			auth, err := test.factory.generateAuthMethodFromCredentials(test.credentialsID)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, test.res, auth)
			}
		})
	}
}

func TestGenerateAuthMethodFromOptions(t *testing.T) {}
