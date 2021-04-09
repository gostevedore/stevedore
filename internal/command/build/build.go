package build

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/engine"
	"github.com/gostevedore/stevedore/internal/types"

	"github.com/spf13/cobra"
)

const (
	setVarsSplitToken = "="
)

type buildCmdFlags struct {
	BuildBuilderName            string
	BuildConfiguraction         string
	BuildDryRun                 bool
	BuildOnCascade              bool
	ConnectionLocal             bool
	CascadeDepth                int
	Debug                       bool
	EnableSemanticVersionTags   bool
	ImageFromName               string
	ImageFromRegistryHost       string
	ImageFromRegistryNamespace  string
	ImageFromVersion            string
	ImageName                   string
	ImageRegistryHost           string
	ImageRegistryNamespace      string
	Inventory                   string
	Limit                       string
	NumWorkers                  int
	PersistentVars              []string
	PushImages                  bool
	SemanticVersionTagsTemplate []string
	Tags                        []string
	Vars                        []string
	Versions                    []string
}

var buildCmdFlagsVar *buildCmdFlags

// init define the arguments and add build command to root command cli
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {
	buildCmdFlagsVar = &buildCmdFlags{}

	buildCmd := &cobra.Command{
		Use:   "build <image>",
		Short: "Stevedore command to build images",
		Long: `Stevedore command to build images

  Example: 
  	Command to build the 'focal' version of an image named 'ubuntu-base':
    stevedore build ubuntu-base --image-version focal
		`,
		RunE: buildHandler(ctx, config),
	}

	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.BuildBuilderName, "builder-name", "b", "", "Intermediate builder's container name [only applies to ansible-playbook builders]")
	buildCmd.Flags().BoolVarP(&buildCmdFlagsVar.BuildOnCascade, "cascade", "C", false, "Build images on cascade. Children's image build is started once the image build finishes")
	buildCmd.Flags().BoolVarP(&buildCmdFlagsVar.ConnectionLocal, "connection-local", "L", false, "Use local connection for ansible [only applies to ansible-playbook builders]")
	buildCmd.Flags().IntVarP(&buildCmdFlagsVar.CascadeDepth, "cascade-depth", "d", -1, "Number images levels to build when build on cascade is executed")
	buildCmd.Flags().BoolVar(&buildCmdFlagsVar.Debug, "debug", false, "Enable debug mode to show build options")
	buildCmd.Flags().BoolVarP(&buildCmdFlagsVar.BuildDryRun, "dry-run", "D", false, "Run a dry-run build")
	buildCmd.Flags().BoolVarP(&buildCmdFlagsVar.EnableSemanticVersionTags, "enable-semver-tags", "S", false, "Generate a set of tags for the image based on the semantic version tree when main version is semver 2.0.0 compliance")
	buildCmd.Flags().StringSliceVarP(&buildCmdFlagsVar.Versions, "image-version", "v", []string{}, "Image versions to be built. One or more image versions could be built")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageFromName, "image-from", "I", "", "Image (FROM) parent's name")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageFromRegistryHost, "image-from-registry", "R", "", "Image (FROM) parent's registry host")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageFromRegistryNamespace, "image-from-namespace", "N", "", "Image (FROM) parent's registry namespace")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageFromVersion, "image-from-version", "V", "", "Image (FROM) parent's version")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageName, "image-name", "i", "", "Image name- It overrides image tree image name")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.Inventory, "inventory", "H", "", "Specify inventory hosts' path or comma separated list of hosts [only applies to Ansible builders]")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.Limit, "limit", "l", "", "Further limit selected hosts to an additional pattern [only applies to Ansible builders]")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageRegistryNamespace, "namespace", "n", "", "Image's registry namespace where image will be stored")
	buildCmd.Flags().IntVarP(&buildCmdFlagsVar.NumWorkers, "num-workers", "w", 0, "Number of workers to execute builds")
	buildCmd.Flags().BoolVarP(&buildCmdFlagsVar.PushImages, "no-push", "P", false, "Do not push the image to registry once it is built")
	buildCmd.Flags().StringVarP(&buildCmdFlagsVar.ImageRegistryHost, "registry", "r", "", "Image's registry host where image will be stored")
	buildCmd.Flags().StringSliceVarP(&buildCmdFlagsVar.SemanticVersionTagsTemplate, "semver-tags-template", "T", []string{}, "List templates to generate tags following semantic version expression")
	buildCmd.Flags().StringSliceVarP(&buildCmdFlagsVar.PersistentVars, "set-persistent", "p", []string{}, "Set persistent variables to use during the build. A persistent variable will be available on child image during its build and could not be overwrite. The format of each variable must be <key>=<value>")
	buildCmd.Flags().StringSliceVarP(&buildCmdFlagsVar.Vars, "set", "s", []string{}, "Set variables to use during the build. The format of each variable must be <key>=<value>")
	buildCmd.Flags().StringSliceVarP(&buildCmdFlagsVar.Tags, "tag", "t", []string{}, "Give an extra tag for the docker image")

	command := &command.StevedoreCommand{
		Command: buildCmd,
	}

	return command
}

func buildHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {

	return func(cmd *cobra.Command, args []string) error {

		var err error
		var imagesEngine *engine.ImagesEngine
		var persistentVars map[string]interface{}
		var vars map[string]interface{}
		var buildImageName string

		if cmd.Flags().NArg() == 0 {
			return errors.New("(command::buildHandler)", "Is required an image name")
		} else {
			buildImageName = cmd.Flags().Arg(0)
			if cmd.Flags().NArg() > 1 {
				args := cmd.Flags().Args()
				fmt.Println("Arguments to be ignored:", args[1:])
			}
		}

		vars, err = varListToMap(buildCmdFlagsVar.Vars)
		if err != nil {
			return errors.New("(command::buildHandler)", "", err)
		}
		persistentVars, err = varListToMap(buildCmdFlagsVar.PersistentVars)
		if err != nil {
			return errors.New("(command::buildHandler)", "", err)
		}

		semverTemplates := config.SemanticVersionTagsTemplates
		if len(buildCmdFlagsVar.SemanticVersionTagsTemplate) > 0 {
			semverTemplates = buildCmdFlagsVar.SemanticVersionTagsTemplate
		}

		pushImages := false
		if !buildCmdFlagsVar.PushImages && config.PushImages {
			pushImages = true
		}

		options := &types.BuildOptions{
			BuilderName:                 buildCmdFlagsVar.BuildBuilderName,
			Cascade:                     buildCmdFlagsVar.BuildOnCascade,
			EnableSemanticVersionTags:   buildCmdFlagsVar.EnableSemanticVersionTags,
			DryRun:                      buildCmdFlagsVar.BuildDryRun,
			PushImages:                  pushImages,
			ConnectionLocal:             buildCmdFlagsVar.ConnectionLocal,
			PersistentVars:              persistentVars,
			Vars:                        vars,
			Tags:                        buildCmdFlagsVar.Tags,
			SemanticVersionTagsTemplate: semverTemplates,
			RegistryNamespace:           buildCmdFlagsVar.ImageRegistryNamespace,
			RegistryHost:                buildCmdFlagsVar.ImageRegistryHost,
			ImageFromName:               buildCmdFlagsVar.ImageFromName,
			ImageFromVersion:            buildCmdFlagsVar.ImageFromVersion,
			ImageFromRegistryNamespace:  buildCmdFlagsVar.ImageFromRegistryNamespace,
			ImageFromRegistryHost:       buildCmdFlagsVar.ImageFromRegistryHost,
		}

		if buildCmdFlagsVar.ImageName != "" {
			// Behavior: Could not override the image name and build all children to avoid unexpected behaviors on the built images names
			if buildCmdFlagsVar.BuildOnCascade {
				return errors.New("(command::buildHandler)", "Could not override image name with build on cascade")
			}
			options.ImageName = buildCmdFlagsVar.ImageName
		}

		// Define num of workers when num workers flag is defined over 0
		if buildCmdFlagsVar.NumWorkers > 0 {
			config.NumWorkers = buildCmdFlagsVar.NumWorkers
		}

		imagesEngine, err = engine.NewImagesEngine(ctx, config.NumWorkers, config.TreePathFile, config.BuilderPathFile)
		if err != nil {
			return errors.New("(command::buildHandler)", "Error creating new image engine", err)
		}

		err = imagesEngine.Build(buildImageName, buildCmdFlagsVar.Versions, options, buildCmdFlagsVar.CascadeDepth)
		if err != nil {
			return errors.New("(command::buildHandler)", fmt.Sprintf("Error building image '%s'", buildImageName), err)
		}

		return nil
	}
}

// someone need to semd vars []string to options
func varListToMap(varsList []string) (map[string]interface{}, error) {

	vars := map[string]interface{}{}

	for _, v := range varsList {
		tokens := strings.Split(v, setVarsSplitToken)

		if len(tokens) != 2 {
			return nil, errors.New("(command::varListToMap)", fmt.Sprintf("Invalid extra variable format on '%v'", v))
		}
		vars[tokens[0]] = tokens[1]
	}

	return vars, nil
}
