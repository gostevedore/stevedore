package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/credentials"
	handler "github.com/gostevedore/stevedore/internal/handler/create/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

const (
	DeprecatedFlagMessageRegistryHost                 = "[DEPRECATED FLAG] credentials id must be passed as command argument instead of using 'registry-host' flag"
	DeprecatedFlagMessageDockerRegistryCredentialsDir = "[DEPRECATED FLAG] 'credentials-dir' is deprecated and will be ignored. Credentials parameters are set through the 'credentials' section of the configuration file or using the flag 'local-storage-path'"
)

// NewCommand return an stevedore command object to create credentials
func NewCommand(ctx context.Context, compatibility Compatibilitier, config *configuration.Configuration, e Entrypointer) *command.StevedoreCommand {

	createCredentialsFlagOptions := &createCredentialsFlagOptions{}

	createCredentialsCmd := &cobra.Command{
		Use: "credentials",
		Aliases: []string{
			"auth",
			"badge",
		},
		Short: "Stevedore subcommand to add a new credentials badge into credentials store",
		Long: `
		Stevedore subcommand to add a new credentials badge into credentials store
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::create::credentials::RunE)"

			handlerOptions := &handler.Options{}
			entrypointOptions := &entrypoint.Options{}

			if createCredentialsFlagOptions.LocalStoragePath != "" {
				entrypointOptions.LocalStoragePath = createCredentialsFlagOptions.LocalStoragePath
			}
			if createCredentialsFlagOptions.Force {
				entrypointOptions.ForceCreate = createCredentialsFlagOptions.Force
			}

			if createCredentialsFlagOptions.AllowUseSSHAgent {
				handlerOptions.AllowUseSSHAgent = createCredentialsFlagOptions.AllowUseSSHAgent
			}
			if createCredentialsFlagOptions.AWSAccessKeyID != "" {
				handlerOptions.AWSAccessKeyID = createCredentialsFlagOptions.AWSAccessKeyID
			}
			if createCredentialsFlagOptions.AWSProfile != "" {
				handlerOptions.AWSProfile = createCredentialsFlagOptions.AWSProfile
			}
			if createCredentialsFlagOptions.AWSRegion != "" {
				handlerOptions.AWSRegion = createCredentialsFlagOptions.AWSRegion
			}
			if createCredentialsFlagOptions.AWSRoleARN != "" {
				handlerOptions.AWSRoleARN = createCredentialsFlagOptions.AWSRoleARN
			}
			if len(createCredentialsFlagOptions.AWSSharedConfigFiles) > 0 {
				handlerOptions.AWSSharedConfigFiles = append([]string{}, createCredentialsFlagOptions.AWSSharedConfigFiles...)
			}
			if len(createCredentialsFlagOptions.AWSSharedCredentialsFiles) > 0 {
				handlerOptions.AWSSharedCredentialsFiles = append([]string{}, createCredentialsFlagOptions.AWSSharedCredentialsFiles...)
			}
			if createCredentialsFlagOptions.AWSUseDefaultCredentialsChain {
				handlerOptions.AWSUseDefaultCredentialsChain = createCredentialsFlagOptions.AWSUseDefaultCredentialsChain
			}
			if createCredentialsFlagOptions.GitSSHUser != "" {
				handlerOptions.GitSSHUser = createCredentialsFlagOptions.GitSSHUser
			}
			if createCredentialsFlagOptions.PrivateKeyFile != "" {
				handlerOptions.PrivateKeyFile = createCredentialsFlagOptions.PrivateKeyFile
			}
			if createCredentialsFlagOptions.PrivateKeyPassword != "" {
				handlerOptions.PrivateKeyPassword = createCredentialsFlagOptions.PrivateKeyPassword
			}
			if createCredentialsFlagOptions.Username != "" {
				handlerOptions.Username = createCredentialsFlagOptions.Username
			}

			if createCredentialsFlagOptions.DEPRECATEDRegistryHost != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageRegistryHost)
				entrypointOptions.DEPRECATEDRegistryHost = createCredentialsFlagOptions.DEPRECATEDRegistryHost
			}

			if createCredentialsFlagOptions.DEPRECATEDDockerRegistryCredentialsDir != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageDockerRegistryCredentialsDir)
				entrypointOptions.LocalStoragePath = createCredentialsFlagOptions.DEPRECATEDDockerRegistryCredentialsDir
			}

			err = e.Execute(ctx, cmd.Flags().Args(), config, entrypointOptions, handlerOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.AllowUseSSHAgent, "allow-use-ssh-agent", false, "When is used that flag, is allowed to use ssh-agent")
	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.AWSUseDefaultCredentialsChain, "aws-use-default-credentials-chain", false, "When is used that flag, AWS default credentials chain is used to achieve credentials from AWS")
	createCredentialsCmd.Flags().StringSliceVar(&createCredentialsFlagOptions.AWSSharedConfigFiles, "aws-shared-config-files", []string{}, "List of AWS shared config files to achieve credentials from AWS")
	createCredentialsCmd.Flags().StringSliceVar(&createCredentialsFlagOptions.AWSSharedCredentialsFiles, "aws-shared-credentials-files", []string{}, "List AWS shared credentials files to achieve credentials from AWS")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSAccessKeyID, "aws-access-key-id", "", "AWS Access Key ID to achieve credentials from AWS to achieve credentials from AWS. AWS Secret asked key is going to be requested")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSProfile, "aws-profile", "", "AWS Profile to achieve credentials from AWS")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSRegion, "aws-region", "", "AWS Region to achieve credentials from AWS")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSRoleARN, "aws-role-arn", "", "AWS Role ARN to achieve credentials from AWS")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.GitSSHUser, "git-ssh-user", "", "Git SSH User")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.LocalStoragePath, "local-storage-path", "", "Path where credentials are stored locally, using local storage type")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.PrivateKeyFile, "private-key-file", "", "Private Key File")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.PrivateKeyPassword, "private-key-password", "", "Private Key Password")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.Username, "username", "", "Username for basic auth method. Password is going to be requested")
	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.Force, "force", false, "When is enabled the flag, credentials creation is forced. It overwrites the existing value")

	createCredentialsCmd.Flags().StringVarP(&createCredentialsFlagOptions.DEPRECATEDDockerRegistryCredentialsDir, "credentials-dir", "d", "", DeprecatedFlagMessageDockerRegistryCredentialsDir)
	createCredentialsCmd.Flags().StringVarP(&createCredentialsFlagOptions.DEPRECATEDRegistryHost, "registry-host", "r", "", DeprecatedFlagMessageRegistryHost)

	command := &command.StevedoreCommand{
		Command: createCredentialsCmd,
	}

	return command
}
