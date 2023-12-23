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

// ConfigurationFileSafePersist is a configuration output which writes the configuration to a file
type ConfigurationFileSafePersist struct {
	ConfigurationFilePersist
}

// NewConfigurationFileSafePersist creates a new ConfigurationFileSafePersist
func NewConfigurationFileSafePersist(options ...OptionsFunc) *ConfigurationFileSafePersist {
	output := &ConfigurationFileSafePersist{}
	output.Options(options...)

	return output
}

// Write writes the configuration to the writer
func (o *ConfigurationFileSafePersist) Write(config *configuration.Configuration) (err error) {
	var configFile afero.File

	errContext := "(configuration::output::ConfigurationFileSafePersist::Write)"

	fileInfo, _ := o.fs.Stat(o.filePath)

	if fileInfo != nil {
		return errors.New(errContext, fmt.Sprintf("Configuration file '%s' already exist and will not be created", o.filePath))
	}

	configFile, err = o.fs.OpenFile(o.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("File '%s' could not be opened", o.filePath), err)
	}

	defer func() {
		closeFileErr := configFile.Close()
		if closeFileErr != nil {
			// here the closeFileErr is appended to the err returned by the function. With that we ensure that the closeFileErr is not lost
			err = errors.New(errContext, fmt.Sprintf("Error closing file '%s'.", o.filePath), closeFileErr, err)
		}
	}()

	err = o.write(configFile, config)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return err
}

func (o *ConfigurationFileSafePersist) write(write io.Writer, config *configuration.Configuration) error {
	var buff bytes.Buffer

	errContext := "(configuration::output::ConfigurationFileSafePersist::write)"
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
