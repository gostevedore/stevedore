package moo

import (
	"context"
	"fmt"
	"stevedore/internal/command"
	"stevedore/internal/configuration"
	"stevedore/internal/ui/console"

	"github.com/spf13/cobra"
)

const moo = `
               (__)
               (oo)
         /------\/
        / |    ||
       *  /\---/\
          ~~   ~~
..."Have you mooed today?"...`

// NewCommmand creates a StevedoreCommand for version
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {
	versionCmd := &cobra.Command{
		Use:    "moo",
		Short:  "Have you mooed today?",
		Long:   `Have you mooed today?`,
		Hidden: true,
		Run:    mooHandler(ctx),
	}

	command := &command.StevedoreCommand{
		Command: versionCmd,
	}

	return command
}

func mooHandler(ctx context.Context) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(console.GetConsole(), moo)
	}
}
