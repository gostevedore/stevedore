package basic

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// BasicAuthMethod is a basic authentication method
type BasicAuthMethod struct {
	Username string
	Password string
}

// NewBasicAuthMethod creates a new BasicAuthMethod
func NewBasicAuthMethod() *BasicAuthMethod {
	return &BasicAuthMethod{}
}

// AuthMethodConstructor return BasicAuthMethod from the given credential
func (auth *BasicAuthMethod) AuthMethodConstructor(credential *credentials.Credential) (repository.AuthMethodReader, error) {

	if credential == nil {
		return nil, nil
	}

	if credential.Username != "" && credential.Password != "" {
		auth = &BasicAuthMethod{
			Username: credential.Username,
			Password: credential.Password,
		}

		return auth, nil
	} else {
		return nil, nil
	}
}

// Name returns the name of the authentication method
func (a *BasicAuthMethod) Name() string {
	return credentials.BasicAuthMethod
}
