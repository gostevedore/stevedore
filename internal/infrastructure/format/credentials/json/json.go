package json

import (
	"encoding/json"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type JSONFormater struct{}

func NewJSONFormater() *JSONFormater {
	return &JSONFormater{}
}

func (f *JSONFormater) Marshal(credential *credentials.Credential) (string, error) {

	var jsoned []byte
	var err error

	errContext := "(JSONFormater::Marshal)"

	if credential == nil {
		return "", errors.New(errContext, "Credential to be formatted must be provided")
	}

	jsoned, err = json.MarshalIndent(credential, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsoned), nil
}

func (f *JSONFormater) Unmarshal(data []byte) (*credentials.Credential, error) {

	errContext := "(JSONFormater::Unmarshal)"

	if data == nil {
		return nil, errors.New(errContext, "Data to be unmarshalled must be provided")
	}

	credential := &credentials.Credential{}

	err := json.Unmarshal(data, credential)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credential, nil
}
