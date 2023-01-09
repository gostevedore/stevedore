package initialize

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	createconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/cli/create/configuration"
)

// NewCommmand creates a StevedoreCommand for init
func NewCommand(ctx context.Context, e Entrypointer) *command.StevedoreCommand {
	init := createconfiguration.NewCommand(ctx, e)
	init.Command.Use = "initialize"
	init.Command.Aliases = []string{"init"}
	init.Command.Short = "Stevedore command to create and initizalize the configuration"
	init.Command.Long = `
Stevedore command to create and initizalize the configuration.
Initialize command is an create configuration subcommand alias.
`
	init.Command.Example = `
Example setting all configuration parameters:
  stevedore initialize --builders-path /builders --concurrency 4 --config /stevedore-config.yaml --credentials-format json --credentials-local-storage-path /credentials --credentials-storage-type local --enable-semver-tags --force --images-path /images --log-path-file /logs --push-images --semver-tags-template "{{ .Major }}" --semver-tags-template "{{ .Major }}_{{ .Minor }}"
	`

	return init
}
