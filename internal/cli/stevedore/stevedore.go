package stevedore

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/cli/build"
	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/gostevedore/stevedore/internal/cli/command/middleware"
	"github.com/gostevedore/stevedore/internal/cli/completion"
	"github.com/gostevedore/stevedore/internal/cli/promote"
	"github.com/gostevedore/stevedore/internal/configuration"
	buildentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	promoteentrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	"github.com/gostevedore/stevedore/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// stevedoreCmdFlags
type stevedoreCmdFlags struct {
	ConfigFile string
}

var stevedoreCmdFlagsVars *stevedoreCmdFlags

// var conf *configuration.Configuration

//  NewCommand return an stevedore
func NewCommand(ctx context.Context, fs afero.Fs, compatibilityStore CompatibilityStorer, compatibilityReport CompatibilityReporter, console Consoler, config *configuration.Configuration) *command.StevedoreCommand {
	var err error
	var log Logger

	errContext := "(stevedore::NewCommand)"

	if config == nil {
		config, err = configuration.New(fs, compatibilityStore)
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
	}

	stevedoreCmdFlagsVars = &stevedoreCmdFlags{}

	stevedoreCmd := &cobra.Command{
		Use:   "stevedore",
		Short: "Stevedore, the docker images factory",
		Long:  `Stevedore is a useful tool when you need to manage a bunch of Docker images in a standardized way, such on a microservices architecture. It lets you to define how to build your Docker images and their parent-child relationship. It builds automatically the children images when parent ones are done. And many other features which improve the Docker image's building process. Is not a Dockerfile's alternative, but how to use them to build your images`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var logWriter io.Writer

			if len(stevedoreCmdFlagsVars.ConfigFile) > 0 {
				err = config.ReloadConfigurationFromFile(fs, stevedoreCmdFlagsVars.ConfigFile, compatibilityStore)
				if err != nil {
					console.Error(err.Error())
					return errors.New(errContext, fmt.Sprintf("Error loading configuration from file '%s'", stevedoreCmdFlagsVars.ConfigFile), err)
				}
			}

			logWriter, err = generateLogWriter(fs, config.LogPathFile)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
			log = logger.NewLogger(logWriter, logger.LogConsoleEncoderName)

			return nil
		},
		Run: stevedoreHandler,
	}

	stevedoreCmd.PersistentFlags().StringVarP(&stevedoreCmdFlagsVars.ConfigFile, "config", "c", "", "Configuration file location path")

	command := &command.StevedoreCommand{
		Command: stevedoreCmd,
	}

	// entrypoint is not created
	buildEntrypoint := buildentrypoint.NewEntrypoint(
		buildentrypoint.WithWriter(console),
		buildentrypoint.WithFileSystem(fs),
	)
	command.AddCommand(middleware.Command(ctx, build.NewCommand(ctx, compatibilityStore, config, buildEntrypoint), compatibilityReport, log, console))

	promoteEntrypoint := promoteentrypoint.NewEntrypoint(
		promoteentrypoint.WithWriter(console),
		promoteentrypoint.WithFileSystem(fs),
	)
	command.AddCommand(middleware.Command(ctx, promote.NewCommand(ctx, compatibilityStore, config, promoteEntrypoint), compatibilityReport, log, console))

	command.AddCommand(middleware.Command(ctx, completion.NewCommand(ctx, config, command, console), compatibilityReport, log, console))

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

func generateLogWriter(fs afero.Fs, path string) (io.Writer, error) {

	errContext := "(cli::stevedore)"

	file, err := fs.Create(path)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return file, nil
}
