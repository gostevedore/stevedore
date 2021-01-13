package middleware

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"stevedore/internal/command"
	"stevedore/internal/logger"
	"stevedore/internal/ui/console"

	"github.com/spf13/cobra"
)

func Middleware(c *command.StevedoreCommand) *command.StevedoreCommand {

	// Middleware is almost ready to use context but it is not used
	ctx := context.TODO()

	// apply to command run
	if c.Command.Run != nil {
		c.Command.Run = InterruptManagement(ctx, PanicRecover(ctx, AuditCommand(ctx, c.Command.Run)))
	}

	if c.Command.RunE != nil {
		c.Command.Run = InterruptManagement(ctx, PanicRecover(ctx, AuditCommand(ctx, ErrorManagement(ctx, c.Command.RunE))))
		c.Command.RunE = nil
	}

	return c
}

// AuditCommand is a middleware which audits the command exection
func AuditCommand(ctx context.Context, f command.CobraRunFunc) command.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) {
		logger.Info("Executing command 'stevedore ", cmd.Use, args, "'")
		f(cmd, args)
	}
}

// ErrorManagement is a middleware which controls errors returned during command execution
func ErrorManagement(ctx context.Context, f command.CobraRunEFunc) command.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) {
		err := f(cmd, args)
		if err != nil {
			logger.Error(err.Error())
			console.Print(err.Error())
			os.Exit(1)
		}
	}
}

func PanicRecover(ctx context.Context, f command.CobraRunFunc) command.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s: %s", "Unexpected panic:", err)
				logger.Error(msg)
				logger.Error(string(debug.Stack()))
				console.Print(msg)
				os.Exit(1)
			}
		}()

		f(cmd, args)
	}
}

func InterruptManagement(ctx context.Context, f command.CobraRunFunc) command.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) {
		var cancel context.CancelFunc
		done := make(chan int8)

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
				logger.Warn(msg)
				console.Warn(msg)

				cancel()
			case <-ctx.Done():
			}
		}()

		go func() {
			f(cmd, args)
			done <- int8(0)
		}()

		select {
		case <-ctx.Done():
			msg := "'" + cmd.Use + "' execution has been cancelled"
			logger.Warn(msg)
			console.Warn(msg)
		case <-done:
		}
	}
}
