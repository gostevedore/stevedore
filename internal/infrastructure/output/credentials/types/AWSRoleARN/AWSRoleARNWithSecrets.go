package awsrolearn

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type AWSRoleARNWithSecretsOutput struct {
	output *AWSRoleARNOutput
}

func NewAWSRoleARNWithSecretsOutput(o *AWSRoleARNOutput) *AWSRoleARNWithSecretsOutput {
	return &AWSRoleARNWithSecretsOutput{
		output: o,
	}
}

func (o *AWSRoleARNWithSecretsOutput) Output(badge *credentials.Badge) (string, string, error) {
	errContext := "(output::credentials::types::AWSRoleARNWithSecretsOutput::Output)"

	if o.output == nil {
		return "", "", errors.New(errContext, "AWS role arn with secret output requieres an output")
	}

	badgeType, details, err := o.output.Output(badge)
	if err != nil {
		return "", "", errors.New(errContext, "", err)
	}

	if badge.AWSAccessKeyID != "" && badge.AWSSecretAccessKey != "" {
		if details != "" {
			details = fmt.Sprintf("%s,", details)
		}
		details = fmt.Sprintf("%s aws_secret_access_key=%s", details, badge.AWSSecretAccessKey)
	}

	return badgeType, details, nil
}
