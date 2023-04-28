package awsdefaultchain

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	errContext := "(output::credentials::types::AWSDefaultCredentialsChain::Output)"

	tests := []struct {
		desc            string
		output          *AWSDefaultCredentialsChainOutput
		credential      *credentials.Credential
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for AWSDefaultCredentialsChain and credential is nil",
			output:          NewAWSDefaultCredentialsChainOutput(),
			credential:      nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "To show credential output, credential must be provided"),
		},
		{
			desc:   "Testing generate output for AWSDefaultCredentialsChain",
			output: NewAWSDefaultCredentialsChainOutput(),
			credential: &credentials.Credential{
				AWSUseDefaultCredentialsChain: true,
				AWSRegion:                     "us-east-1",
				AWSProfile:                    "default",
				AWSSharedConfigFiles:          []string{"/path/to/shared/config/file"},
				AWSSharedCredentialsFiles:     []string{"/path/to/shared/credentials/file"},
			},
			detail:          "Use AWS default credentials chain, region=us-east-1, profile=default, shared_config_files=[/path/to/shared/config/file], shared_credentials_files=[/path/to/shared/credentials/file]",
			credentialsType: AWSDefaultCredentialsChainType,
			err:             &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			credentialsType, detail, err := test.output.Output(test.credential)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.credentialsType, credentialsType)
				assert.Equal(t, test.detail, detail)
			}
		})
	}
}
