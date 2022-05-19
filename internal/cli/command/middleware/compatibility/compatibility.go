package compatibility

import (
	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/spf13/cobra"
)

// Compatibility is a middleware to notice about compatibility missages
func NewCommand(c *command.StevedoreCommand, r CompatibilityReporter) *command.StevedoreCommand {
	var postRunFunc func(cmd *cobra.Command, args []string)

	if c.Command.PostRun != nil {
		postRunFunc = func(cmd *cobra.Command, args []string) {
			defer r.Report()
			c.Command.PostRun(cmd, args)
		}
	} else {
		postRunFunc = func(cmd *cobra.Command, args []string) {
			r.Report()
		}
	}

	c.Command.PostRun = postRunFunc

	return c
}
