package completion

import (
	"context"
	"stevedore/internal/command"
	"stevedore/internal/configuration"
	"stevedore/internal/ui/console"

	"github.com/spf13/cobra"
)

//  NewCommand return an stevedore command object for dev
func NewCommand(ctx context.Context, config *configuration.Configuration, rootCmd *command.StevedoreCommand) *command.StevedoreCommand {

	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run
	
	. <(stevedore completion)
	
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
