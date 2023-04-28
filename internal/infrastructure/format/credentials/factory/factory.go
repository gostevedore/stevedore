package factory

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/json"
	"github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/yaml"
)

type FormatFactory struct{}

func NewFormatFactory() *FormatFactory {
	return &FormatFactory{}
}

func (f *FormatFactory) Get(format string) (repository.Formater, error) {

	errContext := "(credentials::formater::FormatFactory::Get)"

	switch format {
	case credentials.JSONFormat:
		return json.NewJSONFormater(), nil
	case credentials.YAMLFormat:
		return yaml.NewYAMLFormater(), nil
	default:
		return nil, errors.New(errContext, fmt.Sprintf("Credentials format '%s' is not supported", format))
	}
}
