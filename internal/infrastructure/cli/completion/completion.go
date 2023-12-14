package completion

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"

	"github.com/spf13/cobra"
)

// NewCommand return an stevedore command object for completion
func NewCommand(ctx context.Context, config *configuration.Configuration, rootCmd *command.StevedoreCommand, cons Consoler) *command.StevedoreCommand {

	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Stevedore command to generate shell completions",
		Long: `To load stevedore completion run
	$ . <(stevedore completion)
	
	To configure your bash shell to load completions for each session add to your bashrc
	# ~/.bashrc or ~/.profile
	. <(stevedore completion)
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.Command.GenBashCompletion(cons)
		},
	}

	command := &command.StevedoreCommand{
		Command: completionCmd,
	}

	return command
}
