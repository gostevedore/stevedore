package builders

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	handler "github.com/gostevedore/stevedore/internal/handler/get/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration, entrypoint Entrypointer) *command.StevedoreCommand {

	getBuildersFlagOptions := &getBuildersFlagOptions{}

	getBuildersCmd := &cobra.Command{
		Use: "builders",
		Aliases: []string{
			"builder",
			"b",
		},
		Short: "Stevedore subcommand to get builders information",
		Long: `
Stevedore subcommand to get builders information
`,
		Example: `
Get builder filtered by name:
  stevedore get images --filter name=golang-app

Get builder filtered by driver:
  stevedore get images --filter driver=docker
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::get::builders::RunE)"

			handlerOptions := &handler.Options{}

			if len(getBuildersFlagOptions.Filter) > 0 {
				handlerOptions.Filter = append([]string{}, getBuildersFlagOptions.Filter...)
			}

			err = entrypoint.Execute(ctx, cmd.Flags().Args(), config, handlerOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	getBuildersCmd.Flags().StringSliceVarP(&getBuildersFlagOptions.Filter, "filter", "f", []string{}, "List of filters to apply. Filters must be defined on the following format: <attribute>=<value>")

	command := &command.StevedoreCommand{
		Command: getBuildersCmd,
	}

	return command
}
