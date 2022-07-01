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

func (f *JSONFormater) Marshal(badge *credentials.Badge) (string, error) {

	var jsoned []byte
	var err error

	errContext := "(JSONFormater::Marshal)"

	if badge == nil {
		return "", errors.New(errContext, "Badge to be formatted must be provided")
	}

	jsoned, err = json.MarshalIndent(badge, "", "  ")
	if err != nil {
		return err.Error(), err
	}

	return string(jsoned), nil
}

func (f *JSONFormater) Unmarshal(data []byte) (*credentials.Badge, error) {

	errContext := "(JSONFormater::Unmarshal)"

	if data == nil {
		return nil, errors.New(errContext, "Data to be unmarshalled must be provided")
	}

	badge := &credentials.Badge{}

	err := json.Unmarshal(data, badge)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return badge, nil
}
