package version

// import (
// 	"context"

// 	"github.com/gostevedore/stevedore/internal/command"
// 	"github.com/gostevedore/stevedore/internal/configuration"
// 	"github.com/gostevedore/stevedore/internal/release"
// 	"github.com/gostevedore/stevedore/internal/ui/console"

// 	"github.com/spf13/cobra"
// )

// // NewCommmand creates a StevedoreCommand for version
// func NewCommand(ctx context.Context, conf *configuration.Configuration) *command.StevedoreCommand {
// 	versionCmd := &cobra.Command{
// 		Use:   "version",
// 		Short: "Stevedore command to print the binary release version",
// 		Long:  "Stevedore command to print the binary release version",
// 		RunE:  versionHandler(ctx, conf),
// 	}

// 	command := &command.StevedoreCommand{
// 		Configuration: conf,
// 		Command:       versionCmd,
// 	}

// 	return command
// }

// func versionHandler(ctx context.Context, conf *configuration.Configuration) command.CobraRunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		r := release.NewRelease(console.GetConsole())
// 		r.PrintVersion()
// 		return nil
// 	}
// }
