package get

// import (
// 	"context"

// 	"github.com/gostevedore/stevedore/internal/command"
// 	getbuilders "github.com/gostevedore/stevedore/internal/command/get/builders"
// 	getconfiguration "github.com/gostevedore/stevedore/internal/command/get/configuration"
// 	getcredentials "github.com/gostevedore/stevedore/internal/command/get/credentials"
// 	getimages "github.com/gostevedore/stevedore/internal/command/get/images"
// 	getmoo "github.com/gostevedore/stevedore/internal/command/get/moo"
// 	"github.com/gostevedore/stevedore/internal/command/middleware"
// 	"github.com/gostevedore/stevedore/internal/configuration"

// 	"github.com/spf13/cobra"
// )

// type getCmdFlags struct {
// 	All bool
// }

// var getCmdFlagsVar *getCmdFlags

// //  NewCommand return an stevedore command object for get
// func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {
// 	getCmdFlagsVar = &getCmdFlags{}

// 	getCmd := &cobra.Command{
// 		Use:     "get",
// 		Aliases: []string{"list"},
// 		Short:   "Stevedore command to get items information",
// 		Long:    "Stevedore command to get items information",
// 		Run:     getHandler(ctx),
// 	}

// 	command := &command.StevedoreCommand{
// 		Command: getCmd,
// 	}

// 	// getCmd.Flags().BoolVarP(&getCmdFlagsVar.All, "all", "a", false, "Return all kind of elements")

// 	command.AddCommand(middleware.Middleware(getcredentials.NewCommand(ctx, config)))
// 	command.AddCommand(middleware.Middleware(getbuilders.NewCommand(ctx, config)))
// 	command.AddCommand(middleware.Middleware(getimages.NewCommand(ctx, config)))
// 	command.AddCommand(middleware.Middleware(getmoo.NewCommand(ctx, config)))
// 	command.AddCommand(middleware.Middleware(getconfiguration.NewCommand(ctx, config)))

// 	return command
// }

// func getHandler(ctx context.Context) command.CobraRunFunc {
// 	return func(cmd *cobra.Command, args []string) {
// 		cmd.HelpFunc()(cmd, args)
// 	}
// }
