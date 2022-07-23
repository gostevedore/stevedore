package badge

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type CredentialsCompatibility struct {
	compatibility Compatibilitier
}

func NewCredentialsCompatibility(compatibility Compatibilitier) *CredentialsCompatibility {
	return &CredentialsCompatibility{
		compatibility: compatibility,
	}
}

func (c *CredentialsCompatibility) CheckCompatibility(badge *credentials.Badge) error {

	errContext := "(credentials::compatibility::CheckCompatibility)"

	if c.compatibility == nil {
		return errors.New(errContext, "To check credentials badge compatibility, compatibilitier must be provided")
	}

	if badge == nil {
		return errors.New(errContext, "To check credentials badge compatibility, badge must be provided")
	}

	if badge.DEPRECATEDUsername != "" {
		c.compatibility.AddDeprecated("'docker_login_username' is deprecated and will be removed on v0.12.0, please use 'username' instead")
		badge.Username = badge.DEPRECATEDUsername
	}

	if badge.DEPRECATEDPassword != "" {
		c.compatibility.AddDeprecated("'docker_login_password' is deprecated and will be removed on v0.12.0, please use 'password' instead")
		badge.Password = badge.DEPRECATEDPassword
	}

	return nil
}
