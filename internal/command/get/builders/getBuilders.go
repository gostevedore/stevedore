package getbuilders

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/build"
	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/engine"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/spf13/cobra"
)

const (
	columnSeparator = " | "
)

//  NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	getBuildersCmd := &cobra.Command{
		Use: "builders",
		Aliases: []string{
			"builder",
		},
		Short: "get builders return all builders defined",
		Long:  "get builders return all builders defined",
		RunE:  getBuildersHandler(ctx, config),
	}

	command := &command.StevedoreCommand{
		Command: getBuildersCmd,
	}

	return command
}

func getBuildersHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {

	return func(cmd *cobra.Command, args []string) error {
		var err error
		var imagesEngine *engine.ImagesEngine
		var builders [][]string
		var table [][]string

		imagesEngine, err = engine.NewImagesEngine(ctx, 1, config.TreePathFile, config.BuilderPathFile)
		if err != nil {
			return errors.New("(command::getBuildersHandler)", "Error creating images engine", err)
		}

		builders, err = imagesEngine.Builders.ListBuilders()
		if err != nil {
			return errors.New("(command::getBuildersHandler)", "Error listing builders", err)
		}

		table = make([][]string, len(builders)+1)
		table[0] = build.ListBuildersHeader()
		copy(table[1:], builders)

		console.PrintTable(table)

		return nil
	}

}
