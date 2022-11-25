package builders

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// PlainOutputOptionsFunc is a function used to configure the service
type PlainOutputOptionsFunc func(*PlainOutput)

// PlainOutput is an output for the builders
type PlainOutput struct {
	writer repository.BuildersPrinter
}

// NewPlainOutput creates a new Output
func NewPlainOutput(options ...PlainOutputOptionsFunc) *PlainOutput {
	output := &PlainOutput{}
	output.Options(options...)
	return output
}

func WithWriter(w repository.ImagesPlainPrinter) PlainOutputOptionsFunc {
	return func(o *PlainOutput) {
		o.writer = w
	}
}

// Options configure the service
func (o *PlainOutput) Options(opts ...PlainOutputOptionsFunc) {
	for _, opt := range opts {
		opt(o)
	}
}

// builderHeader returns the header for the builders
func outputHeader() []string {
	return []string{"NAME", "DRIVER"}
}

// Output writes into writer the builders from a list in plain format
func (o *PlainOutput) Output(list []*builder.Builder) error {
	errContext := "(output::builders::Output::PlainOutput)"
	content := [][]string{}
	content = append(content, outputHeader())

	if o.writer == nil {
		return errors.New(errContext, "Builders output requires a writer")
	}

	for _, builder := range list {
		builderSlice, err := builderToOutputSlice(builder)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if len(builderSlice) > 0 {
			content = append(content, builderSlice)
		}
	}

	err := o.writer.PrintTable(content)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil

}

func builderToOutputSlice(b *builder.Builder) ([]string, error) {
	res := []string{}

	if b.Name != "" && b.Driver != "" {
		res = append(res, b.Name, b.Driver)
	}

	return res, nil
}
