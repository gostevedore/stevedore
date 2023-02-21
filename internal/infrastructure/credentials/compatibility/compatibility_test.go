package badge

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/stretchr/testify/assert"
)

func TestCheckCompatibility(t *testing.T) {

	errContext := "(credentials::compatibility::CheckCompatibility)"

	tests := []struct {
		desc              string
		compatibility     *CredentialsCompatibility
		prepareAssertFunc func(*CredentialsCompatibility)
		badge             *credentials.Badge
		res               *credentials.Badge
		err               error
	}{
		{
			desc:          "Testing error checking credentials badge compatibility when compatibilitier is not provided",
			compatibility: NewCredentialsCompatibility(nil),
			err:           errors.New(errContext, "To check credentials badge compatibility, compatibilitier must be provided"),
		},
		{
			desc: "Testing error checking credentials badge compatibility when badge is not provided",
			compatibility: NewCredentialsCompatibility(
				compatibility.NewMockCompatibility(),
			),
			badge: nil,
			err:   errors.New(errContext, "To check credentials badge compatibility, badge must be provided"),
		},
		{
			desc: "Testing check credentials badge with non incompatibilities detected",
			compatibility: NewCredentialsCompatibility(
				compatibility.NewMockCompatibility(),
			),
			badge: &credentials.Badge{
				AllowUseSSHAgent:              true,
				AWSAccessKeyID:                "access_key_id",
				AWSProfile:                    "aws_profile",
				AWSRegion:                     "aws_region",
				AWSRoleARN:                    "aws_role_arn",
				AWSSecretAccessKey:            "aws_secret_access_key",
				AWSSharedConfigFiles:          []string{"aws_share_config_files"},
				AWSSharedCredentialsFiles:     []string{"aws_share_credentials_files"},
				AWSUseDefaultCredentialsChain: true,
				GitSSHUser:                    "git_ssh_user",
				ID:                            "id",
				Password:                      "password",
				PrivateKeyFile:                "private_key_file",
				PrivateKeyPassword:            "private_key_password",
				Username:                      "user",
			},
			res: &credentials.Badge{
				AllowUseSSHAgent:              true,
				AWSAccessKeyID:                "access_key_id",
				AWSProfile:                    "aws_profile",
				AWSRegion:                     "aws_region",
				AWSRoleARN:                    "aws_role_arn",
				AWSSecretAccessKey:            "aws_secret_access_key",
				AWSSharedConfigFiles:          []string{"aws_share_config_files"},
				AWSSharedCredentialsFiles:     []string{"aws_share_credentials_files"},
				AWSUseDefaultCredentialsChain: true,
				GitSSHUser:                    "git_ssh_user",
				ID:                            "id",
				Password:                      "password",
				PrivateKeyFile:                "private_key_file",
				PrivateKeyPassword:            "private_key_password",
				Username:                      "user",
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing check credentials badge with incompatibilities detected",
			compatibility: NewCredentialsCompatibility(
				compatibility.NewMockCompatibility(),
			),
			badge: &credentials.Badge{
				AllowUseSSHAgent:              true,
				AWSAccessKeyID:                "access_key_id",
				AWSProfile:                    "aws_profile",
				AWSRegion:                     "aws_region",
				AWSRoleARN:                    "aws_role_arn",
				AWSSecretAccessKey:            "aws_secret_access_key",
				AWSSharedConfigFiles:          []string{"aws_share_config_files"},
				AWSSharedCredentialsFiles:     []string{"aws_share_credentials_files"},
				AWSUseDefaultCredentialsChain: true,
				GitSSHUser:                    "git_ssh_user",
				ID:                            "id",
				DEPRECATEDPassword:            "password",
				PrivateKeyFile:                "private_key_file",
				PrivateKeyPassword:            "private_key_password",
				DEPRECATEDUsername:            "user",
			},
			res: &credentials.Badge{
				AllowUseSSHAgent:              true,
				AWSAccessKeyID:                "access_key_id",
				AWSProfile:                    "aws_profile",
				AWSRegion:                     "aws_region",
				AWSRoleARN:                    "aws_role_arn",
				AWSSecretAccessKey:            "aws_secret_access_key",
				AWSSharedConfigFiles:          []string{"aws_share_config_files"},
				AWSSharedCredentialsFiles:     []string{"aws_share_credentials_files"},
				AWSUseDefaultCredentialsChain: true,
				DEPRECATEDPassword:            "password",
				DEPRECATEDUsername:            "user",
				GitSSHUser:                    "git_ssh_user",
				ID:                            "id",
				Password:                      "password",
				PrivateKeyFile:                "private_key_file",
				PrivateKeyPassword:            "private_key_password",
				Username:                      "user",
			},
			prepareAssertFunc: func(c *CredentialsCompatibility) {
				c.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'docker_login_username' is deprecated and will be removed on v0.12.0, please use 'username' instead"}).Return(nil)
				c.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'docker_login_password' is deprecated and will be removed on v0.12.0, please use 'password' instead"}).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.compatibility)
			}

			err := test.compatibility.CheckCompatibility(test.badge)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, test.badge)
			}
		})
	}
}
