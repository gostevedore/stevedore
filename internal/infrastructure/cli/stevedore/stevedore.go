package stevedore

import (
	"context"
	"fmt"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	buildentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	getcredentialsentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/credentials"
	promoteentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command/middleware"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/completion"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/get"
	getcredentials "github.com/gostevedore/stevedore/internal/infrastructure/cli/get/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/version"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/loader"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stevedoreCmdFlags
type stevedoreCmdFlags struct {
	ConfigFile string
}

var stevedoreCmdFlagsVars *stevedoreCmdFlags

// var conf *configuration.Configuration

//  NewCommand return an stevedore
func NewCommand(ctx context.Context, fs afero.Fs, compatibilityStore CompatibilityStorer, compatibilityReport CompatibilityReporter, console Consoler, log Logger, config *configuration.Configuration) *command.StevedoreCommand {
	var err error
	//	var log *logger.Logger

	errContext := "(stevedore::NewCommand)"

	if config == nil {
		configLoader := loader.NewConfigurationLoader(viper.New())
		config, err = configuration.New(fs, configLoader, compatibilityStore)
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
	}

	stevedoreCmdFlagsVars = &stevedoreCmdFlags{}

	stevedoreCmd := &cobra.Command{
		Use:   "stevedore",
		Short: "Stevedore, the docker images factory",
		Long: `Stevedore is a tool to manage bunches of Docker images builds in just one command. It improves the way you build and promote your Docker images when you have a lot of them. Is not a Dockerfile's alternative, but how to use them to build your images.
You just need to define how each image should be built and the relationship among the images. At this moment, everything is ready to build Docker images: build a single image, build all versions of the same images, build an image and all its children.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(stevedoreCmdFlagsVars.ConfigFile) > 0 {
				err = config.ReloadConfigurationFromFile(stevedoreCmdFlagsVars.ConfigFile)
				if err != nil {
					console.Error(err.Error())
					return errors.New(errContext, fmt.Sprintf("Error loading configuration from file '%s'", stevedoreCmdFlagsVars.ConfigFile), err)
				}
			}

			log.ReloadWithWriter(config.LogWriter)

			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Sync()
		},
		Run: stevedoreHandler,
	}

	stevedoreCmd.PersistentFlags().StringVarP(&stevedoreCmdFlagsVars.ConfigFile, "config", "c", "", "Configuration file location path")

	command := &command.StevedoreCommand{
		Command: stevedoreCmd,
	}

	// Completion
	command.AddCommand(
		middleware.Command(ctx, completion.NewCommand(ctx, config, command, console), compatibilityReport, log, console),
	)

	// Version
	command.AddCommand(
		middleware.Command(ctx, version.NewCommand(ctx, console), compatibilityReport, log, console),
	)

	// Build
	buildEntrypoint := buildentrypoint.NewEntrypoint(
		buildentrypoint.WithWriter(console),
		buildentrypoint.WithFileSystem(fs),
		buildentrypoint.WithCompatibility(compatibilityStore),
	)
	command.AddCommand(
		middleware.Command(ctx, build.NewCommand(ctx, compatibilityStore, config, buildEntrypoint), compatibilityReport, log, console),
	)

	// Promote
	promoteEntrypoint := promoteentrypoint.NewEntrypoint(
		promoteentrypoint.WithWriter(console),
		promoteentrypoint.WithFileSystem(fs),
	)
	command.AddCommand(
		middleware.Command(ctx, promote.NewCommand(ctx, compatibilityStore, config, promoteEntrypoint), compatibilityReport, log, console),
	)

	// Get

	// Get credentials
	getCredentialsEntrypoint := getcredentialsentrypoint.NewEntrypoint(
		getcredentialsentrypoint.WithWriter(console),
		getcredentialsentrypoint.WithFileSystem(fs),
		getcredentialsentrypoint.WithCompatibilitier(compatibilityStore),
	)
	getCredentialsCommand := middleware.Command(ctx, getcredentials.NewCommand(ctx, config, getCredentialsEntrypoint), compatibilityReport, log, console)

	getCommand := get.NewCommand(ctx, getCredentialsCommand)
	command.AddCommand(getCommand)

	// command.AddCommand(middleware.Middleware(create.NewCommand(ctx, config)))
	// command.AddCommand(middleware.Middleware(get.NewCommand(ctx, config)))
	// command.AddCommand(middleware.Middleware(initialize.NewCommand(ctx, config)))
	// command.AddCommand(middleware.Middleware(moo.NewCommand(ctx, config)))
	// command.AddCommand(middleware.Middleware(promote.NewCommand(ctx, config)))
	// command.AddCommand(middleware.Middleware(version.NewCommand(ctx, config)))

	return command
}

func stevedoreHandler(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, args)
}
