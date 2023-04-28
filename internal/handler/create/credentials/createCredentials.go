package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

// OptionsFunc is a function used to configure the handler
type OptionsFunc func(*CreateCredentialsHandler)

// CreateCredentialsHandler is a handler for get credentials commands
type CreateCredentialsHandler struct {
	app Applicationer
}

// NewCreateCredentialsHandler creates a new handler for build commands
func NewCreateCredentialsHandler(options ...OptionsFunc) *CreateCredentialsHandler {
	handler := &CreateCredentialsHandler{}
	handler.Options(options...)

	return handler
}

func WithApplication(app Applicationer) OptionsFunc {
	return func(h *CreateCredentialsHandler) {
		h.app = app
	}
}

// Options configure the service
func (h *CreateCredentialsHandler) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// CreateCredentialsHandler handles build commands
func (h *CreateCredentialsHandler) Handler(ctx context.Context, id string, options *Options) error {
	var err error

	errContext := "(create/credentials::Handler)"

	credential := createCredentialFromOptions(options)

	err = h.app.Run(ctx, id, credential)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func createCredentialFromOptions(options *Options) *credentials.Credential {
	credential := &credentials.Credential{}

	credential.AllowUseSSHAgent = options.AllowUseSSHAgent
	credential.AWSAccessKeyID = options.AWSAccessKeyID
	credential.AWSProfile = options.AWSProfile
	credential.AWSRegion = options.AWSRegion
	credential.AWSRoleARN = options.AWSRoleARN
	credential.AWSSecretAccessKey = options.AWSSecretAccessKey
	credential.AWSSharedConfigFiles = append([]string{}, options.AWSSharedConfigFiles...)
	credential.AWSSharedCredentialsFiles = append([]string{}, options.AWSSharedCredentialsFiles...)
	credential.AWSUseDefaultCredentialsChain = options.AWSUseDefaultCredentialsChain
	credential.GitSSHUser = options.GitSSHUser
	credential.Password = options.Password
	credential.PrivateKeyFile = options.PrivateKeyFile
	credential.PrivateKeyPassword = options.PrivateKeyPassword
	credential.Username = options.Username

	return credential
}
