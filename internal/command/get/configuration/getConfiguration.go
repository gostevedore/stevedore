package getconfiguration

import (
	"context"

	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/cobra"
)

const (
	columnSeparator = " | "
)

//  NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	getConfigurationCmd := &cobra.Command{
		Use: "configuration",
		Aliases: []string{
			"config",
			"conf",
			"cfg",
		},
		Short: "get configuration",
		Long:  "get configuration",
		RunE:  getConfigurationHandler(ctx, config),
	}

	command := &command.StevedoreCommand{
		Command: getConfigurationCmd,
	}

	return command
}

func getConfigurationHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		table := [][]string{}
		table = append(table, configuration.ConfigurationHeaders())
		configArray, err := config.ToArray()
		if err != nil {
			return errors.New("(command::getConfigurationHandler)", "Error converting configuration to an array", err)
		}
		for _, parameter := range configArray {
			table = append(table, parameter)
		}

		console.PrintTable(table)

		return nil
	}
}
