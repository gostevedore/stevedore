package awsstaticcredentials

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type AWSStaticCredentialsWithSecretsOutput struct {
	output *AWSStaticCredentialsOutput
}

func NewAWSStaticCredentialsWithSecretsOutput(o *AWSStaticCredentialsOutput) *AWSStaticCredentialsWithSecretsOutput {
	return &AWSStaticCredentialsWithSecretsOutput{
		output: o,
	}
}

func (o *AWSStaticCredentialsWithSecretsOutput) Output(badge *credentials.Badge) (string, string, error) {
	errContext := "(output::credentials::types::AWSStaticCredentialsWithSecretsOutput::Output)"

	if o.output == nil {
		return "", "", errors.New(errContext, "AWS static credentials with secret output requieres an output")
	}

	badgeType, details, err := o.output.Output(badge)
	if err != nil {
		return "", "", errors.New(errContext, "", err)
	}

	if badge.AWSSecretAccessKey != "" {
		if details != "" {
			details = fmt.Sprintf("%s,", details)
		}
		details = fmt.Sprintf("%s aws_secret_access_key=%s", details, badge.AWSSecretAccessKey)
	}

	return badgeType, details, nil

}
