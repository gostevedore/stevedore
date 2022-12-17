package configuration

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration, entrypoint Entrypointer) *command.StevedoreCommand {

	getConfigurationCmd := &cobra.Command{
		Use: "configuration",
		Aliases: []string{
			"config",
			"conf",
			"cfg",
		},
		Short: "Stevedore subcommand to get configuration information",
		Long: `
		Stevedore subcommand to get configuration information
`,
		Example: `
  stevedore get configuration
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::get::configuration::RunE)"

			err = entrypoint.Execute(ctx, cmd.Flags().Args(), config)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	command := &command.StevedoreCommand{
		Command: getConfigurationCmd,
	}

	return command
}
