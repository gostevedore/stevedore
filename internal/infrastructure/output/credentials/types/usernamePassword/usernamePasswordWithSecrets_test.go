package usernamepassword

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutputWithSecret(t *testing.T) {

	errContext := "(output::credentials::types::UsernamePasswordWithSecretsOutput::Output)"

	tests := []struct {
		desc            string
		output          *UsernamePasswordWithSecretsOutput
		credential      *credentials.Credential
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for UsernamePasswordWithSecretsOutput and output is nil",
			output:          NewUsernamePasswordWithSecretsOutput(nil),
			credential:      nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "Username-password with secret output requieres an output"),
		},
		{
			desc: "Testing generate output for UsernamePasswordWithSecretsOutput",
			output: NewUsernamePasswordWithSecretsOutput(
				NewUsernamePasswordOutput(),
			),
			credential: &credentials.Credential{
				Username: "user",
				Password: "pass",
			},
			detail:          "username=user, password=pass",
			credentialsType: UsernamePasswordType,
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
