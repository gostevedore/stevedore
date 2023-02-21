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

func (f *YAMLFormater) Marshal(badge *credentials.Badge) (string, error) {

	var yamled []byte
	var err error

	errContext := "(YAMLFormater::Marshal)"

	if badge == nil {
		return "", errors.New(errContext, "Badge to be formatted must be provided")
	}

	yamled, err = yaml.Marshal(badge)
	if err != nil {
		return "", err
	}

	return string(yamled), nil
}

func (f *YAMLFormater) Unmarshal(data []byte) (*credentials.Badge, error) {

	errContext := "(YAMLFormater::Unmarshal)"

	if data == nil {
		return nil, errors.New(errContext, "Data to be unmarshalled must be provided")
	}

	badge := &credentials.Badge{}

	err := yaml.Unmarshal(data, badge)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return badge, nil
}
