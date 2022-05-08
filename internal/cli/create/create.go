package create

// import (
// 	"context"

// 	"github.com/gostevedore/stevedore/internal/cli/command"
// 	cmdconf "github.com/gostevedore/stevedore/internal/cli/create/configuration"
// 	"github.com/gostevedore/stevedore/internal/command/create/credentials"
// 	"github.com/gostevedore/stevedore/internal/command/middleware"
// 	"github.com/gostevedore/stevedore/internal/configuration"

// 	"github.com/spf13/cobra"
// )

// //  NewCommand return an stevedore command object for get
// func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

// 	createCmd := &cobra.Command{
// 		Use:     "create",
// 		Aliases: []string{"generate"},
// 		Short:   "Stevedore command to create configuration files",
// 		Long:    "Stevedore command to create configuration files",
// 		RunE:    createHandler(ctx),
// 	}

// 	command := &command.StevedoreCommand{
// 		Command: createCmd,
// 	}

// 	command.AddCommand(middleware.Middleware(credentials.NewCommand(ctx, config)))
// 	command.AddCommand(middleware.Middleware(cmdconf.NewCommand(ctx, config)))

// 	return command
// }

// func createHandler(ctx context.Context) command.CobraRunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		cmd.HelpFunc()(cmd, args)
// 		return nil
// 	}
// }
