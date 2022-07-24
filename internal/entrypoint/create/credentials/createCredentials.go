package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *CreateCredentialsEntrypoint)

// CreateCredentialsEntrypoint defines the entrypoint for the application
type CreateCredentialsEntrypoint struct{}

// NewCreateCredentialsEntrypoint returns a new entrypoint
func NewCreateCredentialsEntrypoint(opts ...OptionsFunc) *CreateCredentialsEntrypoint {
	e := &CreateCredentialsEntrypoint{}
	e.Options(opts...)

	return e
}

// Options provides the options for the entrypoint
func (e *CreateCredentialsEntrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute provides a mock function
func (e *CreateCredentialsEntrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration) error {
	return nil
}
