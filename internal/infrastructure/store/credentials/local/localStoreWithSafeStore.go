package local

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/spf13/afero"
)

// LocalStoreWithSafeStore is a local store for credentials
type LocalStoreWithSafeStore struct {
	store *LocalStore
}

// NewLocalStoreWithSafeStore creates a new local store for credentials
func NewLocalStoreWithSafeStore(fs afero.Fs, path string, f repository.Formater, compatibility CredentialsCompatibilier) *LocalStoreWithSafeStore {

	return &LocalStoreWithSafeStore{
		NewLocalStore(fs, path, f, compatibility),
	}
}

func (s *LocalStoreWithSafeStore) Store(id string, badge *credentials.Badge) error {
	errContext := "(store::credentials::local::LocalStoreWithSafeStore)"

	err := s.store.SafeStore(id, badge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
