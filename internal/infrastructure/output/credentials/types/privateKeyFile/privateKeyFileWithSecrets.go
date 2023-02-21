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

func (o *PrivateKeyFileWithSecretsOutput) Output(badge *credentials.Badge) (string, string, error) {
	errContext := "(output::credentials::types::PrivateKeyFileWithSecretsOutput::Output)"

	if o.output == nil {
		return "", "", errors.New(errContext, "Private key file with secret output requieres an output")
	}

	badgeType, details, err := o.output.Output(badge)
	if err != nil {
		return "", "", errors.New(errContext, "", err)
	}

	if badge.PrivateKeyPassword != "" {
		if details != "" {
			details = fmt.Sprintf("%s,", details)
		}
		details = fmt.Sprintf("%s private_key_password=%s", details, badge.PrivateKeyPassword)
	}

	return badgeType, details, nil

}
