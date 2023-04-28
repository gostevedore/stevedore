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

func (o *PrivateKeyFileOutput) Output(credential *credentials.Credential) (string, string, error) {

	errContext := "(output::credentials::types::PrivateKeyFileOutput::Output)"

	if credential == nil {
		return "", "", errors.New(errContext, "To show credential output, credential must be provided")
	}

	if credential.PrivateKeyFile != "" {
		detail := fmt.Sprintf("private_key_file=%s", credential.PrivateKeyFile)

		if credential.PrivateKeyPassword != "" {
			detail = fmt.Sprintf("%s, protected by password", detail)
		}

		if credential.GitSSHUser != "" {
			detail = fmt.Sprintf("%s, with git user '%s'", detail, credential.GitSSHUser)
		}

		return PrivateKeyFileType, detail, nil
	} else {
		return "", "", nil
	}
}
