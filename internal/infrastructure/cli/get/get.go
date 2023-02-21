package get

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object to get stevedore elements
func NewCommand(ctx context.Context, subcommands ...*command.StevedoreCommand) *command.StevedoreCommand {
	// 	getCmdFlagsVar = &getCmdFlags{}

	getCmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"list"},
		Short:   "Stevedore command to get items information",
		Long:    "Stevedore command to get items information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command := &command.StevedoreCommand{
		Command: getCmd,
	}

	for _, subcommand := range subcommands {
		command.AddCommand(subcommand)
	}

	return command
}
