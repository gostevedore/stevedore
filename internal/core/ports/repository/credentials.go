package repository

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

// CredentialsFactorier
type CredentialsFactorier interface {
	Get(id string) (AuthMethodReader, error)
}

// CredentialsStorer is a repository for credentials
type CredentialsStorer interface {
	Get(id string) (*credentials.Badge, error)
	Store(id string, badge *credentials.Badge) error
}

// CredentialsProviderer interface that provides authentication
type CredentialsProviderer interface {
	Get(badge *credentials.Badge) (AuthMethodReader, error)
}

// AuthMethodReader interface that provides authentication method data
type AuthMethodReader interface {
	Name() string
}

//AuthMethodConstructor interface that creates authentication method data
type AuthMethodConstructor interface {
	AuthMethod(badge *credentials.Badge) (AuthMethodReader, error)
}
