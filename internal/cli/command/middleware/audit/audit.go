package audit

import (
	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

// NewCommand is a middleware to audit commands execution
func NewCommand(c *command.StevedoreCommand, l Logger) *command.StevedoreCommand {

	if c.Command.PersistentPreRun != nil {
		persistentPreRunFunc := func(cmd *cobra.Command, args []string) {
			audit(l, cmd, args)
			c.Command.PersistentPreRun(cmd, args)
		}
		c.Command.PersistentPreRun = persistentPreRunFunc
	}

	if c.Command.PersistentPreRunE != nil {
		persistentPreRunEFunc := func(cmd *cobra.Command, args []string) error {
			audit(l, cmd, args)
			return c.Command.PersistentPreRunE(cmd, args)
		}
		c.Command.PersistentPreRunE = persistentPreRunEFunc
	}

	if c.Command.PreRun != nil {
		preRunFunc := func(cmd *cobra.Command, args []string) {
			audit(l, cmd, args)
			c.Command.PreRun(cmd, args)
		}
		c.Command.PreRun = preRunFunc
	}

	if c.Command.PreRunE != nil {
		preRunEFunc := func(cmd *cobra.Command, args []string) error {
			audit(l, cmd, args)
			return c.Command.PreRunE(cmd, args)
		}
		c.Command.PreRunE = preRunEFunc
	}

	if c.Command.Run != nil {
		runFunc := func(cmd *cobra.Command, args []string) {
			audit(l, cmd, args)
			c.Command.Run(cmd, args)
		}
		c.Command.Run = runFunc
	}

	if c.Command.RunE != nil {
		runEFunc := func(cmd *cobra.Command, args []string) error {
			audit(l, cmd, args)
			return c.Command.RunE(cmd, args)
		}
		c.Command.RunE = runEFunc
	}

	if c.Command.PostRun != nil {
		postRunFunc := func(cmd *cobra.Command, args []string) {
			audit(l, cmd, args)
			c.Command.PostRun(cmd, args)
		}
		c.Command.PostRun = postRunFunc
	}

	if c.Command.PostRunE != nil {
		postRunEFunc := func(cmd *cobra.Command, args []string) error {
			audit(l, cmd, args)
			return c.Command.PostRunE(cmd, args)
		}
		c.Command.PostRunE = postRunEFunc
	}

	if c.Command.PostRun != nil {
		postRunFunc := func(cmd *cobra.Command, args []string) {
			audit(l, cmd, args)
			c.Command.PostRun(cmd, args)
		}
		c.Command.PostRun = postRunFunc
	}

	if c.Command.PersistentPostRunE != nil {
		persistentPostRunEFunc := func(cmd *cobra.Command, args []string) error {
			audit(l, cmd, args)
			return c.Command.PersistentPostRunE(cmd, args)
		}
		c.Command.PersistentPostRunE = persistentPostRunEFunc
	}

	return c
}

func audit(l Logger, cmd *cobra.Command, args []string) {
	l.Info("Executing command 'stevedore ", cmd.Use, args, "'")
}
