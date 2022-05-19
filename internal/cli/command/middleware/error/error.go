package error

import (
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

type ErrorManagement struct {
	cmd     *command.StevedoreCommand
	logger  Logger
	console Consoler
}

// NewCommand is a middleware that manages errors output on commands
func NewCommand(c *command.StevedoreCommand, p Consoler, l Logger) *command.StevedoreCommand {

	if c.Command.PersistentPreRunE != nil {
		c.Command.PersistentPreRunE = errorManagement(l, p, c.Command.PersistentPreRunE)
	}

	if c.Command.PreRunE != nil {
		c.Command.PreRunE = errorManagement(l, p, c.Command.PreRunE)
	}

	if c.Command.RunE != nil {
		c.Command.RunE = errorManagement(l, p, c.Command.RunE)
	}

	if c.Command.PostRunE != nil {
		c.Command.PostRunE = errorManagement(l, p, c.Command.PostRunE)
	}

	if c.Command.PersistentPostRunE != nil {
		c.Command.PersistentPostRunE = errorManagement(l, p, c.Command.PersistentPostRunE)
	}

	return c
}

func errorManagement(l Logger, c Consoler, f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := f(cmd, args)
		if err != nil {
			l.Error(err.(*errors.Error).ErrorWithContext())
			//c.Error(err.Error())
			c.Error(err.(*errors.Error).ErrorWithContext())
			os.Exit(1)
		}

		return nil
	}
}
