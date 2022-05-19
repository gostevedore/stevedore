package panic

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

// PanicRecover is a middleware which recovers from panics
func NewCommand(c *command.StevedoreCommand, p Consoler, l Logger) *command.StevedoreCommand {

	if c.Command.PersistentPreRun != nil {
		persistentPreRunFunc := func(f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
			return func(cmd *cobra.Command, args []string) {
				defer panicRecover(l, p)
				f(cmd, args)
			}
		}
		c.Command.PersistentPreRun = persistentPreRunFunc(c.Command.PersistentPreRun)
	}

	if c.Command.PersistentPreRunE != nil {
		persistentPreRunEFunc := func(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				defer panicRecover(l, p)
				return f(cmd, args)
			}
		}
		c.Command.PersistentPreRunE = persistentPreRunEFunc(c.Command.PersistentPreRunE)
	}

	if c.Command.PreRun != nil {
		preRunFunc := func(f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
			return func(cmd *cobra.Command, args []string) {
				defer panicRecover(l, p)
				f(cmd, args)
			}
		}
		c.Command.PreRun = preRunFunc(c.Command.PreRun)
	}

	if c.Command.PreRunE != nil {
		preRunEFunc := func(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				defer panicRecover(l, p)
				return f(cmd, args)
			}
		}
		c.Command.PreRunE = preRunEFunc(c.Command.PreRunE)
	}

	if c.Command.Run != nil {
		runFunc := func(f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
			return func(cmd *cobra.Command, args []string) {
				defer panicRecover(l, p)
				f(cmd, args)
			}
		}
		c.Command.Run = runFunc(c.Command.Run)
	}

	if c.Command.RunE != nil {
		runEFunc := func(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				defer panicRecover(l, p)
				return f(cmd, args)
			}
		}
		c.Command.RunE = runEFunc(c.Command.RunE)
	}

	if c.Command.PostRun != nil {
		postRunFunc := func(f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
			return func(cmd *cobra.Command, args []string) {
				defer panicRecover(l, p)
				f(cmd, args)
			}
		}
		c.Command.PostRun = postRunFunc(c.Command.PostRun)
	}

	if c.Command.PostRunE != nil {
		postRunEFunc := func(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				defer panicRecover(l, p)
				return f(cmd, args)
			}
		}
		c.Command.PostRunE = postRunEFunc(c.Command.PostRunE)
	}

	if c.Command.PostRun != nil {
		postRunFunc := func(f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
			return func(cmd *cobra.Command, args []string) {
				defer panicRecover(l, p)
				f(cmd, args)
			}
		}
		c.Command.PostRun = postRunFunc(c.Command.PostRun)
	}

	if c.Command.PersistentPostRunE != nil {
		persistentPostRunEFunc := func(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				defer panicRecover(l, p)
				return f(cmd, args)
			}
		}
		c.Command.PersistentPostRunE = persistentPostRunEFunc(c.Command.PersistentPostRunE)
	}

	return c
}

func panicRecover(l Logger, p Consoler) {
	//var errMessage string

	if err := recover(); err != nil {

		// recoveredErr, isGoCommonUtilsError := err.(*errors.Error)
		// if isGoCommonUtilsError {
		// 	errMessage = recoveredErr.ErrorWithContext()
		// } else {
		// 	errMessage = recoveredErr.Error()
		// }

		msg := fmt.Sprintf("Unexpected panic: %s", err)
		l.Error(msg)
		l.Error(string(debug.Stack()))
		p.Error(msg)
		os.Exit(1)
	}
}
