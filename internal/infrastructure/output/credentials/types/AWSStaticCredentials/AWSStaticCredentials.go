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

func (o *AWSStaticCredentialsOutput) Output(credential *credentials.Credential) (string, string, error) {

	errContext := "(output::credentials::types::AWSStaticCredentialsOutput::Output)"

	if credential == nil {
		return "", "", errors.New(errContext, "To show credential output, credential must be provided")
	}

	if credential.AWSAccessKeyID != "" && credential.AWSSecretAccessKey != "" {
		detail := fmt.Sprintf("aws_access_key_id=%s", credential.AWSAccessKeyID)

		if credential.AWSRegion != "" {
			detail = fmt.Sprintf("%s, region=%s", detail, credential.AWSRegion)
		}

		return AWSStaticCredentialsType, detail, nil
	} else {
		return "", "", nil
	}
}
