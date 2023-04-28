package privatekeyfile

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type PrivateKeyFileWithSecretsOutput struct {
	output *PrivateKeyFileOutput
}

func NewPrivateKeyFileWithSecretsOutput(o *PrivateKeyFileOutput) *PrivateKeyFileWithSecretsOutput {
	return &PrivateKeyFileWithSecretsOutput{
		output: o,
	}
}

func (o *PrivateKeyFileWithSecretsOutput) Output(credential *credentials.Credential) (string, string, error) {
	errContext := "(output::credentials::types::PrivateKeyFileWithSecretsOutput::Output)"

	if o.output == nil {
		return "", "", errors.New(errContext, "Private key file with secret output requieres an output")
	}

	credentialType, details, err := o.output.Output(credential)
	if err != nil {
		return "", "", errors.New(errContext, "", err)
	}

	if credential.PrivateKeyPassword != "" {
		if details != "" {
			details = fmt.Sprintf("%s,", details)
		}
		details = fmt.Sprintf("%s private_key_password=%s", details, credential.PrivateKeyPassword)
	}

	return credentialType, details, nil

}
