package command

import (
	"stevedore/internal/configuration"

	"github.com/spf13/cobra"
)

// CobraRunFunc is a cobra handler function
type CobraRunFunc func(cmd *cobra.Command, args []string)

// CobraRunEFunc is a cobra handler function which returns an error
type CobraRunEFunc func(cmd *cobra.Command, args []string) error

// StevedoreCommand defines a stevedore command element
type StevedoreCommand struct {
	Command       *cobra.Command
	Configuration *configuration.Configuration
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
