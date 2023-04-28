package yaml

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"gopkg.in/yaml.v3"
)

type YAMLFormater struct{}

func NewYAMLFormater() *YAMLFormater {
	return &YAMLFormater{}
}

func (f *YAMLFormater) Marshal(credential *credentials.Credential) (string, error) {

	var yamled []byte
	var err error

	errContext := "(YAMLFormater::Marshal)"

	if credential == nil {
		return "", errors.New(errContext, "Credential to be formatted must be provided")
	}

	yamled, err = yaml.Marshal(credential)
	if err != nil {
		return "", err
	}

	return string(yamled), nil
}

func (f *YAMLFormater) Unmarshal(data []byte) (*credentials.Credential, error) {

	errContext := "(YAMLFormater::Unmarshal)"

	if data == nil {
		return nil, errors.New(errContext, "Data to be unmarshalled must be provided")
	}

	credential := &credentials.Credential{}

	err := yaml.Unmarshal(data, credential)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credential, nil
}
