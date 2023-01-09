package stevedore

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	buildentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	createconfigurationentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/configuration"
	createcredentialsentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/credentials"
	getbuildersentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/builders"
	getconfigurationentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/configuration"
	getcredentialsentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/credentials"
	getimagesentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/images"
	promoteentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command/middleware"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/completion"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/create"
	createconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/cli/create/configuration"
	createcredentials "github.com/gostevedore/stevedore/internal/infrastructure/cli/create/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/get"
	getbuilders "github.com/gostevedore/stevedore/internal/infrastructure/cli/get/builders"
	getconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/cli/get/configuration"
	getcredentials "github.com/gostevedore/stevedore/internal/infrastructure/cli/get/credentials"
	getimages "github.com/gostevedore/stevedore/internal/infrastructure/cli/get/images"
	initizalize "github.com/gostevedore/stevedore/internal/infrastructure/cli/initialize"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/version"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// stevedoreCmdFlags
type stevedoreCmdFlags struct {
	ConfigFile string
	Debug      bool
}

var stevedoreCmdFlagsVars *stevedoreCmdFlags

// NewCommand return an stevedore
func NewCommand(ctx context.Context, fs afero.Fs, compatibilityStore CompatibilityStorer, compatibilityReport CompatibilityReporter, console Consoler, log Logger, config *configuration.Configuration) *command.StevedoreCommand {

	errContext := "(cli::stevedore::NewCommand)"
	_ = errContext

	stevedoreCmdFlagsVars = &stevedoreCmdFlags{}

	stevedoreCmd := &cobra.Command{
		Use:   "stevedore [COMMAND] [OPTIONS]",
		Short: "Stevedore, the docker images factory",
		Long: `
Stevedore is a Docker images factory, a tool that helps you to manage bunches of Docker image builds in just one command. It is not an alternative to Dockerfile or Buildkit, but a way to improve your building and promote experience
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if stevedoreCmdFlagsVars.ConfigFile != "" {
				err = config.ReloadConfigurationFromFile(stevedoreCmdFlagsVars.ConfigFile)
				if err != nil {
					return errors.New(errContext, fmt.Sprintf("Error loading configuration from file '%s'", stevedoreCmdFlagsVars.ConfigFile), err)
				}
			}
			log.ReloadWithWriter(config.LogWriter)

			return nil
		},
	}

	stevedoreCmd.PersistentFlags().StringVarP(&stevedoreCmdFlagsVars.ConfigFile, "config", "c", "", "Configuration file location path")
	stevedoreCmd.PersistentFlags().BoolVar(&stevedoreCmdFlagsVars.Debug, "debug", false, "Enable debug mode")

	// Stevedore command
	command := &command.StevedoreCommand{
		Command: stevedoreCmd,
	}

	command = middleware.Command(ctx, command, compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	//
	// Completion command
	//
	command.AddCommand(
		middleware.Command(ctx, completion.NewCommand(ctx, config, command, console), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug),
	)

	//
	// Version ccommand
	//
	command.AddCommand(
		middleware.Command(ctx, version.NewCommand(ctx, console), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug),
	)

	//
	// Build ccommand
	//
	buildEntrypoint := buildentrypoint.NewEntrypoint(
		buildentrypoint.WithWriter(console),
		buildentrypoint.WithFileSystem(fs),
		buildentrypoint.WithCompatibility(compatibilityStore),
	)
	command.AddCommand(
		middleware.Command(ctx, build.NewCommand(ctx, compatibilityStore, config, buildEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug),
	)

	//
	// Create command
	//

	// Create configuration
	createConfigurationEntrypoint := createconfigurationentrypoint.NewCreateConfigurationEntrypoint(
		createconfigurationentrypoint.WithFileSystem(fs),
	)
	createConfigurationCommand := middleware.Command(ctx, createconfiguration.NewCommand(ctx, createConfigurationEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	// Create credentials
	createCredentialsEntrypoint := createcredentialsentrypoint.NewCreateCredentialsEntrypoint(
		createcredentialsentrypoint.WithConsole(console),
		createcredentialsentrypoint.WithFileSystem(fs),
		createcredentialsentrypoint.WithCompatibility(compatibilityStore),
	)
	createCredentialsCommand := middleware.Command(ctx, createcredentials.NewCommand(ctx, compatibilityStore, config, createCredentialsEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	// Create root command
	createCommand := create.NewCommand(
		ctx,
		createConfigurationCommand,
		createCredentialsCommand,
	)
	command.AddCommand(createCommand)

	//
	// Get command
	//

	// Get configuration subcommand
	getConfigurationEntrypoint := getconfigurationentrypoint.NewGetConfigurationEntrypoint(
		getconfigurationentrypoint.WithWriter(console),
	)
	getConfigurationCommand := middleware.Command(ctx, getconfiguration.NewCommand(ctx, config, getConfigurationEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	// Get credentials subcommand
	getCredentialsEntrypoint := getcredentialsentrypoint.NewEntrypoint(
		getcredentialsentrypoint.WithWriter(console),
		getcredentialsentrypoint.WithFileSystem(fs),
		getcredentialsentrypoint.WithCompatibility(compatibilityStore),
	)
	getCredentialsCommand := middleware.Command(ctx, getcredentials.NewCommand(ctx, config, getCredentialsEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	// Get images subcommand
	getImagesEntrypoint := getimagesentrypoint.NewGetImagesEntrypoint(
		getimagesentrypoint.WithWriter(console),
		getimagesentrypoint.WithFileSystem(fs),
		getimagesentrypoint.WithCompatibility(compatibilityStore),
	)
	getImagesCommand := middleware.Command(ctx, getimages.NewCommand(ctx, config, getImagesEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	// Get builders subcommand
	getBuildersEntrypoint := getbuildersentrypoint.NewGetBuildersEntrypoint(
		getbuildersentrypoint.WithWriter(console),
		getbuildersentrypoint.WithFileSystem(fs),
		getbuildersentrypoint.WithCompatibility(compatibilityStore),
	)
	getBuildersCommand := middleware.Command(ctx, getbuilders.NewCommand(ctx, config, getBuildersEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)

	// Get root command
	getCommand := get.NewCommand(
		ctx,
		getBuildersCommand,
		getConfigurationCommand,
		getCredentialsCommand,
		getImagesCommand,
	)
	command.AddCommand(getCommand)

	//
	// Initialize command
	//

	// it uses the entrypoint that create configuration
	initializeCommand := middleware.Command(ctx, initizalize.NewCommand(ctx, createConfigurationEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug)
	command.AddCommand(initializeCommand)

	//
	// Promote command
	//
	promoteEntrypoint := promoteentrypoint.NewEntrypoint(
		promoteentrypoint.WithWriter(console),
		promoteentrypoint.WithFileSystem(fs),
		promoteentrypoint.WithCompatibility(compatibilityStore),
	)
	command.AddCommand(
		middleware.Command(ctx, promote.NewCommand(ctx, compatibilityStore, config, promoteEntrypoint), compatibilityReport, log, console, &stevedoreCmdFlagsVars.Debug),
	)

	return command
}
