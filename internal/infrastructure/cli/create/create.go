package create

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object to create stevedore elements
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

	return command
}
