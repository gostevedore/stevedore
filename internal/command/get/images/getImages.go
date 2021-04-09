package getimages

import (
	"context"

	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/engine"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/cobra"
)

type getImagesCmdFlags struct {
	Tree bool
	// Plantuml bool
}

var getImagesCmdFlagsVar *getImagesCmdFlags

// NewCommand return an stevedore command object for get images
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {
	getImagesCmdFlagsVar = &getImagesCmdFlags{}

	getImagesCmd := &cobra.Command{
		Use: "images",
		Aliases: []string{
			"images",
		},
		Short: "Stevedore subcommand to get images information",
		Long: `Stevedore subcommand to get images information

  Example:
    stevedore get images --tree
`,
		Args: cobra.MaximumNArgs(0),
		RunE: getImagesHandler(ctx, config),
	}

	getImagesCmd.Flags().BoolVarP(&getImagesCmdFlagsVar.Tree, "tree", "t", false, "Return the output as a tree")
	// TODO: pending to be scheduled
	// getImagesCmd.Flags().BoolVarP(&getImagesCmdFlagsVar.Plantuml, "plantuml", "p", false, "Return the output as a Plantuml graph")

	command := &command.StevedoreCommand{
		Command: getImagesCmd,
	}

	return command
}

func getImagesHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {

	return func(cmd *cobra.Command, args []string) error {

		var err error
		var imagesEngine *engine.ImagesEngine

		imagesEngine, err = engine.NewImagesEngine(ctx, 1, config.TreePathFile, config.BuilderPathFile)
		if err != nil {
			return errors.New("(command::getImagesHandler)", "Error creating images engine", err)
		}

		if getImagesCmdFlagsVar.Tree {
			imagesEngine.DrawGraph(ctx)
			return nil
		}
		// TODO: pending to be scheduled
		// if getImagesCmdFlagsVar.Plantuml {
		// 	return nil
		// }

		return listImages(ctx, imagesEngine)
	}
}

func listImages(ctx context.Context, imagesEngine *engine.ImagesEngine) error {
	var table [][]string

	if imagesEngine == nil {
		return errors.New("(listImages)", "Images engine is nil")
	}

	listImages, err := imagesEngine.ListImages()
	if err != nil {
		return errors.New("(command::listImages)", "Error listing registry credentials", err)
	}

	table = make([][]string, len(listImages)+1)
	table[0] = engine.ListImageHeader()
	copy(table[1:], listImages)

	console.PrintTable(table)

	return nil

}
