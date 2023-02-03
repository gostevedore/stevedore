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

func (o *AWSRoleARNOutput) Output(badge *credentials.Badge) (string, string, error) {

	errContext := "(output::credentials::types::AWSRoleARNOutput::Output)"

	if badge == nil {
		return "", "", errors.New(errContext, "To show badge output, badge must be provided")
	}

	if badge.AWSRoleARN != "" {
		detail := fmt.Sprintf("role_arn=%s", badge.AWSRoleARN)

		if badge.AWSRegion != "" {
			detail = fmt.Sprintf("%s, region=%s", detail, badge.AWSRegion)
		}

		if badge.AWSProfile != "" {
			detail = fmt.Sprintf("%s, profile=%s", detail, badge.AWSProfile)
		}

		if badge.AWSAccessKeyID != "" && badge.AWSSecretAccessKey != "" {
			detail = fmt.Sprintf("%s, aws_access_key_id=%s", detail, badge.AWSAccessKeyID)
		}

		return AWSRoleARNType, detail, nil
	} else {
		return "", "", nil
	}
}
