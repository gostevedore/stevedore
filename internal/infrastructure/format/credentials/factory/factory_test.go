package factory

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/json"
	"github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/yaml"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	errContext := "(credentials::formater::FormatFactory::Get)"

	tests := []struct {
		desc    string
		factory *FormatFactory
		format  string
		res     repository.Formater
		err     error
	}{
		{
			desc:    "Testing get JSON formater",
			factory: NewFormatFactory(),
			format:  credentials.JSONFormat,
			res:     json.NewJSONFormater(),
			err:     &errors.Error{},
		},
		{
			desc:    "Testing get YAML formater",
			factory: NewFormatFactory(),
			format:  credentials.YAMLFormat,
			res:     yaml.NewYAMLFormater(),
			err:     &errors.Error{},
		},
		{
			desc:    "Testing error getting formater when format is not supported",
			factory: NewFormatFactory(),
			format:  "invalid",
			err:     errors.New(errContext, "Credentials format 'invalid' is not supported"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			formater, err := test.factory.Get(test.format)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, test.res, formater)
			}
		})
	}
}
