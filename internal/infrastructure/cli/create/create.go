package create

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/spf13/cobra"
)

// import (
// 	"context"

// 	"github.com/gostevedore/stevedore/internal/cli/command"
// 	cmdconf "github.com/gostevedore/stevedore/internal/cli/create/configuration"
// 	"github.com/gostevedore/stevedore/internal/command/create/credentials"
// 	"github.com/gostevedore/stevedore/internal/command/middleware"
// 	"github.com/gostevedore/stevedore/internal/configuration"

// 	"github.com/spf13/cobra"
// )

//  NewCommand return an stevedore command object for get
func NewCommand(ctx context.Context, subcommands ...*command.StevedoreCommand) *command.StevedoreCommand {

	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"generate"},
		Short:   "Stevedore command to create items",
		Long:    "Stevedore command to create items",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command := &command.StevedoreCommand{
		Command: createCmd,
	}

	for _, subcommand := range subcommands {
		command.AddCommand(subcommand)
	}
	// command.AddCommand(middleware.Middleware(credentials.NewCommand(ctx, config)))
	// command.AddCommand(middleware.Middleware(cmdconf.NewCommand(ctx, config)))

	return command
}

// func createHandler(ctx context.Context) command.CobraRunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		cmd.HelpFunc()(cmd, args)
// 		return nil
// 	}
// }
