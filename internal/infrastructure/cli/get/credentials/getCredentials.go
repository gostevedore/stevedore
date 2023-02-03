package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	getcredentialsentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object to get credentials
func NewCommand(ctx context.Context, config *configuration.Configuration, entrypoint Entrypointer) *command.StevedoreCommand {

	getCredentialsFlagOptions := &getCredentialsFlagOptions{}

	getCredentialsCmd := &cobra.Command{
		Use: "credentials",
		Aliases: []string{
			"auth",
			"auths",
			"cred",
			"creds",
			"credential",
		},
		Short: "Stevedore subcommand to get credentials information",
		Long: `Stevedore subcommand to get credentials information

  Example:
    stevedore get credentials
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::get::credentials::RunE)"

			entrypointOptions := &getcredentialsentrypoint.Options{}
			entrypointOptions.ShowSecrets = getCredentialsFlagOptions.ShowSecrets

			err = entrypoint.Execute(ctx, cmd.Flags().Args(), config, entrypointOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	getCredentialsCmd.Flags().BoolVar(&getCredentialsFlagOptions.ShowSecrets, "show-secrets", false, "When this flag is enabled, the output provide secrets")

	command := &command.StevedoreCommand{
		Command: getCredentialsCmd,
	}

	return command
}
