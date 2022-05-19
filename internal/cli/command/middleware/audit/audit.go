package audit

import (
	"fmt"

	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

// NewCommand is a middleware to audit commands execution
func NewCommand(c *command.StevedoreCommand, l Logger) *command.StevedoreCommand {

	if c.Command.Run != nil {
		c.Command.Run = audit(l, c.Command.Run)
	}

	if c.Command.RunE != nil {
		c.Command.RunE = auditE(l, c.Command.RunE)
	}

	return c
}

// auditE is a function that audit commands execution
func auditE(l Logger, f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		l.Info(fmt.Sprintf("Executing command 'stevedore %s %s'", cmd.Use, args))
		return f(cmd, args)
	}
}

// audit is a function that audit commands execution
func audit(l Logger, f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		l.Info(fmt.Sprintf("Executing command 'stevedore %s %+v'", cmd.Use, args))
		f(cmd, args)
	}
}
