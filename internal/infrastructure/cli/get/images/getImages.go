package images

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/images"
	handler "github.com/gostevedore/stevedore/internal/handler/get/images"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration, e Entrypointer) *command.StevedoreCommand {

	getImagesFlagOptions := &getImagesFlagOptions{}

	getImagesCmd := &cobra.Command{
		Use: "images",
		Aliases: []string{
			"image",
			"i",
			"img",
		},
		Short: "Stevedore subcommand that shows detail about the defined images",
		Long: `
Stevedore subcommand that shows detail about the defined images
`,
		Example: `
Get images filtered by name:
  stevedore get images --filter name=app1

Get images filtered by registry:
  stevedore get images --filter registry=registry.test
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::get::images::RunE)"

			entrypointOptions := &entrypoint.Options{}
			handlerOptions := &handler.Options{}

			entrypointOptions.Tree = getImagesFlagOptions.Tree
			entrypointOptions.UseDockerNormalizedName = getImagesFlagOptions.UseDockerNormalizedName

			if len(getImagesFlagOptions.Filter) > 0 {
				handlerOptions.Filter = append([]string{}, getImagesFlagOptions.Filter...)
			}

			err = e.Execute(ctx, cmd.Flags().Args(), config, entrypointOptions, handlerOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	getImagesCmd.Flags().BoolVarP(&getImagesFlagOptions.Tree, "tree", "t", false, "When this flag is enabled, output is returned in tree format")
	getImagesCmd.Flags().StringSliceVarP(&getImagesFlagOptions.Filter, "filter", "f", []string{}, "List of filters to apply. Filters must be defined on the following format: <attribute>=<value>")
	getImagesCmd.Flags().BoolVar(&getImagesFlagOptions.UseDockerNormalizedName, "use-docker-normalized-name", false, "Use Docker normalized name references")

	command := &command.StevedoreCommand{
		Command: getImagesCmd,
	}

	return command
}
