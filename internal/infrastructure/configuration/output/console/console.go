package console

import (
	"fmt"
	"io"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
)

type ConfigurationConsoleOutput struct {
	writer io.Writer
}

func NewConfigurationConsoleOutput(w io.Writer) *ConfigurationConsoleOutput {
	return &ConfigurationConsoleOutput{
		writer: w,
	}
}

func (o *ConfigurationConsoleOutput) Write(conf *configuration.Configuration) error {

	fmt.Println()
	fmt.Fprintf(o.writer, " %s: %s\n", configuration.BuildersPathKey, conf.BuildersPath)
	fmt.Fprintf(o.writer, " %s: %d\n", configuration.ConcurrencyKey, conf.Concurrency)
	fmt.Fprintf(o.writer, " %s: %t\n", configuration.EnableSemanticVersionTagsKey, conf.EnableSemanticVersionTags)
	fmt.Fprintf(o.writer, " %s: %s\n", configuration.ImagesPathKey, conf.ImagesPath)
	fmt.Fprintf(o.writer, " %s: %s\n", configuration.LogPathFileKey, conf.LogPathFile)
	fmt.Fprintf(o.writer, " %s: %t\n", configuration.PushImagesKey, conf.PushImages)
	if len(conf.SemanticVersionTagsTemplates) > 0 {
		fmt.Fprintf(o.writer, " %s:\n", configuration.SemanticVersionTagsTemplatesKey)
		for _, tmpl := range conf.SemanticVersionTagsTemplates {
			fmt.Fprintf(o.writer, "   - %s\n", tmpl)
		}
	}
	if conf.Credentials != nil {
		fmt.Fprintf(o.writer, " %s:\n", configuration.CredentialsKey)
		fmt.Fprintf(o.writer, "   %s: %s\n", configuration.CredentialsStorageTypeKey, conf.Credentials.StorageType)
		fmt.Fprintf(o.writer, "   %s: %s\n", configuration.CredentialsFormatKey, conf.Credentials.Format)
		if conf.Credentials.StorageType == credentials.LocalStore {
			fmt.Fprintf(o.writer, "   %s: %s\n", configuration.CredentialsLocalStoragePathKey, conf.Credentials.LocalStoragePath)
		}
	}
	fmt.Println()

	return nil
}
