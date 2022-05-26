package release

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/release"
)

// Print is used to print the release information
type Output struct {
	writer io.Writer
}

const (
	versionTmpl = `Stevedore {{ .Version }} Commit: {{ .Commit }} {{ .OsArch }} BuildDate: {{ .BuildDate }}`
)

// NewPrint creates a new Print
func NewOutput(w io.Writer) *Output {
	return &Output{
		writer: w,
	}
}

// Print prints the release information
func (o *Output) Print(r *release.Release) error {

	var w bytes.Buffer
	errContext := "(release::Print)"

	if o.writer == nil {
		o.writer = os.Stdout
	}

	tmpl, err := template.New("version").Parse(versionTmpl)
	if err != nil {
		return errors.New(errContext, "Error parsing version template", err)
	}

	err = tmpl.Execute(io.Writer(&w), r)
	if err != nil {
		return errors.New(errContext, "Error appling version parsed template", err)
	}

	fmt.Fprintln(o.writer, w.String())

	return nil
}
