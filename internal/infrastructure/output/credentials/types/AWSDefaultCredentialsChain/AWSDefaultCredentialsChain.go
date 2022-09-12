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

func (o *AWSDefaultCredentialsChainOutput) Output(badge *credentials.Badge) (string, string, error) {

	errContext := "(output::credentials::types::AWSDefaultCredentialsChain::Output)"

	if badge == nil {
		return "", "", errors.New(errContext, "To show badge output, badge must be provided")
	}

	if badge.AWSUseDefaultCredentialsChain {
		detail := "Use AWS default credentials chain"

		if badge.AWSRegion != "" {
			detail = fmt.Sprintf("%s, region=%s", detail, badge.AWSRegion)
		}

		if badge.AWSProfile != "" {
			detail = fmt.Sprintf("%s, profile=%s", detail, badge.AWSProfile)
		}

		if len(badge.AWSSharedConfigFiles) > 0 {
			detail = fmt.Sprintf("%s, shared_config_files=%s", detail, badge.AWSSharedConfigFiles)
		}

		if len(badge.AWSSharedCredentialsFiles) > 0 {
			detail = fmt.Sprintf("%s, shared_credentials_files=%s", detail, badge.AWSSharedCredentialsFiles)
		}

		return AWSDefaultCredentialsChainType, detail, nil
	} else {
		return "", "", nil
	}
}
