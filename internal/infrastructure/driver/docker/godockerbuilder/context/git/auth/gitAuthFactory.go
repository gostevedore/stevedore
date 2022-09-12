package gitauth

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	gitcontextbasicauth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
	gitcontextkeyauth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	gitcontextsshagentauth "github.com/apenella/go-docker-builder/pkg/auth/git/sshagent"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/sshagent"
)

// GitAuthFactory is a factory for creating GitAuther
type GitAuthFactory struct {
	Credentials repository.CredentialsFactorier
}

// NewGitAuthFactory creates a new GitAuthFactory
func NewGitAuthFactory(credentials repository.CredentialsFactorier) *GitAuthFactory {
	return &GitAuthFactory{
		Credentials: credentials,
	}
}

// GenerateAuthMethod returns a new auth method based on the given context
func (f *GitAuthFactory) GenerateAuthMethod(options *builder.DockerDriverGitContextAuthOptions) (GitAuther, error) {

	var err error
	var auth GitAuther

	errContext := "(GitAuthFactory::GenerateAuthMethod)"

	if options == nil {
		return nil, errors.New(errContext, "Git context auth options is required to generate an an auth method")
	}

	if options.CredentialsID != "" {
		if f.Credentials == nil {
			return nil, errors.New(errContext, "Credentials store is expected when a credentials id is configured")
		}

		auth, err = f.generateAuthMethodFromCredentials(options.CredentialsID)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
	}

	if auth == nil {
		auth, err = f.generateAuthMethodFromOptions(options)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
	}

	return auth, nil
}

func (f *GitAuthFactory) generateAuthMethodFromCredentials(id string) (GitAuther, error) {

	var authMethod repository.AuthMethodReader
	var auth GitAuther
	var err error
	errContext := "(gitauth::generateAuthMethodFromCredentials)"

	authMethod, err = f.Credentials.Get(id)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if authMethod == nil {
		return nil, errors.New(errContext, fmt.Sprintf("Credentials with id '%s' not found", id))
	}

	switch authMethod.Name() {
	case credentials.BasicAuthMethod:
		username := authMethod.(*basic.BasicAuthMethod).Username
		password := authMethod.(*basic.BasicAuthMethod).Password

		if username != "" && password != "" {
			auth = &gitcontextbasicauth.BasicAuth{
				Username: username,
				Password: password,
			}
		}
	case credentials.KeyFileAuthMethod:
		key := authMethod.(*keyfile.KeyFileAuthMethod).PrivateKeyFile
		if key != "" {
			auth = &gitcontextkeyauth.KeyAuth{
				PkFile: key,
			}

			password := authMethod.(*keyfile.KeyFileAuthMethod).PrivateKeyPassword
			if password != "" {
				auth.(*gitcontextkeyauth.KeyAuth).PkPassword = password
			}

			user := authMethod.(*keyfile.KeyFileAuthMethod).GitSSHUser
			if user != "" {
				auth.(*gitcontextkeyauth.KeyAuth).GitSSHUser = user
			}
		}
	default:
		auth = &gitcontextsshagentauth.SSHAgentAuth{}

		user := authMethod.(*sshagent.SSHAgentAuthMethod).GitSSHUser
		if user != "" {
			auth.(*gitcontextsshagentauth.SSHAgentAuth).GitSSHUser = user
		}
	}

	return auth, nil
}

func (f *GitAuthFactory) generateAuthMethodFromOptions(options *builder.DockerDriverGitContextAuthOptions) (GitAuther, error) {
	var auth GitAuther

	if options == nil {
		return nil, nil
	}

	if options.Username != "" && options.Password != "" {
		auth = &gitcontextbasicauth.BasicAuth{
			Username: options.Username,
			Password: options.Password,
		}

		return auth, nil
	}

	if options.PrivateKeyFile != "" {
		auth = &gitcontextkeyauth.KeyAuth{
			PkFile: options.PrivateKeyFile,
		}

		if options.PrivateKeyPassword != "" {
			auth.(*gitcontextkeyauth.KeyAuth).PkPassword = options.PrivateKeyPassword
		}

		if options.GitSSHUser != "" {
			auth.(*gitcontextkeyauth.KeyAuth).GitSSHUser = options.GitSSHUser
		}

		return auth, nil
	}

	auth = &gitcontextsshagentauth.SSHAgentAuth{}
	if options.GitSSHUser != "" {
		auth.(*gitcontextsshagentauth.SSHAgentAuth).GitSSHUser = options.GitSSHUser
	}

	return auth, nil
}
