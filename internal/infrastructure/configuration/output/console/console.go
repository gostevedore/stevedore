package console

import (
	"io"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

type ConfigurationConsoleOutput struct {
	writer io.Writer
}

func NewConfigurationFileOutput(w io.Writer) *ConfigurationConsoleOutput {
	return &ConfigurationConsoleOutput{
		writer: w,
	}
}

func (o *ConfigurationConsoleOutput) Write(configuration *configuration.Configuration) error {
	return nil
}
