package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/cobra"
)

const (
	DeprecatedFlagMessageRemoveTargetImageTags        = "[DEPRECATED FLAG] use `remove-local-images-after-push` instead of `remove-promote-tags`"
	DeprecatedFlagMessageTargetImageRegistryHost      = "[DEPRECATED FLAG] use `promote-image-registry-host` instead of `promote-image-registry`"
	DeprecatedFlagMessageTargetImageRegistryNamespace = "[DEPRECATED FLAG] use `promote-image-registry-namespace` instead of `promote-image-namespace`"
)

// NewCommand returns a new command to promote images
func NewCommand(ctx context.Context, compatibility Compatibilitier, conf *configuration.Configuration, promote Entrypointer) *command.StevedoreCommand {

	promoteFlagOptions := &promoteFlagOptions{}

	promoteCmd := &cobra.Command{
		Use:     "promote",
		Aliases: []string{"publish", "copy"},
		Short:   "Stevedore command to promote, publish or copy images to a docker registry or namespace",
		Long:    "Stevedore command to promote, publish or copy images to a docker registry or namespace",
		Example: "stevedore promote ubuntu:impish --promote-image-registry myregistry.example.com --promote-image-namespace mynamespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(cli::promote::RunE)"
			handlerOptions := &handler.Options{}
			entrypointOptions := &entrypoint.Options{}

			// Transitorial flags
			if promoteFlagOptions.DEPRECATEDRemoveTargetImageTags {
				compatibility.AddDeprecated(DeprecatedFlagMessageRemoveTargetImageTags)
				if !promoteFlagOptions.RemoveTargetImageTags {
					promoteFlagOptions.RemoveTargetImageTags = promoteFlagOptions.DEPRECATEDRemoveTargetImageTags
				}
			}

			if promoteFlagOptions.DEPRECATEDTargetImageRegistryHost != image.UndefinedStringValue {
				compatibility.AddDeprecated(DeprecatedFlagMessageTargetImageRegistryHost)
				if promoteFlagOptions.TargetImageRegistryHost == image.UndefinedStringValue {
					promoteFlagOptions.TargetImageRegistryHost = promoteFlagOptions.DEPRECATEDTargetImageRegistryHost
				}
			}

			if promoteFlagOptions.DEPRECATEDTargetImageRegistryNamespace != image.UndefinedStringValue {
				compatibility.AddDeprecated(DeprecatedFlagMessageTargetImageRegistryNamespace)
				if promoteFlagOptions.TargetImageRegistryNamespace == image.UndefinedStringValue {
					promoteFlagOptions.TargetImageRegistryNamespace = promoteFlagOptions.DEPRECATEDTargetImageRegistryNamespace
				}
			}

			entrypointOptions.UseDockerNormalizedName = promoteFlagOptions.UseDockerNormalizedName

			handlerOptions.DryRun = promoteFlagOptions.DryRun
			handlerOptions.EnableSemanticVersionTags = promoteFlagOptions.EnableSemanticVersionTags
			handlerOptions.TargetImageName = promoteFlagOptions.TargetImageName
			handlerOptions.TargetImageRegistryNamespace = promoteFlagOptions.TargetImageRegistryNamespace
			handlerOptions.TargetImageRegistryHost = promoteFlagOptions.TargetImageRegistryHost
			handlerOptions.TargetImageTags = append([]string{}, promoteFlagOptions.TargetImageTags...)
			handlerOptions.RemoveTargetImageTags = promoteFlagOptions.RemoveTargetImageTags
			handlerOptions.SemanticVersionTagsTemplates = append([]string{}, promoteFlagOptions.SemanticVersionTagsTemplates...)
			handlerOptions.PromoteSourceImageTag = promoteFlagOptions.PromoteSourceImageTag
			handlerOptions.RemoteSourceImage = promoteFlagOptions.RemoteSourceImage

			err = promote.Execute(ctx, cmd.Flags().Args(), conf, entrypointOptions, handlerOptions)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			return nil
		},
	}

	promoteCmd.Flags().BoolVar(&promoteFlagOptions.DEPRECATEDRemoveTargetImageTags, "remove-promote-tags", false, DeprecatedFlagMessageRemoveTargetImageTags)
	promoteCmd.Flags().StringVar(&promoteFlagOptions.DEPRECATEDTargetImageRegistryHost, "promote-image-host", image.UndefinedStringValue, DeprecatedFlagMessageTargetImageRegistryNamespace)
	promoteCmd.Flags().StringVar(&promoteFlagOptions.DEPRECATEDTargetImageRegistryNamespace, "promote-image-namespace", image.UndefinedStringValue, DeprecatedFlagMessageTargetImageRegistryNamespace)

	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.EnableSemanticVersionTags, "enable-semver-tags", "S", false, "When this flag is enabled, and main version is semver 2.0.0 compliance extra tag are created based on the semantic version tree")
	promoteCmd.Flags().StringSliceVarP(&promoteFlagOptions.SemanticVersionTagsTemplates, "semver-tags-template", "T", []string{}, "List templates to generate tags following semantic version expression")
	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.DryRun, "dry-run", "D", false, "Dry run promotion")
	// using UndefinedStringValue rather than the empty string let you to overwrite those values with an empty value
	promoteCmd.Flags().StringVarP(&promoteFlagOptions.TargetImageName, "promote-image-name", "i", image.UndefinedStringValue, "Target image name")
	// using UndefinedStringValue rather than the empty string let you to overwrite those values with an empty value
	promoteCmd.Flags().StringVarP(&promoteFlagOptions.TargetImageRegistryNamespace, "promote-image-registry-namespace", "n", image.UndefinedStringValue, "Target image registry mamespace")
	// using UndefinedStringValue rather than the empty string let you to overwrite those values with an empty value
	promoteCmd.Flags().StringVarP(&promoteFlagOptions.TargetImageRegistryHost, "promote-image-registry-host", "r", image.UndefinedStringValue, "Target image registry host")
	promoteCmd.Flags().StringSliceVarP(&promoteFlagOptions.TargetImageTags, "promote-image-tag", "t", []string{}, "List of target image tags")
	promoteCmd.Flags().BoolVar(&promoteFlagOptions.RemoveTargetImageTags, "remove-local-images-after-push", false, "When this flag is enabled, images are removed from local after push")
	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.PromoteSourceImageTag, "force-promote-source-image", "s", false, "When this flag is enabled, the source image is also promoted, along with any other target image")
	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.RemoteSourceImage, "use-source-image-from-remote", "R", false, "When this flag is enabled, source images is downloaded from remote Docker registry")
	promoteCmd.Flags().BoolVar(&promoteFlagOptions.UseDockerNormalizedName, "use-docker-normalized-name", false, "Use Docker normalized name references")

	command := &command.StevedoreCommand{
		Command: promoteCmd,
	}

	return command
}
