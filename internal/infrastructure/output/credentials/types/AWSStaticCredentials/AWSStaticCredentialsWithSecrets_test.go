package awsstaticcredentials

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutputWithSecrets(t *testing.T) {

	errContext := "(output::credentials::types::AWSStaticCredentialsWithSecretsOutput::Output)"

	tests := []struct {
		desc            string
		output          *AWSStaticCredentialsWithSecretsOutput
		badge           *credentials.Badge
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for AWSStaticCredentialsWithSecretsOutput and output is nil",
			output:          NewAWSStaticCredentialsWithSecretsOutput(nil),
			badge:           nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "AWS static credentials with secret output requieres an output"),
		},
		{
			desc: "Testing generate output for AWSStaticCredentialsWithSecretsOutput",
			output: NewAWSStaticCredentialsWithSecretsOutput(
				NewAWSStaticCredentialsOutput(),
			),
			badge: &credentials.Badge{
				AWSAccessKeyID:     "accessKeyID",
				AWSSecretAccessKey: "secretAccessKey",
				AWSRegion:          "region",
			},
			detail:          "aws_access_key_id=accessKeyID, region=region, aws_secret_access_key=secretAccessKey",
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
