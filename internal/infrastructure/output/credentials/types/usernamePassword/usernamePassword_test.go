package usernamepassword

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	errContext := "(output::credentials::types::UsernamePasswordOutput::Output)"

	tests := []struct {
		desc            string
		output          *UsernamePasswordOutput
		credential      *credentials.Credential
		detail          string
		credentialsType string
		err             error
	}{
		{
			desc:            "Testing error when creating the output for UsernamePasswordOutput and credential is nil",
			output:          NewUsernamePasswordOutput(),
			credential:      nil,
			detail:          "",
			credentialsType: "",
			err:             errors.New(errContext, "To show credential output, credential must be provided"),
		},
		{
			desc:   "Testing generate output for UsernamePasswordOutput",
			output: NewUsernamePasswordOutput(),
			credential: &credentials.Credential{
				Username: "user",
				Password: "pass",
			},
			detail:          "username=user",
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
