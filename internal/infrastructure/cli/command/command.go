package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CobraRunFunc is a cobra handler function
type CobraRunFunc func(cmd *cobra.Command, args []string)

// CobraRunEFunc is a cobra handler function which returns an error
type CobraRunEFunc func(cmd *cobra.Command, args []string) error

type CommandOptionsFunc func(c *StevedoreCommand)

// StevedoreCommand defines a stevedore command element
type StevedoreCommand struct {
	Command *cobra.Command
	// TODO: remove configuration from stevedore command, must be injected to services
	// Configuration *configuration.Configuration
}

// AddCommand method add a new subcommand to stevedore command
func (c *StevedoreCommand) AddCommand(cmd *StevedoreCommand) {
	c.Command.AddCommand(cmd.Command)
}

// Execute executes cobra command
func (c *StevedoreCommand) Execute() error {
	if err := c.Command.Execute(); err != nil {
		return err
	}

	return nil
}

// Options configure the stevedore command
func (c *StevedoreCommand) Options(opts ...CommandOptionsFunc) {
	for _, opt := range opts {
		opt(c)
	}
}

func FullCommand(cmd *cobra.Command) string {
	if cmd.Parent() == nil {
		return cmd.Use
	} else {
		parentCommand := FullCommand(cmd.Parent())
		return fmt.Sprintf("%s %s", parentCommand, cmd.Use)
	}
}

// WithRun is a command options function to setup the cobra Run option
func WithRun(f CobraRunFunc) CommandOptionsFunc {
	return func(c *StevedoreCommand) {
		c.Command.Run = f
	}
}

// WithRunE is a command options function to setup the cobra RunE option
func WithRunE(f CobraRunEFunc) CommandOptionsFunc {
	return func(c *StevedoreCommand) {
		c.Command.RunE = f
	}
}
