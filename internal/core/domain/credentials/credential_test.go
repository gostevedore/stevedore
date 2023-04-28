package credentials

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	errContext := "(core::domain::credentials::IsValid)"

	tests := []struct {
		desc       string
		credential *Credential
		valid      bool
		err        error
	}{
		{
			desc: "Testing a valid credential with username and password",
			credential: &Credential{
				Username: "username",
				Password: "password",
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid credential with aws access key",
			credential: &Credential{
				AWSAccessKeyID:     "id",
				AWSSecretAccessKey: "secret",
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid credential with aws default credentials chain",
			credential: &Credential{
				AWSUseDefaultCredentialsChain: true,
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid credential with private key file",
			credential: &Credential{
				PrivateKeyFile: "file",
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid credential with ssh agent",
			credential: &Credential{
				AllowUseSSHAgent: true,
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc:       "Testing an invalid credential",
			credential: &Credential{},
			valid:      true,
			err:        errors.New(errContext, "Invalid credential. Unknown reason"),
		},
		{
			desc: "Testing an invalid credential with user provided and password not provided",
			credential: &Credential{
				Username: "user",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid credential. Missing password"),
		},
		{
			desc: "Testing an invalid credential with AWS access key provided and AWS secret access key not provided",
			credential: &Credential{
				AWSAccessKeyID: "accesskey",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid credential. Missing AWS secret access key"),
		},
		{
			desc: "Testing an invalid credential with password provided and user not provided",
			credential: &Credential{
				Password: "pass",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid credential. Missing username"),
		},
		{
			desc: "Testing an invalid credential with AWS secret access key provided and AWS access key not provided",
			credential: &Credential{
				AWSSecretAccessKey: "secretaccesskey",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid credential. Missing AWS access key"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			valid, err := test.credential.IsValid()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.valid, valid)
			}

		})
	}
}
