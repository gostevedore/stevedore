package build

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

const (
	DeprecatedFlagMessageConnectionLocal  = "[DEPRECATED FLAG] use 'ansible-connection-local' instead of 'connection-local'"
	DeprecatedFlagMessageBuildBuilderName = "[DEPRECATED FLAG] use 'ansible-intermediate-container-name' instead of 'builder-name'"
	DeprecatedFlagMessageInventory        = "[DEPRECATED FLAG] use 'ansible-inventory-path' instead of 'inventory'"
	DeprecatedFlagMessageLimit            = "[DEPRECATED FLAG] use 'ansible-limit' instead of 'limit'"
	DeprecatedFlagMessageImageFrom        = "[DEPRECATED FLAG] use 'image-from-name' instead of 'image-from'"
	DeprecatedFlagMessageRegistry         = "[DEPRECATED FLAG] use 'image-registry-host' instead of 'registry'"
	DeprecatedFlagMessageNamespace        = "[DEPRECATED FLAG] use 'image-registry-namespace' instead of 'namespace'"
	DeprecatedFlagMessageSetPersistent    = "[DEPRECATED FLAG] use 'persistent-variable' instead of 'set-persistent'"
	DeprecatedFlagMessageSet              = "[DEPRECATED FLAG] use 'variable' instead of 'set'"
	DeprecatedFlagMessageCascade          = "[DEPRECATED FLAG] use 'build-on-cascade' instead of 'cascade'"
	DeprecatedFlagMessageNumWorkers       = "[DEPRECATED FLAG] use 'concurrency' instead of 'num-workers'"
	DeprecatedFlagMessagePushImages       = "[DEPRECATED FLAG] 'no-push' is the stevedore default behavior, use --push-after-build to push image"
)

// NewCommand returns a new command to build images
func NewCommand(ctx context.Context, compatibility Compatibilitier, conf *configuration.Configuration, build Entrypointer) *command.StevedoreCommand {

	buildFlagOptions := &buildFlagOptions{}

	buildCmd := &cobra.Command{
		Use:     "build <image>",
		Short:   "Stevedore command to build images",
		Long:    "Stevedore command to build images",
		Example: "stevedore build ubuntu-base --image-version impish --tag 21.10 --pull-parent-image --push-after-build --remove-local-images-after-push",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::build::RunE)"
			handlerOptions := &handler.Options{}
			entrypointOptions := &entrypoint.Options{}

			if buildFlagOptions.DEPRECATEDConnectionLocal {
				compatibility.AddDeprecated(DeprecatedFlagMessageConnectionLocal)
				buildFlagOptions.AnsibleConnectionLocal = buildFlagOptions.DEPRECATEDConnectionLocal
			}

			if buildFlagOptions.DEPRECATEDBuildBuilderName != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageBuildBuilderName)
				buildFlagOptions.AnsibleIntermediateContainerName = buildFlagOptions.DEPRECATEDBuildBuilderName
			}

			if buildFlagOptions.DEPRECATEDInventory != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageInventory)
				buildFlagOptions.AnsibleInventoryPath = buildFlagOptions.DEPRECATEDInventory
			}

			if buildFlagOptions.DEPRECATEDLimit != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageLimit)
				buildFlagOptions.AnsibleLimit = buildFlagOptions.DEPRECATEDLimit
			}

			if buildFlagOptions.DEPRECATEDImageFrom != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageImageFrom)
				buildFlagOptions.ImageFromName = buildFlagOptions.DEPRECATEDImageFrom
			}

			if buildFlagOptions.DEPRECATEDRegistry != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageRegistry)
				buildFlagOptions.ImageRegistryHost = buildFlagOptions.DEPRECATEDRegistry
			}

			if buildFlagOptions.DEPRECATEDNamespace != "" {
				compatibility.AddDeprecated(DeprecatedFlagMessageNamespace)
				buildFlagOptions.ImageRegistryNamespace = buildFlagOptions.DEPRECATEDNamespace
			}

			if len(buildFlagOptions.DEPRECATEDSetPersistent) > 0 {
				compatibility.AddDeprecated(DeprecatedFlagMessageSetPersistent)
				buildFlagOptions.PersistentVars = append([]string{}, buildFlagOptions.DEPRECATEDSetPersistent...)
			}

			if len(buildFlagOptions.DEPRECATEDSet) > 0 {
				compatibility.AddDeprecated(DeprecatedFlagMessageSet)
				buildFlagOptions.Vars = append([]string{}, buildFlagOptions.DEPRECATEDSet...)
			}

			if buildFlagOptions.DEPRECATEDCascade {
				compatibility.AddDeprecated(DeprecatedFlagMessageCascade)
				buildFlagOptions.BuildOnCascade = buildFlagOptions.DEPRECATEDCascade
			}
			if buildFlagOptions.DEPRECATEDNumWorkers > 0 {
				compatibility.AddDeprecated(DeprecatedFlagMessageNumWorkers)
				buildFlagOptions.Concurrency = buildFlagOptions.DEPRECATEDNumWorkers
			}

			if buildFlagOptions.DEPRECATEDPushImages {
				compatibility.AddDeprecated(DeprecatedFlagMessagePushImages)
			}

			entrypointOptions.Concurrency = buildFlagOptions.Concurrency
			entrypointOptions.Debug = buildFlagOptions.Debug
			entrypointOptions.DryRun = buildFlagOptions.DryRun
			entrypointOptions.UseDockerNormalizedName = buildFlagOptions.UseDockerNormalizedName

			handlerOptions.AnsibleConnectionLocal = buildFlagOptions.AnsibleConnectionLocal
			handlerOptions.AnsibleIntermediateContainerName = buildFlagOptions.AnsibleIntermediateContainerName
			handlerOptions.AnsibleInventoryPath = buildFlagOptions.AnsibleInventoryPath
			handlerOptions.AnsibleLimit = buildFlagOptions.AnsibleLimit
			handlerOptions.BuildOnCascade = buildFlagOptions.BuildOnCascade
			handlerOptions.CascadeDepth = buildFlagOptions.CascadeDepth
			handlerOptions.EnableSemanticVersionTags = buildFlagOptions.EnableSemanticVersionTags
			handlerOptions.ImageFromName = buildFlagOptions.ImageFromName
			handlerOptions.ImageFromRegistryHost = buildFlagOptions.ImageFromRegistryHost
			handlerOptions.ImageFromRegistryNamespace = buildFlagOptions.ImageFromRegistryNamespace
			handlerOptions.ImageFromVersion = buildFlagOptions.ImageFromVersion
			handlerOptions.ImageName = buildFlagOptions.ImageName
			handlerOptions.ImageRegistryHost = buildFlagOptions.ImageRegistryHost
			handlerOptions.ImageRegistryNamespace = buildFlagOptions.ImageRegistryNamespace
			handlerOptions.Versions = append([]string{}, buildFlagOptions.ImageVersions...)
			handlerOptions.Labels = append([]string{}, buildFlagOptions.Labels...)
			handlerOptions.PersistentLabels = append([]string{}, buildFlagOptions.PersistentLabels...)
			handlerOptions.PersistentVars = append([]string{}, buildFlagOptions.PersistentVars...)
			handlerOptions.PullParentImage = buildFlagOptions.PullParentImage
			handlerOptions.PushImagesAfterBuild = buildFlagOptions.PushImagesAfterBuild
			handlerOptions.RemoveImagesAfterPush = buildFlagOptions.RemoveImagesAfterPush
			handlerOptions.SemanticVersionTagsTemplates = append([]string{}, buildFlagOptions.SemanticVersionTagsTemplates...)
			handlerOptions.Tags = append([]string{}, buildFlagOptions.Tags...)
			handlerOptions.Vars = append([]string{}, buildFlagOptions.Vars...)

			err = build.Execute(ctx, cmd.Flags().Args(), conf, entrypointOptions, handlerOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	buildCmd.Flags().BoolVar(&buildFlagOptions.DEPRECATEDCascade, "cascade", false, DeprecatedFlagMessageCascade)
	buildCmd.Flags().BoolVar(&buildFlagOptions.DEPRECATEDConnectionLocal, "connection-local", false, DeprecatedFlagMessageConnectionLocal)
	buildCmd.Flags().IntVar(&buildFlagOptions.DEPRECATEDNumWorkers, "num-workers", 0, DeprecatedFlagMessageNumWorkers)
	buildCmd.Flags().StringSliceVar(&buildFlagOptions.DEPRECATEDSet, "set", []string{}, DeprecatedFlagMessageSet)
	buildCmd.Flags().StringSliceVar(&buildFlagOptions.DEPRECATEDSetPersistent, "set-persistent", []string{}, DeprecatedFlagMessageSetPersistent)
	buildCmd.Flags().StringVar(&buildFlagOptions.DEPRECATEDBuildBuilderName, "builder-name", "", DeprecatedFlagMessageBuildBuilderName)
	buildCmd.Flags().StringVar(&buildFlagOptions.DEPRECATEDImageFrom, "image-from", "", DeprecatedFlagMessageImageFrom)
	buildCmd.Flags().StringVar(&buildFlagOptions.DEPRECATEDInventory, "inventory", "", DeprecatedFlagMessageInventory)
	buildCmd.Flags().StringVar(&buildFlagOptions.DEPRECATEDLimit, "limit", "", DeprecatedFlagMessageLimit)
	buildCmd.Flags().StringVar(&buildFlagOptions.DEPRECATEDNamespace, "namespace", "", DeprecatedFlagMessageNamespace)
	buildCmd.Flags().StringVar(&buildFlagOptions.DEPRECATEDRegistry, "registry", "", DeprecatedFlagMessageRegistry)

	// ansible driver flags
	buildCmd.Flags().BoolVar(&buildFlagOptions.AnsibleConnectionLocal, "ansible-connection-local", false, "When this flag is enabled, ansible uses local connection [only applies to ansible-playbook driver]")
	buildCmd.Flags().StringVar(&buildFlagOptions.AnsibleIntermediateContainerName, "ansible-intermediate-container-name", "", "Name of an intermediate container that can be used during ansible build process [only applies to ansible-playbook driver]")
	buildCmd.Flags().StringVar(&buildFlagOptions.AnsibleInventoryPath, "ansible-inventory-path", "", "Specify inventory hosts' path or comma separated list of hosts [only applies to ansible-playbook driver]")
	buildCmd.Flags().StringVar(&buildFlagOptions.AnsibleLimit, "ansible-limit", "", "Further limit selected hosts to an additional pattern [only applies to ansible-playbook driver]")

	// image definition flags
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.ImageVersions, "image-version", "v", []string{}, "List of versions to build")
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.Labels, "label", "l", []string{}, "List of labels to assign to the image")
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.PersistentVars, "persistent-variable", "p", []string{}, "List of persistent variables to set during the build process. Persistent variable inherited from parent has precedence over the current ones. The format of each variable must be <key>=<value>")
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.PersistentLabels, "persistent-label", "L", []string{}, "List of persistent labels to set during the build process. Persistent labels inherited from parent has precedence over the current ones. The format of each variable must be <key>=<value>")
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.SemanticVersionTagsTemplates, "semver-tags-template", "T", []string{}, "List of templates to generate tags following semantic version expression")
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.Tags, "tag", "t", []string{}, "List of extra tags to generate")
	buildCmd.Flags().StringSliceVarP(&buildFlagOptions.Vars, "variable", "x", []string{}, "Variables to set during the build process. The format of each variable must be <key>=<value>")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageFromName, "image-from-name", "I", image.UndefinedStringValue, "Image parent's name")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageFromRegistryHost, "image-from-registry", "R", image.UndefinedStringValue, "Image parent's registry host")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageFromRegistryNamespace, "image-from-namespace", "N", image.UndefinedStringValue, "Image parent's registry namespace")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageFromVersion, "image-from-version", "V", image.UndefinedStringValue, "Image parent's version")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageName, "image-name", "i", image.UndefinedStringValue, "Image name. Its value overrides the name on the images tree definition")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageRegistryHost, "image-registry-host", "r", image.UndefinedStringValue, "Image registry host")
	buildCmd.Flags().StringVarP(&buildFlagOptions.ImageRegistryNamespace, "image-registry-namespace", "n", image.UndefinedStringValue, "Image namespace")
	buildCmd.Flags().BoolVar(&buildFlagOptions.UseDockerNormalizedName, "use-docker-normalized-name", false, "Use Docker normalized name references")

	// behavior flags
	buildCmd.Flags().BoolVar(&buildFlagOptions.BuildOnCascade, "build-on-cascade", false, "When this flag is enabled, children images are also built")
	buildCmd.Flags().IntVar(&buildFlagOptions.CascadeDepth, "cascade-depth", -1, "Number children levels to build when build on cascade is executed")
	buildCmd.Flags().IntVar(&buildFlagOptions.Concurrency, "concurrency", 0, "Number of images builds that can be excuted at the same time")

	// buildCmd.Flags().BoolVar(&buildFlagOptions.Debug, "debug", false, "Enable debug mode to show build options")
	buildCmd.Flags().BoolVar(&buildFlagOptions.DEPRECATEDPushImages, "no-push", false, DeprecatedFlagMessagePushImages)
	buildCmd.Flags().BoolVar(&buildFlagOptions.DryRun, "dry-run", false, "When this flag is enabled, the built is executed in dry-run mode")
	buildCmd.Flags().BoolVar(&buildFlagOptions.EnableSemanticVersionTags, "enable-semver-tags", false, "When this flag is enabled, and main version is semver 2.0.0 compliance extra tag are created based on the semantic version tree")
	buildCmd.Flags().BoolVar(&buildFlagOptions.PullParentImage, "pull-parent-image", false, "When this flag is enabled, parent image is pulled from docker registry")
	buildCmd.Flags().BoolVar(&buildFlagOptions.PushImagesAfterBuild, "push-after-build", false, "When this flag is enabled, the image is pushed to docker registry after the build")
	buildCmd.Flags().BoolVar(&buildFlagOptions.RemoveImagesAfterPush, "remove-local-images-after-push", false, "When this flag is enabled, images are removed from local after push")

	command := &command.StevedoreCommand{
		Command: buildCmd,
	}

	return command
}
