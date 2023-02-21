package awsstaticcredentials

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	errContext := "(output::credentials::types::AWSStaticCredentialsOutput::Output)"

	tests := []struct {
		desc            string
		output          *AWSStaticCredentialsOutput
		badge           *credentials.Badge
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for AWSStaticCredentialsOutput and badge is nil",
			output:          NewAWSStaticCredentialsOutput(),
			badge:           nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "To show badge output, badge must be provided"),
		},
		{
			desc:   "Testing generate output for AWSStaticCredentialsOutput",
			output: NewAWSStaticCredentialsOutput(),
			badge: &credentials.Badge{
				AWSAccessKeyID:     "accessKeyID",
				AWSSecretAccessKey: "secretAccessKey",
				AWSRegion:          "region",
			},
			detail:          "aws_access_key_id=accessKeyID, region=region",
			credentialsType: AWSStaticCredentialsType,
			err:             &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			credentialsType, detail, err := test.output.Output(test.badge)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.credentialsType, credentialsType)
				assert.Equal(t, test.detail, detail)
			}
		})
	}
}
