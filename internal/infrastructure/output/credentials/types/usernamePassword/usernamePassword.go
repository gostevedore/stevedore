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

func (o *UsernamePasswordOutput) Output(badge *credentials.Badge) (string, string, error) {

	errContext := "(output::credentials::types::UsernamePasswordOutput::Output)"

	if badge == nil {
		return "", "", errors.New(errContext, "To show badge output, badge must be provided")
	}

	if badge.Username != "" && badge.Password != "" {
		return UsernamePasswordType, fmt.Sprintf("username=%s", badge.Username), nil
	} else {
		return "", "", nil
	}
}
