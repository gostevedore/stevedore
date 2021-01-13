package stevedore

import (
	"context"
	"fmt"
	"os"
	"stevedore/internal/command"
	"stevedore/internal/command/build"
	"stevedore/internal/command/completion"
	"stevedore/internal/command/create"
	"stevedore/internal/command/get"
	"stevedore/internal/command/initialize"
	"stevedore/internal/command/middleware"
	"stevedore/internal/command/moo"
	"stevedore/internal/command/promote"
	"stevedore/internal/command/version"
	"stevedore/internal/configuration"
	"stevedore/internal/credentials"
	"stevedore/internal/logger"
	"stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/cobra"
)

// stevedoreCmdFlags
type stevedoreCmdFlags struct {
	ConfigFile string
}

var stevedoreCmdFlagsVars *stevedoreCmdFlags
var cancelContext context.Context

//  NewCommand return an stevedore command object
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {
	var err error

	console.Init(os.Stdout)

	if config == nil {
		config, err = configuration.New()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
	}

	stevedoreCmdFlagsVars = &stevedoreCmdFlags{}

	stevedoreCmd := &cobra.Command{
		Use:   "stevedore",
		Short: "Stevedore orchestrates image building",
		Long:  `Stevedore orchestrates image building`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			err = logger.Init(config.LogPathFile, logger.LogConsoleEncoderName)
			if err != nil {
				return errors.New("(stevedore::NewCommand)", "Error initializing logger", err)
			}

			if len(stevedoreCmdFlagsVars.ConfigFile) > 0 {
				err = config.ReloadConfigurationFromFile(stevedoreCmdFlagsVars.ConfigFile)
				if err != nil {
					console.Print(err.Error())
					return errors.New("(stevedore::NewCommand)", fmt.Sprintf("Error loading configuration from file '%s'", stevedoreCmdFlagsVars.ConfigFile), err)
				}
				logger.Info(fmt.Sprintf("Configuration reloaded from '%s'", stevedoreCmdFlagsVars.ConfigFile))
			}

			err = credentials.LoadCredentials(config.DockerCredentialsDir)
			if err != nil {
				err := errors.New("(stevedore::NewCommand)", fmt.Sprintf("Credentials loading credentials from directory  '%s'", config.DockerCredentialsDir), err)
				console.Print(err.Error())
				logger.Info(err.ErrorWithContext())
			}

			return nil
		},
		Run: stevedoreHandler,
	}

	stevedoreCmd.PersistentFlags().StringVarP(&stevedoreCmdFlagsVars.ConfigFile, "config", "c", "", "Configuration file location path")

	command := &command.StevedoreCommand{
		Configuration: config,
		Command:       stevedoreCmd,
	}

	command.AddCommand(middleware.Middleware(build.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(create.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(completion.NewCommand(ctx, config, command)))
	command.AddCommand(middleware.Middleware(get.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(initialize.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(moo.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(promote.NewCommand(ctx, config)))
	command.AddCommand(middleware.Middleware(version.NewCommand(ctx, config)))

	return command
}

func stevedoreHandler(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, args)
}
