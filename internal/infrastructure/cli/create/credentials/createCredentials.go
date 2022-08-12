package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

const (
	DeprecatedFlagMessageRegistryHost      = "[DEPRECATED FLAG] use 'credentials-id' instead of 'registry-host'"
	DeprecatedDockerRegistryCredentialsDir = "[DEPRECATED FLAG] 'credentials-dir' is deprecated and will be ignored. Credentials parameters are set through the 'credentials' section of the configuration file"
)

//  NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration, entrypoint Entrypointer) *command.StevedoreCommand {

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

			err = entrypoint.Execute(ctx, cmd.Flags().Args(), config)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.Username, "username", "", "Docker registry username")
	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.AskPassword, "password", false, "Docker registry password")
	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.AllowUseSSHAgent, "allow-use-ssh-agent", false, "Allow use of ssh-agent")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSAccessKeyID, "aws-access-key-id", "", "AWS Access Key ID")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSProfile, "aws-profile", "", "AWS Profile")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSRegion, "aws-region", "", "AWS Region")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.AWSRoleARN, "aws-role-arn", "", "AWS Role ARN")
	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.AskAWSSecretAccessKey, "aws-secret-access-key", false, "AWS Secret Access Key")
	createCredentialsCmd.Flags().StringSliceVar(&createCredentialsFlagOptions.AWSSharedConfigFiles, "aws-shared-config-files", []string{}, "AWS Shared Config Files")
	createCredentialsCmd.Flags().StringSliceVar(&createCredentialsFlagOptions.AWSSharedCredentialsFiles, "aws-shared-credentials-files", []string{}, "AWS Shared Credentials Files")
	createCredentialsCmd.Flags().BoolVar(&createCredentialsFlagOptions.AWSUseDefaultCredentialsChain, "aws-use-default-credentials-chain", false, "AWS Use Default Credentials Chain")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.GitSSHUser, "git-ssh-user", "", "Git SSH User")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.PrivateKeyFile, "private-key-file", "", "Private Key File")
	createCredentialsCmd.Flags().StringVar(&createCredentialsFlagOptions.PrivateKeyPassword, "private-key-password", "", "Private Key Password")

	createCredentialsCmd.Flags().StringVarP(&createCredentialsFlagOptions.DEPRECATEDDockerRegistryCredentialsDir, "credentials-dir", "d", "", DeprecatedDockerRegistryCredentialsDir)
	createCredentialsCmd.Flags().StringVarP(&createCredentialsFlagOptions.DEPRECATEDRegistryHost, "registry-host", "r", "", "Docker registry host to register credentials")

	command := &command.StevedoreCommand{
		Command: createCredentialsCmd,
	}

	return command
}
