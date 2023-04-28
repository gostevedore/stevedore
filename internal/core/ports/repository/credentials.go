package repository

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

// AuthFactorier
type AuthFactorier interface {
	Get(id string) (AuthMethodReader, error)
}

// CredentialsStorer is a repository for credentials
type CredentialsStorer interface {
	Get(id string) (*credentials.Credential, error)
	Store(id string, credential *credentials.Credential) error
}

// AuthProviderer interface that provides authentication
type AuthProviderer interface {
	Get(credential *credentials.Credential) (AuthMethodReader, error)
}

// AuthMethodReader interface that provides authentication method data
type AuthMethodReader interface {
	Name() string
}

// AuthMethodConstructor interface that creates authentication method data
type AuthMethodConstructor interface {
	AuthMethodConstructor(credential *credentials.Credential) (AuthMethodReader, error)
}

// Formater interface to marshal or unmarshal bagde data
type Formater interface {
	Marshaler
	Unmarshaler
}

// Marshaler is used to format the credential before persisting it, such as JSON, YAML,...
type Marshaler interface {
	Marshal(credential *credentials.Credential) (string, error)
}

// Unmarshaler is used to parse the credential after retrieving it, such as JSON, YAML,...
type Unmarshaler interface {
	Unmarshal(data []byte) (*credentials.Credential, error)
}

// CredentialsFilterer is an interface for filtering credentials content output
type CredentialsFilterer interface {
	All() ([]*credentials.Credential, error)
	Get(id string) (*credentials.Credential, error)
}

// CredentialsPrinter is an interface for printing credentials content output
type CredentialsPrinter interface {
	Print(credentials []*credentials.Credential) error
}
