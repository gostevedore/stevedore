package file

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/afero"
)

// OptionsFunc is a function used to configure the file console
type OptionsFunc func(*ConfigurationFilePersist)

// ConfigurationFilePersist is a configuration output which writes the configuration to a file
type ConfigurationFilePersist struct {
	filePath string
	fs       afero.Fs
}

// NewConfigurationFilePersist creates a new ConfigurationFilePersist
func NewConfigurationFilePersist(options ...OptionsFunc) *ConfigurationFilePersist {
	output := &ConfigurationFilePersist{}
	output.Options(options...)

	return output
}

func WithFileSystem(fs afero.Fs) OptionsFunc {
	return func(o *ConfigurationFilePersist) {
		o.fs = fs
	}
}

func WithFilePath(file string) OptionsFunc {
	return func(o *ConfigurationFilePersist) {
		o.filePath = file
	}
}

// Options configure the service
func (h *ConfigurationFilePersist) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(h)
	}
}

// Write writes the configuration to the writer
func (o *ConfigurationFilePersist) Write(config *configuration.Configuration) error {

	errContext := "(configuration::output::ConfigurationFilePersist::Write)"

	// configFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	configFile, err := o.fs.OpenFile(o.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("File '%s' could not be opened", o.filePath), err)
	}
	defer configFile.Close()

	err = o.write(configFile, config)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (o *ConfigurationFilePersist) write(write io.Writer, config *configuration.Configuration) error {
	var buff bytes.Buffer

	errContext := "(configuration::output::ConfigurationFilePersist::write)"
	_ = errContext

	tmpl, err := template.New("configuration").Parse(configurationTemplate)
	if err != nil {
		return errors.New(errContext, "Configuration template could not be parsed", err)
	}

	err = tmpl.Execute(&buff, config)
	if err != nil {
		return errors.New(errContext, "Error applying variables to configuration template", err)

	}

	// golang does not support some charaters on raw strings and must be reprecented by another symbols
	// "`" is reprecented by "#u0060" and must be replaced to all its occurrences
	// Though there are some templating variables which must not be replaced by parser symbols "{" and "}" are also represented by "#u007b" and "#u007b"
	replacer := strings.NewReplacer("#u0060", "`", "#u007b", "{", "#u007d", "}")
	outputConfig := replacer.Replace(buff.String())

	_, err = fmt.Fprintln(write, outputConfig)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
