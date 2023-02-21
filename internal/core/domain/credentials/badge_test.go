package credentials

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	errContext := "(core::domain::credentials::IsValid)"

	tests := []struct {
		desc  string
		badge *Badge
		valid bool
		err   error
	}{
		{
			desc: "Testing a valid badge with username and password",
			badge: &Badge{
				Username: "username",
				Password: "password",
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid badge with aws access key",
			badge: &Badge{
				AWSAccessKeyID:     "id",
				AWSSecretAccessKey: "secret",
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid badge with aws default credentials chain",
			badge: &Badge{
				AWSUseDefaultCredentialsChain: true,
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid badge with private key file",
			badge: &Badge{
				PrivateKeyFile: "file",
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc: "Testing a valid badge with ssh agent",
			badge: &Badge{
				AllowUseSSHAgent: true,
			},
			valid: true,
			err:   &errors.Error{},
		},
		{
			desc:  "Testing an invalid badge",
			badge: &Badge{},
			valid: true,
			err:   errors.New(errContext, "Invalid badge. Unknown reason"),
		},
		{
			desc: "Testing an invalid badge with user provided and password not provided",
			badge: &Badge{
				Username: "user",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid badge. Missing password"),
		},
		{
			desc: "Testing an invalid badge with AWS access key provided and AWS secret access key not provided",
			badge: &Badge{
				AWSAccessKeyID: "accesskey",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid badge. Missing AWS secret access key"),
		},
		{
			desc: "Testing an invalid badge with password provided and user not provided",
			badge: &Badge{
				Password: "pass",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid badge. Missing username"),
		},
		{
			desc: "Testing an invalid badge with AWS secret access key provided and AWS access key not provided",
			badge: &Badge{
				AWSSecretAccessKey: "secretaccesskey",
			},
			valid: true,
			err:   errors.New(errContext, "Invalid badge. Missing AWS access key"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			valid, err := test.badge.IsValid()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.valid, valid)
			}

		})
	}
}
