package envvars

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

// EnvvarsStore is a store for credentials
type EnvvarsStore struct{}

// NewEnvvarsStore creates a new mocked store for credentials
func NewEnvvarsStore() *EnvvarsStore {
	return &EnvvarsStore{}
}

// Store stores a badge
func (m *EnvvarsStore) Store(id string, badge *credentials.Badge) error {
	// TODO: create envvars
	return nil
}

// Get returns a auth for the badge id
func (m *EnvvarsStore) Get(id string) (*credentials.Badge, error) {
	// TODO: get envvars
	return nil, nil
}

// All returns all badges
func (m *EnvvarsStore) All() []*credentials.Badge {

	// TODO: get all envvars

	return nil
}
