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

	badge := createBadgeFromOptions(options)

	err = h.app.Run(ctx, id, badge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func createBadgeFromOptions(options *Options) *credentials.Badge {
	badge := &credentials.Badge{}

	badge.AllowUseSSHAgent = options.AllowUseSSHAgent
	badge.AWSAccessKeyID = options.AWSAccessKeyID
	badge.AWSProfile = options.AWSProfile
	badge.AWSRegion = options.AWSRegion
	badge.AWSRoleARN = options.AWSRoleARN
	badge.AWSSecretAccessKey = options.AWSSecretAccessKey
	badge.AWSSharedConfigFiles = append([]string{}, options.AWSSharedConfigFiles...)
	badge.AWSSharedCredentialsFiles = append([]string{}, options.AWSSharedCredentialsFiles...)
	badge.AWSUseDefaultCredentialsChain = options.AWSUseDefaultCredentialsChain
	badge.GitSSHUser = options.GitSSHUser
	badge.Password = options.Password
	badge.PrivateKeyFile = options.PrivateKeyFile
	badge.PrivateKeyPassword = options.PrivateKeyPassword
	badge.Username = options.Username

	return badge
}
