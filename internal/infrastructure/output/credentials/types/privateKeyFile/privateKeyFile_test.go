package privatekeyfile

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	errContext := "(output::credentials::types::PrivateKeyFileOutput::Output)"

	tests := []struct {
		desc            string
		output          *PrivateKeyFileOutput
		credential      *credentials.Credential
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for PrivateKeyFileOutput and credential is nil",
			output:          NewPrivateKeyFileOutput(),
			credential:      nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "To show credential output, credential must be provided"),
		},
		{
			desc:   "Testing generate output for PrivateKeyFileOutput",
			output: NewPrivateKeyFileOutput(),
			credential: &credentials.Credential{
				PrivateKeyFile:     "privateKeyFile",
				PrivateKeyPassword: "privateKeyPassword",
				GitSSHUser:         "gitSSHUser",
			},
			detail:          "private_key_file=privateKeyFile, protected by password, with git user 'gitSSHUser'",
			credentialsType: PrivateKeyFileType,
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
