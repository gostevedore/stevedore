package awsrolearn

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutputWithSecrets(t *testing.T) {

	errContext := "(output::credentials::types::AWSRoleARNWithSecretsOutput::Output)"

	tests := []struct {
		desc            string
		output          *AWSRoleARNWithSecretsOutput
		badge           *credentials.Badge
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for AWSRoleARNWithSecretsOutput and badge is nil",
			output:          NewAWSRoleARNWithSecretsOutput(nil),
			badge:           nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "AWS role arn with secret output requieres an output"),
		},
		{
			desc: "Testing generate output for AWSRoleARNWithSecretsOutput",
			output: NewAWSRoleARNWithSecretsOutput(
				NewAWSRoleARNOutput(),
			),
			badge: &credentials.Badge{
				AWSRoleARN:         "arn:aws:iam::123456789012:role/stevedore-role",
				AWSRegion:          "eu-west-1",
				AWSProfile:         "default",
				AWSAccessKeyID:     "accessKeyID",
				AWSSecretAccessKey: "secretAccessKey",
			},
			detail:          "role_arn=arn:aws:iam::123456789012:role/stevedore-role, region=eu-west-1, profile=default, aws_access_key_id=accessKeyID, aws_secret_access_key=secretAccessKey",
			credentialsType: AWSRoleARNType,
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
