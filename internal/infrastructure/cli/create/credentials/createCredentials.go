package credentials

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

//  NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration, entrypoint Entrypointer) *command.StevedoreCommand {

	createCredentialsCmd := &cobra.Command{
		Use: "credentials",
		Short: "Stevedore subcommand",
		Long: `
Stevedore subcommand
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(create/credentials::RunE)"

			err = entrypoint.Execute(ctx, cmd.Flags().Args(), config)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	command := &command.StevedoreCommand{
		Command: createCredentialsCmd,
	}

	return command
}
