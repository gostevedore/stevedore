package build

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	dockerclient "github.com/docker/docker/client"
	"github.com/gostevedore/stevedore/internal/command"
	handler "github.com/gostevedore/stevedore/internal/command/build/handler"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/spf13/cobra"
)

var buildHandler Handlerer

// NewCommand returns a new command to build images
func NewCommand(ctx context.Context, conf *configuration.Configuration) *command.StevedoreCommand {

	handlerOptions := &handler.HandlerOptions{}

	buildCmd := &cobra.Command{
		Use:     "build <image>",
		Short:   "Stevedore command to build images",
		Long:    "Stevedore command to build images",
		Example: "stevedore build ubuntu-base --image-version impish --tag 21.10 --pull-parent-image --push-after-build --remove-local-images-after-push",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			errContext := "(build::PreRunE)"

			dockerClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
			_ = dockerClient

			// drivers

			// service

			// handlers

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			return runeHandler(ctx, buildHandler, handlerOptions)
		},
	}

	// ansible driver flags
	buildCmd.Flags().BoolVar(&handlerOptions.AnsibleConnectionLocal, "connection-local", false, "Use ansible local connection [only applies to ansible-playbook driver]")
	DeprecatedFlagMessageConnectionLocal := "[DEPRECATED FLAG] use 'ansible-connection-local' instead of 'connection-local'"
	DEPRECATEDConnectionLocal := buildCmd.Flags().Bool("connection-local", false, DeprecatedFlagMessageConnectionLocal)
	if *DEPRECATEDConnectionLocal {
		console.Warn(DeprecatedFlagMessageConnectionLocal)
		handlerOptions.AnsibleConnectionLocal = *DEPRECATEDConnectionLocal
	}
	buildCmd.Flags().StringVar(&handlerOptions.AnsibleIntermediateContainerName, "ansible-intermediate-container-name", "", "Name of an intermediate container that can be used during ansible build process [only applies to ansible-playbook driver]")
	DeprecatedFlagMessageBuildBuilderName := "[DEPRECATED FLAG] use 'ansible-intermediate-container-name' instead of 'builder-name'"
	DEPRECATEDBuildBuilderName := buildCmd.Flags().String("builder-name", "", DeprecatedFlagMessageBuildBuilderName)
	if *DEPRECATEDBuildBuilderName != "" {
		console.Warn(DeprecatedFlagMessageBuildBuilderName)
		handlerOptions.AnsibleIntermediateContainerName = *DEPRECATEDBuildBuilderName
	}
	buildCmd.Flags().StringVar(&handlerOptions.AnsibleInventoryPath, "ansible-inventory-path", "", "Specify inventory hosts' path or comma separated list of hosts [only applies to ansible-playbook driver]")
	DeprecatedFlagMessageInventory := "[DEPRECATED FLAG] use 'ansible-inventory-path' instead of 'inventory'"
	DEPRECATEDInventory := buildCmd.Flags().String("inventory", "", DeprecatedFlagMessageInventory)
	if *DEPRECATEDInventory != "" {
		console.Warn(DeprecatedFlagMessageInventory)
		handlerOptions.AnsibleInventoryPath = *DEPRECATEDInventory
	}
	buildCmd.Flags().StringVar(&handlerOptions.AnsibleLimit, "ansible-limit", "", "Further limit selected hosts to an additional pattern [only applies to ansible-playbook driver]")
	DeprecatedFlagMessageLimit := "[DEPRECATED FLAG] use 'ansible-limit' instead of 'limit'"
	DEPRECATEDLimit := buildCmd.Flags().String("limit", "", DeprecatedFlagMessageLimit)
	if *DEPRECATEDLimit != "" {
		console.Warn(DeprecatedFlagMessageLimit)
		handlerOptions.AnsibleInventoryPath = *DEPRECATEDLimit
	}

	// image definition flags
	buildCmd.Flags().StringVarP(&handlerOptions.ImageFromName, "image-from-name", "I", "", "Image parent's name")
	DeprecatedFlagMessageImageFrom := "[DEPRECATED FLAG] use 'image-from-name' instead of 'image-from'"
	DEPRECATEDImageFrom := buildCmd.Flags().String("image-from", "", DeprecatedFlagMessageImageFrom)
	if *DEPRECATEDImageFrom != "" {
		console.Warn(DEPRECATEDImageFrom)
		handlerOptions.ImageFromName = *DEPRECATEDImageFrom
	}
	buildCmd.Flags().StringVarP(&handlerOptions.ImageFromRegistryNamespace, "image-from-namespace", "N", "", "Image parent's registry namespace")
	buildCmd.Flags().StringVarP(&handlerOptions.ImageFromRegistryHost, "image-from-registry", "R", "", "Image parent's registry host")
	buildCmd.Flags().StringVarP(&handlerOptions.ImageFromVersion, "image-from-version", "V", "", "Image parent's version")

	buildCmd.Flags().StringVarP(&handlerOptions.ImageName, "image-name", "i", "", "Image name. Its value overrides the name on the images tree definition")
	buildCmd.Flags().StringVarP(&handlerOptions.ImageRegistryHost, "image-registry-host", "r", "", "Image registry host")
	DeprecatedFlagMessageRegistry := "[DEPRECATED FLAG] use 'image-registry-host' instead of 'registry'"
	DEPRECATEDRegistry := buildCmd.Flags().String("registry", "", DeprecatedFlagMessageRegistry)
	if *DEPRECATEDRegistry != "" {
		console.Warn(DEPRECATEDRegistry)
		handlerOptions.ImageFromRegistryHost = *DEPRECATEDRegistry
	}
	buildCmd.Flags().StringVarP(&handlerOptions.ImageRegistryNamespace, "image-registry-namespace", "n", "", "Image namespace")
	DeprecatedFlagMessageNamespace := "[DEPRECATED FLAG] use 'image-registry-namespace' instead of 'namespace'"
	DEPRECATEDNamespace := buildCmd.Flags().String("namespace", "", DeprecatedFlagMessageNamespace)
	if *DEPRECATEDNamespace != "" {
		console.Warn(DEPRECATEDNamespace)
		handlerOptions.ImageRegistryNamespace = *DEPRECATEDNamespace
	}
	buildCmd.Flags().StringSliceVarP(&handlerOptions.Versions, "image-version", "v", []string{}, "List of versions to build")
	buildCmd.Flags().StringSliceVarP(&handlerOptions.PersistentVars, "persistent-variable", "p", []string{}, "List of persistent variables to set during the build process. Persistent variable that child image inherits from its parent and could not be override. The format of each variable must be <key>=<value>")
	DeprecatedFlagMessageSetPersistent := "[DEPRECATED FLAG] use 'persistent-variable' instead of 'set-persistent'"
	DEPRECATEDSetPersistent := buildCmd.Flags().StringSlice("set-persistent", []string{}, DeprecatedFlagMessageSetPersistent)
	if len(*DEPRECATEDSetPersistent) > 0 {
		console.Warn(DeprecatedFlagMessageSetPersistent)
		handlerOptions.PersistentVars = *DEPRECATEDSetPersistent
	}
	buildCmd.Flags().StringSliceVarP(&handlerOptions.Vars, "variable", "x", []string{}, "Variables to set during the build process. The format of each variable must be <key>=<value>")
	DeprecatedFlagMessageSet := "[DEPRECATED FLAG] use 'variable' instead of 'set'"
	DEPRECATEDSet := buildCmd.Flags().StringSlice("set", []string{}, DeprecatedFlagMessageSet)
	if len(*DEPRECATEDSet) > 0 {
		console.Warn(DeprecatedFlagMessageSet)
		handlerOptions.PersistentVars = *DEPRECATEDSet
	}
	buildCmd.Flags().StringSliceVarP(&handlerOptions.Tags, "tag", "t", []string{}, "List of extra tags to generate")
	buildCmd.Flags().StringSliceVarP(&handlerOptions.Labels, "label", "l", []string{}, "List of labels to assign to the image")
	buildCmd.Flags().StringSliceVarP(&handlerOptions.SemanticVersionTagsTemplates, "semver-tags-template", "T", []string{}, "List of templates to generate tags following semantic version expression")

	// behavior flags
	buildCmd.Flags().BoolVar(&handlerOptions.BuildOnCascade, "build-on-cascade", false, "Build images on cascade. Children's image build is started once the image build finishes")
	DeprecatedFlagMessageCascade := "[DEPRECATED FLAG] use 'build-on-cascade' instead of 'cascade'"
	DEPRECATEDCascade := buildCmd.Flags().Bool("cascade", false, DeprecatedFlagMessageCascade)
	if *DEPRECATEDCascade {
		console.Warn(DeprecatedFlagMessageCascade)
		handlerOptions.BuildOnCascade = *DEPRECATEDCascade
	}
	buildCmd.Flags().IntVar(&handlerOptions.CascadeDepth, "cascade-depth", -1, "Number images levels to build when build on cascade is executed")
	buildCmd.Flags().IntVar(&handlerOptions.Concurrency, "concurrency", 0, "Number of images builds that can be excuted at the same time")
	DeprecatedFlagMessageNumWorkers := "[DEPRECATED FLAG] use 'concurrency' instead of 'num-workers'"
	DEPRECATEDNumWorkers := buildCmd.Flags().Int("num-workers", 0, DeprecatedFlagMessageNumWorkers)
	if *DEPRECATEDNumWorkers > 0 {
		console.Warn(DEPRECATEDNumWorkers)
		handlerOptions.Concurrency = *DEPRECATEDNumWorkers
	}
	// buildCmd.Flags().BoolVar(&handlerOptions.Debug, "debug", false, "Enable debug mode to show build options")
	buildCmd.Flags().BoolVar(&handlerOptions.DryRun, "dry-run", false, "Run build on dry-run mode")
	buildCmd.Flags().BoolVar(&handlerOptions.EnableSemanticVersionTags, "enable-semver-tags", false, "Generate a set of tags for the image based on the semantic version tree when main version is semver 2.0.0 compliance")
	buildCmd.Flags().BoolVar(&handlerOptions.PullParentImage, "pull-parent-image", false, "When is defined parent image is pulled from docker registry")
	buildCmd.Flags().BoolVar(&handlerOptions.PushImagesAfterBuild, "push-after-build", false, "When is defined the image is pushed to docker registry after the build")
	DeprecatedFlagMessagePushImages := "[DEPRECATED FLAG] 'no-push' has no effect because is the default behavior"
	DEPRECATEDPushImages := buildCmd.Flags().Bool("no-push", false, DeprecatedFlagMessagePushImages)
	if *DEPRECATEDPushImages {
		console.Warn(DEPRECATEDPushImages)
	}
	buildCmd.Flags().BoolVar(&handlerOptions.RemoveImagesAfterPush, "remove-local-images-after-push", false, "When is defined images are removed from local after push")

	command := &command.StevedoreCommand{
		Command: buildCmd,
	}

	return command

}

func runeHandler(ctx context.Context, handler Handlerer, options *handler.HandlerOptions) error {

	errContext := "(build::runeHandler)"

	err := handler.Handler(ctx, options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil

}
