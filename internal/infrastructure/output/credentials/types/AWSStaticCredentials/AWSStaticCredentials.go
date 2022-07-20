package awsstaticcredentials

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

const (
	// AWSStaticCredentialsType is the name of the basic authentication method
	AWSStaticCredentialsType = "AWS static credentials"
)

type AWSStaticCredentialsOutput struct{}

func NewAWSStaticCredentialsOutput() *AWSStaticCredentialsOutput {
	return &AWSStaticCredentialsOutput{}
}

func (o *AWSStaticCredentialsOutput) Output(badge *credentials.Badge) (string, string, error) {

	errContext := "(output::credentials::types::AWSStaticCredentialsOutput::Output)"

	if badge == nil {
		return "", "", errors.New(errContext, "To show badge output, badge must be provided")
	}

	if badge.AWSAccessKeyID != "" && badge.AWSSecretAccessKey != "" {
		detail := fmt.Sprintf("aws_access_key_id=%s", badge.AWSAccessKeyID)

		if badge.AWSRegion != "" {
			detail = fmt.Sprintf("%s, region=%s", detail, badge.AWSRegion)
		}

		return AWSStaticCredentialsType, detail, nil
	} else {
		return "", "", nil
	}
}
