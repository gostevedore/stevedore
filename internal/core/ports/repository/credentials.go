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

// AuthMethodConstructor interface that creates authentication method data
type AuthMethodConstructor interface {
	AuthMethodConstructor(badge *credentials.Badge) (AuthMethodReader, error)
}

// Formater interface to marshal or unmarshal bagde data
type Formater interface {
	Marshaler
	Unmarshaler
}

// Marshaler is used to format the badge before persisting it, such as JSON, YAML,...
type Marshaler interface {
	Marshal(badge *credentials.Badge) (string, error)
}

// Unmarshaler is used to parse the badge after retrieving it, such as JSON, YAML,...
type Unmarshaler interface {
	Unmarshal(data []byte) (*credentials.Badge, error)
}

// CredentialsFilterer is an interface for filtering credentials content output
type CredentialsFilterer interface {
	All() ([]*credentials.Badge, error)
	Get(id string) (*credentials.Badge, error)
}

// CredentialsPrinter is an interface for printing credentials content output
type CredentialsPrinter interface {
	Print(badges []*credentials.Badge) error
}
