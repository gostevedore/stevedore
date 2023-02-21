package awsrolearn

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	errContext := "(output::credentials::types::AWSRoleARNOutput::Output)"

	tests := []struct {
		desc            string
		output          *AWSRoleARNOutput
		badge           *credentials.Badge
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for AWSRoleARNOutput and badge is nil",
			output:          NewAWSRoleARNOutput(),
			badge:           nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "To show badge output, badge must be provided"),
		},
		{
			desc:   "Testing generate output for AWSRoleARNOutput",
			output: NewAWSRoleARNOutput(),
			badge: &credentials.Badge{
				AWSRoleARN:         "arn:aws:iam::123456789012:role/stevedore-role",
				AWSRegion:          "eu-west-1",
				AWSProfile:         "default",
				AWSAccessKeyID:     "accessKeyID",
				AWSSecretAccessKey: "secretAccessKey",
			},
			detail:          "role_arn=arn:aws:iam::123456789012:role/stevedore-role, region=eu-west-1, profile=default, aws_access_key_id=accessKeyID",
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
