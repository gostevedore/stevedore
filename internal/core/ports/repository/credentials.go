package repository

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

// CredentialsStorer
type CredentialsStorer interface {
	Get(id string) (*credentials.UserPasswordAuth, error)
}
