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

// AuthMethodConstructor return BasicAuthMethod from the given badge
func (a *BasicAuthMethod) AuthMethodConstructor(badge *credentials.Badge) (repository.AuthMethodReader, error) {

	if badge == nil {
		return nil, nil
	}

	if badge.Username != "" && badge.Password != "" {
		a = &BasicAuthMethod{
			Username: badge.Username,
			Password: badge.Password,
		}

		return a, nil
	} else {
		return nil, nil
	}
}

// Name returns the name of the authentication method
func (a *BasicAuthMethod) Name() string {
	return credentials.BasicAuthMethod
}
