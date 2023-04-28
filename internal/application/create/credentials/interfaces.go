package credentials

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

// CredentialsStorer interface defines the storage of credentials
type CredentialsStorer interface {
	Store(id string, credential *credentials.Credential) error
}
