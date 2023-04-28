package awsrolearn

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

const (
	// AWSRoleARNType is the name of the basic authentication method
	AWSRoleARNType = "AWS role arn"
)

type AWSRoleARNOutput struct{}

func NewAWSRoleARNOutput() *AWSRoleARNOutput {
	return &AWSRoleARNOutput{}
}

func (o *AWSRoleARNOutput) Output(credential *credentials.Credential) (string, string, error) {

	errContext := "(output::credentials::types::AWSRoleARNOutput::Output)"

	if credential == nil {
		return "", "", errors.New(errContext, "To show credential output, credential must be provided")
	}

	if credential.AWSRoleARN != "" {
		detail := fmt.Sprintf("role_arn=%s", credential.AWSRoleARN)

		if credential.AWSRegion != "" {
			detail = fmt.Sprintf("%s, region=%s", detail, credential.AWSRegion)
		}

		if credential.AWSProfile != "" {
			detail = fmt.Sprintf("%s, profile=%s", detail, credential.AWSProfile)
		}

		if credential.AWSAccessKeyID != "" && credential.AWSSecretAccessKey != "" {
			detail = fmt.Sprintf("%s, aws_access_key_id=%s", detail, credential.AWSAccessKeyID)
		}

		return AWSRoleARNType, detail, nil
	} else {
		return "", "", nil
	}
}
