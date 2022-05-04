package panic

import (
	"fmt"
	"os"
	"runtime/debug"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

// PanicRecover is a middleware which recovers from panics
func NewCommand(c *command.StevedoreCommand, p Consoler, l Logger) *command.StevedoreCommand {

	if c.Command.PersistentPreRun != nil {
		persistentPreRunFunc := func(cmd *cobra.Command, args []string) {
			defer panicRecover(l, p)
			c.Command.PersistentPreRun(cmd, args)
		}
		c.Command.PersistentPreRun = persistentPreRunFunc
	}

	if c.Command.PersistentPreRunE != nil {
		persistentPreRunEFunc := func(cmd *cobra.Command, args []string) error {
			defer panicRecover(l, p)
			return c.Command.PersistentPreRunE(cmd, args)
		}
		c.Command.PersistentPreRunE = persistentPreRunEFunc
	}

	if c.Command.PreRun != nil {
		preRunFunc := func(cmd *cobra.Command, args []string) {
			defer panicRecover(l, p)
			c.Command.PreRun(cmd, args)
		}
		c.Command.PreRun = preRunFunc
	}

	if c.Command.PreRunE != nil {
		preRunEFunc := func(cmd *cobra.Command, args []string) error {
			defer panicRecover(l, p)
			return c.Command.PreRunE(cmd, args)
		}
		c.Command.PreRunE = preRunEFunc
	}

	if c.Command.Run != nil {
		runFunc := func(cmd *cobra.Command, args []string) {
			defer panicRecover(l, p)
			c.Command.Run(cmd, args)
		}
		c.Command.Run = runFunc
	}

	if c.Command.RunE != nil {
		runEFunc := func(cmd *cobra.Command, args []string) error {
			defer panicRecover(l, p)
			return c.Command.RunE(cmd, args)
		}
		c.Command.RunE = runEFunc
	}

	if c.Command.PostRun != nil {
		postRunFunc := func(cmd *cobra.Command, args []string) {
			defer panicRecover(l, p)
			c.Command.PostRun(cmd, args)
		}
		c.Command.PostRun = postRunFunc
	}

	if c.Command.PostRunE != nil {
		postRunEFunc := func(cmd *cobra.Command, args []string) error {
			defer panicRecover(l, p)
			return c.Command.PostRunE(cmd, args)
		}
		c.Command.PostRunE = postRunEFunc
	}

	if c.Command.PostRun != nil {
		postRunFunc := func(cmd *cobra.Command, args []string) {
			defer panicRecover(l, p)
			c.Command.PostRun(cmd, args)
		}
		c.Command.PostRun = postRunFunc
	}

	if c.Command.PersistentPostRunE != nil {
		persistentPostRunEFunc := func(cmd *cobra.Command, args []string) error {
			defer panicRecover(l, p)
			return c.Command.PersistentPostRunE(cmd, args)
		}
		c.Command.PersistentPostRunE = persistentPostRunEFunc
	}

	return c
}

func panicRecover(l Logger, p Consoler) {
	var errMessage string
	if err := recover(); err != nil {

		recoveredErr, isGoCommonUtilsError := err.(*errors.Error)
		if isGoCommonUtilsError {
			errMessage = recoveredErr.ErrorWithContext()
		} else {
			errMessage = recoveredErr.Error()
		}

		msg := fmt.Sprintf("Unexpected panic: %s", errMessage)
		l.Error(msg)
		l.Error(string(debug.Stack()))
		p.Error(msg)
		os.Exit(1)
	}
}
