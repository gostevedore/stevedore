package completion

import (
	"context"

	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"

	"github.com/spf13/cobra"
)

//  NewCommand return an stevedore command object for dev
func NewCommand(ctx context.Context, config *configuration.Configuration, rootCmd *command.StevedoreCommand) *command.StevedoreCommand {

	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Stevedore command to generate shell completions",
		Long: `To load stevedore completion run
	$ . <(stevedore completion)
	
	To configure your bash shell to load completions for each session add to your bashrc
	# ~/.bashrc or ~/.profile
	. <(stevedore completion)
	`,
		Run: completionHandler(ctx, rootCmd),
	}

	command := &command.StevedoreCommand{
		Command: completionCmd,
	}

	return command
}

func completionHandler(ctx context.Context, command *command.StevedoreCommand) command.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) {
		command.Command.GenBashCompletion(console.GetConsole())
	}
}
