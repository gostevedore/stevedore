package awsdefaultchain

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

const (
	// AWSDefaultCredentialsChainType is the name of the basic authentication method
	AWSDefaultCredentialsChainType = "AWS default Credentials chain"
)

type AWSDefaultCredentialsChainOutput struct{}

func NewAWSDefaultCredentialsChainOutput() *AWSDefaultCredentialsChainOutput {
	return &AWSDefaultCredentialsChainOutput{}
}

func (o *AWSDefaultCredentialsChainOutput) Output(credential *credentials.Credential) (string, string, error) {

	errContext := "(output::credentials::types::AWSDefaultCredentialsChain::Output)"

	if credential == nil {
		return "", "", errors.New(errContext, "To show credential output, credential must be provided")
	}

	if credential.AWSUseDefaultCredentialsChain {
		detail := "Use AWS default credentials chain"

		if credential.AWSRegion != "" {
			detail = fmt.Sprintf("%s, region=%s", detail, credential.AWSRegion)
		}

		if credential.AWSProfile != "" {
			detail = fmt.Sprintf("%s, profile=%s", detail, credential.AWSProfile)
		}

		if len(credential.AWSSharedConfigFiles) > 0 {
			detail = fmt.Sprintf("%s, shared_config_files=%s", detail, credential.AWSSharedConfigFiles)
		}

		if len(credential.AWSSharedCredentialsFiles) > 0 {
			detail = fmt.Sprintf("%s, shared_credentials_files=%s", detail, credential.AWSSharedCredentialsFiles)
		}

		return AWSDefaultCredentialsChainType, detail, nil
	} else {
		return "", "", nil
	}
}
