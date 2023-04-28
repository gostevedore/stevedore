package usernamepassword

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

const (
	// UsernamePasswordType is the name of the basic authentication type
	UsernamePasswordType = "username-password"
)

type UsernamePasswordOutput struct{}

func NewUsernamePasswordOutput() *UsernamePasswordOutput {
	return &UsernamePasswordOutput{}
}

func (o *UsernamePasswordOutput) Output(credential *credentials.Credential) (string, string, error) {

	errContext := "(output::credentials::types::UsernamePasswordOutput::Output)"

	if credential == nil {
		return "", "", errors.New(errContext, "To show credential output, credential must be provided")
	}

	if credential.Username != "" && credential.Password != "" {
		return UsernamePasswordType, fmt.Sprintf("username=%s", credential.Username), nil
	} else {
		return "", "", nil
	}
}
