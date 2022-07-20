package privatekeyfile

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

const (
	// PrivateKeyFileType is the name of the basic authentication type
	PrivateKeyFileType = "Private key file"
)

type PrivateKeyFileOutput struct{}

func NewPrivateKeyFileOutput() *PrivateKeyFileOutput {
	return &PrivateKeyFileOutput{}
}

func (o *PrivateKeyFileOutput) Output(badge *credentials.Badge) (string, string, error) {

	errContext := "(output::credentials::types::PrivateKeyFileOutput::Output)"

	if badge == nil {
		return "", "", errors.New(errContext, "To show badge output, badge must be provided")
	}

	if badge.PrivateKeyFile != "" {
		detail := fmt.Sprintf("private_key_file=%s", badge.PrivateKeyFile)

		if badge.PrivateKeyPassword != "" {
			detail = fmt.Sprintf("%s, protected by password", detail)
		}

		if badge.GitSSHUser != "" {
			detail = fmt.Sprintf("%s, with git user '%s'", detail, badge.GitSSHUser)
		}

		return PrivateKeyFileType, detail, nil
	} else {
		return "", "", nil
	}
}
