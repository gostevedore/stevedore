package gitauth

import (
	errors "github.com/apenella/go-common-utils/error"
	gitcontextbasicauth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
	gitcontextkeyauth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	gitcontextsshagentauth "github.com/apenella/go-docker-builder/pkg/auth/git/sshagent"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// GitAuthFactory is a factory for creating GitAuther
type GitAuthFactory struct {
	Credentials repository.CredentialsStorer
}

// NewGitAuthFactory creates a new GitAuthFactory
func NewGitAuthFactory(credentials repository.CredentialsStorer) *GitAuthFactory {
	return &GitAuthFactory{
		Credentials: credentials,
	}
}

// GenerateAuthMethod returns a new auth method based on the given context
func (f *GitAuthFactory) GenerateAuthMethod(options *builder.DockerDriverGitContextAuthOptions) (GitAuther, error) {

	errContext := "(GitAuthFactory::GenerateAuthMethod)"

	if options == nil {
		return nil, errors.New(errContext, "Git context auth options is required to generate an an auth method")
	}

	if options.CredentialsID != "" {
		if f.Credentials == nil {
			return nil, errors.New(errContext, "Credentials store is expected when a credentials id is configured")
		}

		cred, err := f.Credentials.Get(options.CredentialsID)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		options.Username = cred.Username
		options.Password = cred.Password
	}

	if options.Username != "" && options.Password != "" {
		auth := &gitcontextbasicauth.BasicAuth{
			Username: options.Username,
			Password: options.Password,
		}
		return auth, nil
	}

	if options.PrivateKeyFile != "" {
		auth := &gitcontextkeyauth.KeyAuth{
			PkFile: options.PrivateKeyFile,
		}

		if options.PrivateKeyPassword != "" {
			auth.PkPassword = options.PrivateKeyPassword
		}

		if options.GitSSHUser != "" {
			auth.GitSSHUser = options.GitSSHUser
		}

		return auth, nil
	}

	auth := &gitcontextsshagentauth.SSHAgentAuth{}
	if options.GitSSHUser != "" {
		auth.GitSSHUser = options.GitSSHUser
	}

	return auth, nil
}
