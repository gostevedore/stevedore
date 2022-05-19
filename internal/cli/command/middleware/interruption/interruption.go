package interruption

import (
	"context"
	"os"
	"os/signal"

	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

// NewCommand is a middleware to manage interruptions
func NewCommand(ctx context.Context, c *command.StevedoreCommand, p Consoler, l Logger) *command.StevedoreCommand {

	if c.Command.PersistentPreRun != nil {
		c.Command.PersistentPreRun = interruptManagement(ctx, l, p, c.Command.PersistentPreRun)
	}

	if c.Command.PersistentPreRunE != nil {
		c.Command.PersistentPreRunE = interruptManagementWithError(ctx, l, p, c.Command.PersistentPreRunE)
	}

	if c.Command.PreRun != nil {
		c.Command.PreRun = interruptManagement(ctx, l, p, c.Command.PreRun)
	}

	if c.Command.PreRunE != nil {
		c.Command.PreRunE = interruptManagementWithError(ctx, l, p, c.Command.PreRunE)
	}

	if c.Command.Run != nil {
		c.Command.Run = interruptManagement(ctx, l, p, c.Command.Run)
	}

	if c.Command.RunE != nil {
		c.Command.RunE = interruptManagementWithError(ctx, l, p, c.Command.RunE)
	}

	if c.Command.PostRun != nil {
		c.Command.PostRun = interruptManagement(ctx, l, p, c.Command.PostRun)
	}

	if c.Command.PostRunE != nil {
		c.Command.PostRunE = interruptManagementWithError(ctx, l, p, c.Command.PostRunE)
	}

	if c.Command.PersistentPostRun != nil {
		c.Command.PersistentPostRun = interruptManagement(ctx, l, p, c.Command.PersistentPostRun)
	}

	if c.Command.PersistentPostRunE != nil {
		c.Command.PersistentPostRunE = interruptManagementWithError(ctx, l, p, c.Command.PersistentPostRunE)
	}

	return c
}

func interruptManagement(ctx context.Context, l Logger, p Consoler, f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var cancel context.CancelFunc
		done := make(chan struct{})

		ctx, cancel = context.WithCancel(ctx)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		defer func() {
			signal.Stop(c)
			cancel()
		}()

		go func() {
			select {
			case <-c:
				msg := "'" + cmd.Use + "' execution has been interrupted"
				l.Warn(msg)
				p.Warn(msg)

				cancel()
			case <-ctx.Done():
			}
		}()

		go func() {
			f(cmd, args)
			done <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			msg := "'" + cmd.Use + "' execution has been cancelled"
			l.Warn(msg)
			p.Warn(msg)
		case <-done:
		}
	}
}

func interruptManagementWithError(ctx context.Context, l Logger, p Consoler, f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var cancel context.CancelFunc
		done := make(chan struct{})
		errChan := make(chan error)

		ctx, cancel = context.WithCancel(ctx)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		defer func() {
			signal.Stop(c)
			cancel()
		}()

		go func() {
			select {
			case <-c:
				msg := "'" + cmd.Use + "' execution has been interrupted"
				l.Warn(msg)
				p.Warn(msg)

				cancel()
			case <-ctx.Done():
			}
		}()

		go func() {
			err := f(cmd, args)
			if err != nil {
				errChan <- err
				return
			}
			done <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			msg := "'" + cmd.Use + "' execution has been cancelled"
			l.Warn(msg)
			p.Warn(msg)
		case <-errChan:
			return <-errChan
		case <-done:
		}

		return nil
	}
}
