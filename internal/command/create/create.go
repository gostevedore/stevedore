package create

import (
	"context"
	"stevedore/internal/command"
	cmdconf "stevedore/internal/command/create/configuration"
	"stevedore/internal/command/create/credentials"
	"stevedore/internal/command/middleware"
	"stevedore/internal/configuration"

	"github.com/spf13/cobra"
)

//  NewCommand return an stevedore command object for get
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"generate"},
		Short:   "Create stevedore elements",
		Long:    "",
		RunE:    createHandler(ctx),
	}

	command := &command.StevedoreCommand{
		Command: createCmd,
	}

	command.AddCommand(middleware.Middleware(credentials.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(cmdconf.NewCommand(ctx, config)))

	return command
}

func createHandler(ctx context.Context) command.CobraRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
