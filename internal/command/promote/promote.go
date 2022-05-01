package promote

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/spf13/cobra"
)

const (
	DeprecatedFlagMessageRemoveTargetImageTags = "[DEPRECATED FLAG] use `remove-local-images-after-push` instead of `remove-promote-tags`"
)

// NewCommand returns a new command to promote images
func NewCommand(ctx context.Context, compatibility Compatibilitier, conf *configuration.Configuration, promote Entrypointer) *command.StevedoreCommand {

	promoteFlagOptions := &promoteFlagOptions{}

	promoteCmd := &cobra.Command{
		Use:     "promote",
		Aliases: []string{"publish", "copy"},
		Short:   "Stevedore command to promote, publish or copy images to a docker registry or namespace",
		Long:    "Stevedore command to promote, publish or copy images to a docker registry or namespace",
		Example: "stevedore promote ubuntu:impish --romote-image-registry myregistry.example.com --promote-image-namespace mynamespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			errContext := "(promote::RunE)"
			handlerOptions := &handler.Options{}

			fmt.Println(">>>>", promoteFlagOptions.DEPRECATEDRemoveTargetImageTags)

			// Transitorial flags
			if promoteFlagOptions.DEPRECATEDRemoveTargetImageTags && !handlerOptions.RemoveTargetImageTags {
				promoteFlagOptions.RemoveTargetImageTags = promoteFlagOptions.DEPRECATEDRemoveTargetImageTags
				compatibility.AddDeprecated(DeprecatedFlagMessageRemoveTargetImageTags)
			}

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

			err = promote.Execute(ctx, cmd.Flags().Args(), conf, handlerOptions)
			if err != nil {
				return errors.New(errContext, err.Error())
			}

			return nil
		},
	}

	promoteCmd.Flags().BoolVar(&promoteFlagOptions.DEPRECATEDRemoveTargetImageTags, "remove-promote-tags", false, DeprecatedFlagMessageRemoveTargetImageTags)

	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.EnableSemanticVersionTags, "enable-semver-tags", "S", false, "Generate extra tags based on semantic version tree, when main version is semver 2.0.0 compliance")
	promoteCmd.Flags().StringSliceVarP(&promoteFlagOptions.SemanticVersionTagsTemplates, "semver-tags-template", "T", []string{}, "List templates to generate tags following semantic version expression")
	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.DryRun, "dry-run", "D", false, "Dry run promotion")
	promoteCmd.Flags().StringVarP(&promoteFlagOptions.TargetImageName, "promote-image-name", "i", "", "Target image name")
	promoteCmd.Flags().StringVarP(&promoteFlagOptions.TargetImageRegistryNamespace, "promote-image-registry-namespace", "n", "", "Target image registry mamespace")
	promoteCmd.Flags().StringVarP(&promoteFlagOptions.TargetImageRegistryHost, "promote-image-registry-host", "r", "", "Target image registry host")
	promoteCmd.Flags().StringSliceVarP(&promoteFlagOptions.TargetImageTags, "promote-image-tag", "t", []string{}, "Target image tag")
	promoteCmd.Flags().BoolVar(&promoteFlagOptions.RemoveTargetImageTags, "remove-local-images-after-push", false, "Remove source image tags")
	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.PromoteSourceImageTag, "force-promote-source-image", "s", false, "Force to promote source image tag, although promote-image-tag is set")
	promoteCmd.Flags().BoolVarP(&promoteFlagOptions.RemoteSourceImage, "image-from-remote-source", "R", false, "Promote an image stored on a Docker registry")

	command := &command.StevedoreCommand{
		Command: promoteCmd,
	}

	return command
}
