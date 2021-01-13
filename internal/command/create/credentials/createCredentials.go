package credentials

import (
	"context"
	"fmt"
	"stevedore/internal/command"
	"stevedore/internal/configuration"
	"stevedore/internal/credentials"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/cobra"
)

type createCredentialsCmdFlags struct {
	DockerRegistryCredentialsDir string
	RegistryHost                 string
	Username                     string
	Password                     string
}

var createCredentialsCmdFlagsVar *createCredentialsCmdFlags

//  NewCommand return an stevedore command object for get
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	createCredentialsCmdFlagsVar = &createCredentialsCmdFlags{}

	createCredentialsCmd := &cobra.Command{
		Use:     "credentials",
		Aliases: []string{"auth"},
		Short:   "Create stevedore docker registry credentials",
		Long:    "",
		RunE:    createCredentialsHandler(ctx, config),
	}

	createCredentialsCmd.Flags().StringVarP(&createCredentialsCmdFlagsVar.RegistryHost, "registry-host", "r", "", "Docker registry host to register credentials")
	createCredentialsCmd.Flags().StringVarP(&createCredentialsCmdFlagsVar.Username, "username", "u", "", "Docker registry username")
	createCredentialsCmd.Flags().StringVarP(&createCredentialsCmdFlagsVar.Password, "password", "p", "", "Docker registry password")
	createCredentialsCmd.Flags().StringVarP(&createCredentialsCmdFlagsVar.DockerRegistryCredentialsDir, "credentials-dir", "d", config.DockerCredentialsDir, "Location path to store docker registry credentials")

	command := &command.StevedoreCommand{
		Command: createCredentialsCmd,
	}

	return command
}

func createCredentialsHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		err := credentials.CreateCredential(createCredentialsCmdFlagsVar.DockerRegistryCredentialsDir, createCredentialsCmdFlagsVar.Username, createCredentialsCmdFlagsVar.Password, createCredentialsCmdFlagsVar.RegistryHost)
		if err != nil {
			return errors.New("(command::createCredentialsHandler)", fmt.Sprintf("Error creating credentials for '%s' on '%s'", createCredentialsCmdFlagsVar.RegistryHost, createCredentialsCmdFlagsVar.DockerRegistryCredentialsDir), err)
		}
		return nil
	}
}
