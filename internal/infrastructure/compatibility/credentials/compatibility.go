package credentials

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

func (c *CredentialsCompatibility) CheckCompatibility(credential *credentials.Credential) error {

	errContext := "(credentials::compatibility::CheckCompatibility)"

	if c.compatibility == nil {
		return errors.New(errContext, "To check credentials credential compatibility, compatibilitier must be provided")
	}

	if credential == nil {
		return errors.New(errContext, "To check credentials credential compatibility, credential must be provided")
	}

	if credential.DEPRECATEDUsername != "" {
		c.compatibility.AddDeprecated("'docker_login_username' is deprecated and will be removed on v0.12.0, please use 'username' instead")
		credential.Username = credential.DEPRECATEDUsername
	}

	if credential.DEPRECATEDPassword != "" {
		c.compatibility.AddDeprecated("'docker_login_password' is deprecated and will be removed on v0.12.0, please use 'password' instead")
		credential.Password = credential.DEPRECATEDPassword
	}

	return nil
}
