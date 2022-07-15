package credentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// OptionsFunc defines the signature for an option function to set entrypoint attributes
type OptionsFunc func(opts *Entrypoint)

// Entrypoint defines the entrypoint for the application
type Entrypoint struct{}

// NewEntrypoint returns a new entrypoint
func NewEntrypoint(opts ...OptionsFunc) *Entrypoint {
	e := &Entrypoint{}
	e.Options(opts...)

	return e
}

// Options provides the options for the entrypoint
func (e *Entrypoint) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// Execute provides a mock function
func (e *Entrypoint) Execute(ctx context.Context, args []string, conf *configuration.Configuration) error {
	return nil
}
