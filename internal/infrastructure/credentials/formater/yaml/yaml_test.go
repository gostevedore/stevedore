package yaml

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	errContext := "(YAMLFormater::Marshal)"
	tests := []struct {
		desc     string
		formater *YAMLFormater
		badge    *credentials.Badge
		res      string
		err      error
	}{
		{
			desc:     "Testing error when formating a nil badge to YAML",
			formater: NewYAMLFormater(),
			badge:    nil,
			err:      errors.New(errContext, "Badge to be formatted must be provided"),
		},
		{
			desc:     "Testing formating a badge to YAML",
			formater: NewYAMLFormater(),
			badge: &credentials.Badge{
				AWSAccessKeyID:                "awsaccesskeyid",
				AWSRegion:                     "awsregion",
				AWSRoleARN:                    "awsrolearn",
				AWSSecretAccessKey:            "awssecretaccesskey",
				AWSProfile:                    "awsprofile",
				AWSSharedCredentialsFiles:     []string{"awssharedcredentialsfiles"},
				AWSSharedConfigFiles:          []string{"awssharedconfigfiles"},
				AWSUseDefaultCredentialsChain: true,
				DEPRECATEDPassword:            "deprecatedpassword",
				DEPRECATEDUsername:            "deprecatedusername",
				Password:                      "password",
				Username:                      "username",
				PrivateKeyFile:                "privatekeyfile",
				PrivateKeyPassword:            "privatekeypassword",
				GitSSHUser:                    "gitsshuser",
				AllowUseSSHAgent:              true,
			},
			res: `id: ""
aws_access_key_id: awsaccesskeyid
aws_region: awsregion
aws_role_arn: awsrolearn
aws_secret_access_key: awssecretaccesskey
aws_profile: awsprofile
aws_shared_credentials_files:
    - awssharedcredentialsfiles
aws_shared_config_files:
    - awssharedconfigfiles
aws_use_default_credentials_chain: true
docker_login_password: deprecatedpassword
docker_login_username: deprecatedusername
password: password
username: username
private_key_file: privatekeyfile
private_key_password: privatekeypassword
git_ssh_user: gitsshuser
use_ssh_agent: true
`,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.formater.Marshal(test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}

}

func TestUnmashal(t *testing.T) {

	input := `
aws_access_key_id: awsaccesskeyid
aws_region: awsregion
aws_role_arn: awsrolearn
aws_secret_access_key: awssecretaccesskey
aws_profile: awsprofile
aws_shared_credentials_files:
- awssharedcredentialsfiles
aws_shared_config_files:
- awssharedconfigfiles
aws_use_default_credentials_chain: true
docker_login_password: deprecatedpassword
docker_login_username: deprecatedusername
password: password
username: username
private_key_file: privatekeyfile
private_key_password: privatekeypassword
git_ssh_user: gitsshuser
use_ssh_agent: true
`
	expected := &credentials.Badge{
		AWSAccessKeyID:                "awsaccesskeyid",
		AWSRegion:                     "awsregion",
		AWSRoleARN:                    "awsrolearn",
		AWSSecretAccessKey:            "awssecretaccesskey",
		AWSProfile:                    "awsprofile",
		AWSSharedCredentialsFiles:     []string{"awssharedcredentialsfiles"},
		AWSSharedConfigFiles:          []string{"awssharedconfigfiles"},
		AWSUseDefaultCredentialsChain: true,
		DEPRECATEDPassword:            "deprecatedpassword",
		DEPRECATEDUsername:            "deprecatedusername",
		Password:                      "password",
		Username:                      "username",
		PrivateKeyFile:                "privatekeyfile",
		PrivateKeyPassword:            "privatekeypassword",
		GitSSHUser:                    "gitsshuser",
		AllowUseSSHAgent:              true,
	}

	formater := NewYAMLFormater()
	badge, err := formater.Unmarshal([]byte(input))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, badge)

}
