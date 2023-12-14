package credentials

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

// Output is an output for the builders
type Output struct {
	write   OutputWriter
	methods []Outputter
}

// NewOutput creates a new Output
func NewOutput(write OutputWriter, output ...Outputter) *Output {
	return &Output{
		write:   write,
		methods: output,
	}
}

// Output prints the credentials
func (o *Output) Print(credentials []*credentials.Credential) error {

	errContext := "(output::credentials::Output::PrintAll)"

	if o.write == nil {
		return errors.New(errContext, "To print credentials, you must provide a writer")
	}

	content := [][]string{}
	content = append(content, []string{"ID", "TYPE", "CREDENTIALS"})

	for _, credential := range credentials {
		for _, method := range o.methods {
			credentialsType, detail, err := method.Output(credential)
			if err != nil {
				continue
			}

			if detail != "" && credentialsType != "" {
				content = append(content, []string{credential.ID, credentialsType, detail})
				break
			}
		}
	}

	err := o.write.PrintTable(content)
	if err != nil {
		return errors.New(errContext, "error printing credentials.", err)
	}

	return nil
}
