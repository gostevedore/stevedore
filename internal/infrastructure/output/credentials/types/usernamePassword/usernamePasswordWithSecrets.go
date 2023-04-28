package usernamepassword

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type UsernamePasswordWithSecretsOutput struct {
	output *UsernamePasswordOutput
}

func NewUsernamePasswordWithSecretsOutput(o *UsernamePasswordOutput) *UsernamePasswordWithSecretsOutput {
	return &UsernamePasswordWithSecretsOutput{
		output: o,
	}
}

func (o *UsernamePasswordWithSecretsOutput) Output(credential *credentials.Credential) (string, string, error) {
	errContext := "(output::credentials::types::UsernamePasswordWithSecretsOutput::Output)"

	if o.output == nil {
		return "", "", errors.New(errContext, "Username-password with secret output requieres an output")
	}

	credentialType, details, err := o.output.Output(credential)
	if err != nil {
		return "", "", errors.New(errContext, "", err)
	}

	if credential.Password != "" {
		if details != "" {
			details = fmt.Sprintf("%s,", details)
		}
		details = fmt.Sprintf("%s password=%s", details, credential.Password)
	}

	return credentialType, details, nil

}
