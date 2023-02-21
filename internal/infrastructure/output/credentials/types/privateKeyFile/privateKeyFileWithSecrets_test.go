package privatekeyfile

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutputWithSecrets(t *testing.T) {

	errContext := "(output::credentials::types::PrivateKeyFileWithSecretsOutput::Output)"

	tests := []struct {
		desc            string
		output          *PrivateKeyFileWithSecretsOutput
		badge           *credentials.Badge
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for PrivateKeyFileWithSecretsOutput and output is nil",
			output:          NewPrivateKeyFileWithSecretsOutput(nil),
			badge:           nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "Private key file with secret output requieres an output"),
		},
		{
			desc: "Testing generate output for PrivateKeyFileWithSecretsOutput",
			output: NewPrivateKeyFileWithSecretsOutput(
				NewPrivateKeyFileOutput(),
			),
			badge: &credentials.Badge{
				PrivateKeyFile:     "privateKeyFile",
				PrivateKeyPassword: "privateKeyPassword",
				GitSSHUser:         "gitSSHUser",
			},
			detail:          "private_key_file=privateKeyFile, protected by password, with git user 'gitSSHUser', private_key_password=privateKeyPassword",
			credentialsType: PrivateKeyFileType,
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
