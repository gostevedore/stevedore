package local

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

// LocalStoreWithSafeStore is a local store for credentials
type LocalStoreWithSafeStore struct {
	store *LocalStore
}

// NewLocalStoreWithSafeStore creates a new local store for credentials
func NewLocalStoreWithSafeStore(opts ...OptionsFunc) *LocalStoreWithSafeStore {
	s := &LocalStoreWithSafeStore{
		NewLocalStore(opts...),
	}

	return s
}

func (s *LocalStoreWithSafeStore) Store(id string, credential *credentials.Credential) error {
	errContext := "(store::credentials::local::LocalStoreWithSafeStore)"

	err := s.store.SafeStore(id, credential)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
