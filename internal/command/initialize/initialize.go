package initialize

import (
	"context"
	"stevedore/internal/command"
	cmdconf "stevedore/internal/command/create/configuration"
	"stevedore/internal/configuration"
)

// NewCommmand creates a StevedoreCommand for init
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {
	init := cmdconf.NewCommand(ctx, config)
	init.Command.Use = "init"
	return init
}
