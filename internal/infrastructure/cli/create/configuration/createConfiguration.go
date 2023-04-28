package configuration

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object for get createConfiguration
func NewCommand(ctx context.Context, e Entrypointer) *command.StevedoreCommand {

	createConfigurationFlagOptions := &createConfigurationFlagOptions{}
	entrypointOptions := &entrypoint.Options{}

	createConfigurationCmd := &cobra.Command{
		Use: "configuration",
		Aliases: []string{
			"config",
			"conf",
			"cfg",
		},
		Short: "Stevedore subcommand to create and initialize the configuration",
		Long: `
Stevedore subcommand to create and initialize the configuration
`,
		Example: `
Example setting all configuration parameters:
  stevedore create configuration --builders-path /builders --concurrency 4 --config /stevedore-config.yaml --credentials-format json --credentials-local-storage-path /credentials --credentials-storage-type local --enable-semver-tags --force --images-path /images --log-path-file /logs --push-images --semver-tags-template "{{ .Major }}" --semver-tags-template "{{ .Major }}_{{ .Minor }}"
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::create::configuration::RunE)"

			entrypointOptions.BuildersPath = createConfigurationFlagOptions.BuildersPath
			entrypointOptions.Concurrency = createConfigurationFlagOptions.Concurrency
			entrypointOptions.ConfigurationFilePath = createConfigurationFlagOptions.ConfigurationFilePath
			entrypointOptions.CredentialsEncryptionKey = createConfigurationFlagOptions.CredentialsEncryptionKey
			entrypointOptions.CredentialsFormat = createConfigurationFlagOptions.CredentialsFormat
			entrypointOptions.CredentialsLocalStoragePath = createConfigurationFlagOptions.CredentialsLocalStoragePath
			entrypointOptions.CredentialsStorageType = createConfigurationFlagOptions.CredentialsStorageType
			entrypointOptions.EnableSemanticVersionTags = createConfigurationFlagOptions.EnableSemanticVersionTags
			entrypointOptions.Force = createConfigurationFlagOptions.Force
			entrypointOptions.GenerateCredentialsEncryptionKey = createConfigurationFlagOptions.GenerateCredentialsEncryptionKey
			entrypointOptions.ImagesPath = createConfigurationFlagOptions.ImagesPath
			entrypointOptions.LogPathFile = createConfigurationFlagOptions.LogPathFile
			entrypointOptions.PushImages = createConfigurationFlagOptions.PushImages
			entrypointOptions.SemanticVersionTagsTemplates = createConfigurationFlagOptions.SemanticVersionTagsTemplates

			err = e.Execute(ctx, entrypointOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	defaultConfiguration := configuration.DefaultConfig()

	createConfigurationCmd.Flags().StringVarP(&createConfigurationFlagOptions.BuildersPath, "builders-path", "b", defaultConfiguration.BuildersPath, fmt.Sprintf("It defines the path to locate the builders definition. Its default value is '%s'", defaultConfiguration.BuildersPath))
	createConfigurationCmd.PersistentFlags().StringVarP(&createConfigurationFlagOptions.ConfigurationFilePath, "config", "C", "", "Configuration file location path")
	createConfigurationCmd.Flags().IntVarP(&createConfigurationFlagOptions.Concurrency, "concurrency", "c", defaultConfiguration.Concurrency, fmt.Sprintf("It defines the number of concurrent workers created to build images. Its default value is '%d'", defaultConfiguration.Concurrency))
	createConfigurationCmd.Flags().StringVar(&createConfigurationFlagOptions.CredentialsEncryptionKey, "credentials-encryption-key", "", "Is the encryption key used on the credentials store")
	createConfigurationCmd.Flags().StringVar(&createConfigurationFlagOptions.CredentialsFormat, "credentials-format", defaultConfiguration.Credentials.Format, fmt.Sprintf("Format used to store credentials. The accepted formats are: %s and %s", credentials.JSONFormat, credentials.YAMLFormat))
	createConfigurationCmd.Flags().StringVar(&createConfigurationFlagOptions.CredentialsLocalStoragePath, "credentials-local-storage-path", defaultConfiguration.Credentials.LocalStoragePath, fmt.Sprintf("When is used the '%s' storage, it defines the path to store the credentials. Its default value is '%s'", credentials.LocalStore, defaultConfiguration.Credentials.LocalStoragePath))
	createConfigurationCmd.Flags().StringVar(&createConfigurationFlagOptions.CredentialsStorageType, "credentials-storage-type", defaultConfiguration.Credentials.StorageType, fmt.Sprintf("It defines the storage type. Its default value is '%s'", defaultConfiguration.Credentials.StorageType))
	createConfigurationCmd.Flags().BoolVarP(&createConfigurationFlagOptions.EnableSemanticVersionTags, "enable-semver-tags", "s", defaultConfiguration.EnableSemanticVersionTags, fmt.Sprintf("Generate extra tags when the main image tags is semver 2.0.0 compliance. Its default value is '%t'", defaultConfiguration.EnableSemanticVersionTags))
	createConfigurationCmd.Flags().BoolVar(&createConfigurationFlagOptions.Force, "force", false, "Force to create configuration file when the file already exists")
	createConfigurationCmd.Flags().BoolVar(&createConfigurationFlagOptions.GenerateCredentialsEncryptionKey, "generate-credentials-encryption-key", false, "It creates a random encryption key for the credentials store")
	createConfigurationCmd.Flags().StringVarP(&createConfigurationFlagOptions.ImagesPath, "images-path", "i", defaultConfiguration.ImagesPath, fmt.Sprintf("It defines the path to locate the images definition. Its default value is '%s'", defaultConfiguration.ImagesPath))
	createConfigurationCmd.Flags().StringVarP(&createConfigurationFlagOptions.LogPathFile, "log-path-file", "l", defaultConfiguration.LogPathFile, fmt.Sprintf("Log file location path. Its default value is '%s'", defaultConfiguration.LogPathFile))
	createConfigurationCmd.Flags().BoolVarP(&createConfigurationFlagOptions.PushImages, "push-images", "p", defaultConfiguration.PushImages, fmt.Sprintf("On build, push images automatically after it finishes. Its default value is '%t'", configuration.DefaultPushImages))
	createConfigurationCmd.Flags().StringSliceVarP(&createConfigurationFlagOptions.SemanticVersionTagsTemplates, "semver-tags-template", "t", defaultConfiguration.SemanticVersionTagsTemplates, fmt.Sprintf("List of templates which define those extra tags to generate when 'semantic_version_tags_enabled' is enabled. Its default value is '%v'", defaultConfiguration.SemanticVersionTagsTemplates))

	command := &command.StevedoreCommand{
		Command: createConfigurationCmd,
	}

	return command
}
