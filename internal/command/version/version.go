package version

import (
	"context"
	"stevedore/internal/command"
	"stevedore/internal/configuration"
	"stevedore/internal/release"
	"stevedore/internal/ui/console"

	"github.com/spf13/cobra"
)

// NewCommmand creates a StevedoreCommand for version
func NewCommand(ctx context.Context, conf *configuration.Configuration) *command.StevedoreCommand {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print Stevedore version number",
		Long:  `Stevedore version, the docker images orchestrator`,
		RunE:  versionHandler(ctx, conf),
	}

	command := &command.StevedoreCommand{
		Configuration: conf,
		Command:       versionCmd,
	}

	return command
}

func versionHandler(ctx context.Context, conf *configuration.Configuration) command.CobraRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		r := release.NewRelease(console.GetConsole())
		r.PrintVersion()
		return nil
	}
}
