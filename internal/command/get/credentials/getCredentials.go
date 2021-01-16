package getcredentials

import (
	"context"

	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/cobra"
)

const (
	columnSeparator = " | "
)

type getCredentialsCmdFlags struct {
	Wide bool
}

var getCredentialsCmdFlagsVar *getCredentialsCmdFlags

//  NewCommand return an stevedore command object for get builders
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	getCredentialsCmdFlagsVar = &getCredentialsCmdFlags{}

	getCredentialsCmd := &cobra.Command{
		Use: "credentials",
		Aliases: []string{
			"auth",
			"auths",
			"credential",
		},
		Short: "get credentials return all credentials defined",
		Long:  "get credentials return all credentials defined",
		RunE:  getCredentialssHandler(ctx, config),
	}

	command := &command.StevedoreCommand{
		Command: getCredentialsCmd,
	}

	getCredentialsCmd.Flags().BoolVarP(&getCredentialsCmdFlagsVar.Wide, "wide", "w", false, "Show wide docker registry credentials information")

	return command
}

func getCredentialssHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {

	return func(cmd *cobra.Command, args []string) error {

		var err error
		var creds [][]string
		var table [][]string

		creds, err = credentials.ListRegistryCredentials(getCredentialsCmdFlagsVar.Wide)
		if err != nil {
			return errors.New("(command::getCredentialssHandler)", "Error listing registry credentials", err)
		}

		table = make([][]string, len(creds)+1)
		table[0] = credentials.ListRegistryCredentialsHeader(getCredentialsCmdFlagsVar.Wide)
		copy(table[1:], creds)

		console.PrintTable(table)

		return nil
	}
}
