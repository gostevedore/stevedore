package version

import (
	"context"
	"os"

	domainrelease "github.com/gostevedore/stevedore/internal/core/domain/release"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/gostevedore/stevedore/internal/infrastructure/output/release"
	"github.com/spf13/cobra"
)

// NewCommmand creates a StevedoreCommand for version
func NewCommand(ctx context.Context, console Consoler) *command.StevedoreCommand {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Stevedore command to print the binary release version",
		Long:  "Stevedore command to print the binary release version",
		RunE:  versionHandler(ctx),
	}

	command := &command.StevedoreCommand{
		Command: versionCmd,
	}

	return command
}

func versionHandler(ctx context.Context) command.CobraRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var r *domainrelease.Release
		var p *release.Output
		cons := console.NewConsole(os.Stdout, nil)
		r = domainrelease.NewRelease()
		p = release.NewOutput(cons)
		err := p.Print(r)
		if err != nil {
			return err
		}

		return nil
	}
}
