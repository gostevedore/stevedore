package configuration

import (
	"context"
	"fmt"
	"os"

	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/logger"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/cobra"
)

type createConfigurationCmdFlags struct {
	Force                        bool
	ConfigFile                   string
	TreePathFile                 string
	BuilderPathFile              string
	LogPathFile                  string
	NumWorkers                   int
	PushImages                   bool
	BuildOnCascade               bool
	DockerCredentialsDir         string
	EnableSemanticVersionTags    bool
	SemanticVersionTagsTemplates []string
	CredentialsRegistryHost      string
	CredentialsUsername          string
	CredentialsPassword          string
}

var createConfigurationCmdFlagsVar *createConfigurationCmdFlags

//  NewCommand return an stevedore command object for get
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	var err error

	createConfigurationCmdFlagsVar = &createConfigurationCmdFlags{}

	createConfigurationCmd := &cobra.Command{
		Use:     "configuration",
		Aliases: []string{"config"},
		Short:   "Create stevedore configuration file",
		Long:    "",
		RunE:    createConfigurationHandler(ctx, config),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err = logger.Init(config.LogPathFile, logger.LogConsoleEncoderName)
			if err != nil {
				return errors.New("(stevedore::NewCommand)", "Error initializing logger", err)
			}

			return nil
		},
	}

	createConfigurationCmd.PersistentFlags().StringVarP(&createConfigurationCmdFlagsVar.ConfigFile, "config", "c", "", "Configuration file location path")
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.TreePathFile, "tree-path-file", "t", "", fmt.Sprintf("Images tree location path. Its default value is '%s'", configuration.DefaultTreePathFile))
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.BuilderPathFile, "builder-path-file", "b", "", fmt.Sprintf("Builders location path. Its default value is '%s'", configuration.DefaultBuilderPathFile))
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.LogPathFile, "log-path-file", "l", "", fmt.Sprintf("Log file location path. Its default value is '%s'", configuration.DefaultLogPathFile))
	createConfigurationCmd.Flags().IntVarP(&createConfigurationCmdFlagsVar.NumWorkers, "num-workers", "w", -1, fmt.Sprintf("It defines the number of workers to build images which corresponds to the number of images that can be build concurrently. Its default value is '%d'", configuration.DefaultNumWorker))
	createConfigurationCmd.Flags().BoolVarP(&createConfigurationCmdFlagsVar.PushImages, "no-push-images", "P", false, fmt.Sprintf("On build, push images automatically after it finishes. Its default value is '%t'", configuration.DefaultPushImages))
	createConfigurationCmd.Flags().BoolVarP(&createConfigurationCmdFlagsVar.BuildOnCascade, "build-on-cascade", "C", false, fmt.Sprintf("On build, start children images building once an image build is finished. Its default value is '%t'", configuration.DefaultBuildOnCascade))
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.DockerCredentialsDir, "credentials-dir", "d", "", fmt.Sprintf("Location path to store docker registry credentials. Its default value is '%s'", configuration.DefaultDockerCredentialsDir))
	createConfigurationCmd.Flags().BoolVarP(&createConfigurationCmdFlagsVar.EnableSemanticVersionTags, "enable-semver-tags", "s", false, fmt.Sprintf("Generate extra tags when the main image tags is semver 2.0.0 compliance. Its default value is '%t'", configuration.DefaultEnableSemanticVersionTags))
	createConfigurationCmd.Flags().StringSliceVarP(&createConfigurationCmdFlagsVar.SemanticVersionTagsTemplates, "semver-tags-template", "T", []string{}, fmt.Sprintf("List of templates which define those extra tags to generate when 'semantic_version_tags_enabled' is enabled. Its default value is '%v'", configuration.DefaultSemanticVersionTagsTemplates))
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.CredentialsRegistryHost, "credentials-registry-host", "r", "", "Docker registry host to register credentials")
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.CredentialsUsername, "credentials-username", "u", "", "Docker registry username. It is ignored unless `credentials-regristry` value is defined")
	createConfigurationCmd.Flags().StringVarP(&createConfigurationCmdFlagsVar.CredentialsPassword, "credentials-password", "p", "", "Docker registry password. It is ignored unless `credentials-regristry` value is defined")
	createConfigurationCmd.Flags().BoolVar(&createConfigurationCmdFlagsVar.Force, "force", false, "Force to create configuration file when the file already exists")

	command := &command.StevedoreCommand{
		Command: createConfigurationCmd,
	}

	return command
}

func createConfigurationHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {

	return func(cmd *cobra.Command, args []string) error {
		var err error
		var configFile *os.File
		var fileInfo os.FileInfo

		cfg := configuration.Configuration{}
		file := configuration.ConfigFileUsed()

		if createConfigurationCmdFlagsVar.ConfigFile != "" {
			file = createConfigurationCmdFlagsVar.ConfigFile
		}
		if file == "" {
			file = configuration.DefaultConfigFile
		}

		if createConfigurationCmdFlagsVar.TreePathFile != "" {
			cfg.TreePathFile = createConfigurationCmdFlagsVar.TreePathFile
		}

		if createConfigurationCmdFlagsVar.BuilderPathFile != "" {
			cfg.BuilderPathFile = createConfigurationCmdFlagsVar.BuilderPathFile
		}

		if createConfigurationCmdFlagsVar.LogPathFile != "" {
			cfg.LogPathFile = createConfigurationCmdFlagsVar.LogPathFile
		}

		if createConfigurationCmdFlagsVar.NumWorkers != -1 {
			cfg.NumWorkers = createConfigurationCmdFlagsVar.NumWorkers
		}

		if createConfigurationCmdFlagsVar.PushImages {
			cfg.PushImages = false
		} else {
			cfg.PushImages = true
		}

		if createConfigurationCmdFlagsVar.BuildOnCascade {
			cfg.BuildOnCascade = createConfigurationCmdFlagsVar.BuildOnCascade
		}

		if createConfigurationCmdFlagsVar.DockerCredentialsDir != "" {
			cfg.DockerCredentialsDir = createConfigurationCmdFlagsVar.DockerCredentialsDir
		}

		if createConfigurationCmdFlagsVar.EnableSemanticVersionTags {
			cfg.EnableSemanticVersionTags = createConfigurationCmdFlagsVar.EnableSemanticVersionTags
		}

		if len(createConfigurationCmdFlagsVar.SemanticVersionTagsTemplates) > 0 {
			for _, tmpl := range createConfigurationCmdFlagsVar.SemanticVersionTagsTemplates {
				cfg.SemanticVersionTagsTemplates = append(cfg.SemanticVersionTagsTemplates, tmpl)
			}
		}

		fileInfo, _ = os.Stat(file)
		if fileInfo != nil && !createConfigurationCmdFlagsVar.Force {
			return errors.New("(command::createConfigurationHandler)", fmt.Sprintf("Configuration file '%s' already exist and will not be created", file))
		}

		configFile, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return errors.New("(command::createConfigurationHandler)", fmt.Sprintf("File '%s' could not be opened", file), err)
		}
		defer configFile.Close()

		cfg.WriteConfigurationFile(configFile)
		console.Print(fmt.Sprintf("Stevedore configuration file '%s' has been created", file))

		if createConfigurationCmdFlagsVar.CredentialsRegistryHost != "" {

			dir := config.DockerCredentialsDir
			if cfg.DockerCredentialsDir != "" {
				dir = cfg.DockerCredentialsDir
			}

			err = credentials.CreateCredential(dir, createConfigurationCmdFlagsVar.CredentialsUsername, createConfigurationCmdFlagsVar.CredentialsPassword, createConfigurationCmdFlagsVar.CredentialsRegistryHost)
			if err != nil {
				return errors.New("(command::createConfigurationHandler)", fmt.Sprintf("Error creating credentials for '%s' on '%s'", createConfigurationCmdFlagsVar.CredentialsRegistryHost, dir), err)
			}
			console.Print(fmt.Sprintf("Credentials for '%s' stored on '%s'", createConfigurationCmdFlagsVar.CredentialsRegistryHost, dir))
		}

		return nil
	}
}
