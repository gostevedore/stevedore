package release

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"text/template"

	errors "github.com/apenella/go-common-utils/error"
)

var (
	Version, Commit, BuildDate, OsArch string
)

const (
	header = `
           __                     __              
     _____/ /____ _   _____  ____/ /___  ________ 
    / ___/ __/ _ \ | / / _ \/ __  / __ \/ ___/ _ \
   (__  ) /_/  __/ |/ /  __/ /_/ / /_/ / /  /  __/
  /____/\__/\___/|___/\___/\__,_/\____/_/   \___/ 
                                                 
`

	versionTmpl = `Stevedore {{ .Version }} Commit: {{ .Commit }} {{ .OsArch }} BuildDate: {{ .BuildDate }}`
)

// Release
type Release struct {
	BuildDate string
	Commit    string
	Header    string
	OsArch    string
	Version   string
	Writer    io.Writer
}

// NewRelease
func NewRelease(w io.Writer) *Release {
	return &Release{
		Header:    header,
		BuildDate: strings.TrimSpace(BuildDate),
		Version:   strings.TrimSpace(Version),
		Commit:    strings.TrimSpace(Commit),
		OsArch:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Writer:    w,
	}
}

// PrintVersion show the output version
func (r *Release) PrintVersion() error {

	var w bytes.Buffer

	if r.Writer == nil {
		r.Writer = os.Stdout
	}

	tmpl, err := template.New("version").Parse(versionTmpl)
	if err != nil {
		return errors.New("(release::Version)", "Error parsing version template", err)
	}

	err = tmpl.Execute(io.Writer(&w), r)
	if err != nil {
		return errors.New("(release::Version)", "Error appling version parsed template", err)
	}

	fmt.Fprintln(r.Writer, w.String())

	return nil
}
