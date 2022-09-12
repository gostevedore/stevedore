package local

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

type CredentialsCompatibilier interface {
	CheckCompatibility(badge *credentials.Badge) error
}
