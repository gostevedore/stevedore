package file

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

// ConfigurationFileOutput is a configuration output which writes the configuration to a file
type ConfigurationFileOutput struct {
	writer io.Writer
}

// NewConfigurationFileOutput creates a new ConfigurationFileOutput
func NewConfigurationFileOutput(w io.Writer) *ConfigurationFileOutput {
	return &ConfigurationFileOutput{
		writer: w,
	}
}

// Write writes the configuration to the writer
func (o *ConfigurationFileOutput) Write(config *configuration.Configuration) error {

	var buff bytes.Buffer

	errContext := "(configuration::output::file::CreateConfigurationFile)"

	tmpl, err := template.New("configuration").Parse(configurationTemplate)
	if err != nil {
		return errors.New(errContext, "Configuration template could not be parsed", err)
	}

	err = tmpl.Execute(&buff, config)
	if err != nil {
		return errors.New("(configuration::CreateConfigurationFile)", "Error applying variables to configuration template", err)

	}

	// golang does not support some charaters on raw strings and must be reprecented by another symbols
	// "`" is reprecented by "#u0060" and must be replaced to all its occurrences
	// Though there are some templating variables which must not be replaced by parser symbols "{" and "}" are also represented by "#u007b" and "#u007b"
	replacer := strings.NewReplacer("#u0060", "`", "#u007b", "{", "#u007d", "}")
	outputConfig := replacer.Replace(buff.String())

	fmt.Fprintln(o.writer, outputConfig)

	return nil
}
