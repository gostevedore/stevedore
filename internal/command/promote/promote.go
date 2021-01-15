package promote

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/command"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/engine"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/spf13/cobra"
)

type promoteCmdFlags struct {
	DryRun                        bool
	EnableSemanticVersionTags     bool
	ImageName                     string
	ImagePromoteName              string
	ImagePromoteRegistryNamespace string
	ImagePromoteRegistryHost      string
	ImagePromoteTags              []string
	RemovePromoteTags             bool
	SemanticVersionTagsTemplate   []string
}

var promoteCmdFlagsVar *promoteCmdFlags

//  NewCommand return an stevedore command object for dev
func NewCommand(ctx context.Context, config *configuration.Configuration) *command.StevedoreCommand {

	promoteCmdFlagsVar = &promoteCmdFlags{}

	promoteCmd := &cobra.Command{
		Use:     "promote",
		Aliases: []string{"publish"},
		Short:   "Promote images",
		Long:    "",
		Hidden:  true,
		RunE:    promoteHandler(ctx, config),
	}

	promoteCmd.Flags().BoolVarP(&promoteCmdFlagsVar.EnableSemanticVersionTags, "enable-semver-tags", "S", false, "Generate extra tags based on semantic version tree when main version is semver 2.0.0 compliance")
	promoteCmd.Flags().StringSliceVarP(&promoteCmdFlagsVar.SemanticVersionTagsTemplate, "semver-tags-template", "T", []string{}, "List templates to generate tags following semantic version expression")
	promoteCmd.Flags().BoolVarP(&promoteCmdFlagsVar.DryRun, "dry-run", "D", false, "Dry run show the promote parameters")
	promoteCmd.Flags().StringVarP(&promoteCmdFlagsVar.ImagePromoteName, "promote-image-name", "i", "", "Name for the image to be promoted")
	promoteCmd.Flags().StringVarP(&promoteCmdFlagsVar.ImagePromoteRegistryNamespace, "promote-image-namespace", "n", "", "Registry's mamespace for the image to be promoted")
	promoteCmd.Flags().StringVarP(&promoteCmdFlagsVar.ImagePromoteRegistryHost, "promote-image-registry", "r", "", "Registry's host for the image to be promoted")
	promoteCmd.Flags().StringSliceVarP(&promoteCmdFlagsVar.ImagePromoteTags, "promote-image-tag", "t", []string{}, "Extra tag for the image to be promoted")
	promoteCmd.Flags().BoolVarP(&promoteCmdFlagsVar.RemovePromoteTags, "remove-promote-tags", "R", false, "Remove remoted tags from local docker host")

	command := &command.StevedoreCommand{
		Command: promoteCmd,
	}

	return command
}

func promoteHandler(ctx context.Context, config *configuration.Configuration) command.CobraRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		var imagesEngine *engine.ImagesEngine

		if cmd.Flags().NArg() == 0 {
			return errors.New("(command::promoteHandler)", "Is required an image name")
		} else {
			promoteCmdFlagsVar.ImageName = cmd.Flags().Arg(0)
			if cmd.Flags().NArg() > 1 {
				args := cmd.Flags().Args()
				fmt.Println("Arguments to be ignored:", args[1:])
			}
		}

		options := &types.PromoteOptions{
			DryRun:                      promoteCmdFlagsVar.DryRun,
			EnableSemanticVersionTags:   promoteCmdFlagsVar.EnableSemanticVersionTags,
			ImageName:                   promoteCmdFlagsVar.ImageName,
			RemovePromotedTags:          promoteCmdFlagsVar.RemovePromoteTags,
			SemanticVersionTagsTemplate: promoteCmdFlagsVar.SemanticVersionTagsTemplate,
		}

		if promoteCmdFlagsVar.ImagePromoteName != "" {
			options.ImagePromoteName = promoteCmdFlagsVar.ImagePromoteName
		}

		if promoteCmdFlagsVar.ImagePromoteRegistryNamespace != "" {
			options.ImagePromoteRegistryNamespace = promoteCmdFlagsVar.ImagePromoteRegistryNamespace
		}

		if promoteCmdFlagsVar.ImagePromoteRegistryHost != "" {
			options.ImagePromoteRegistryHost = promoteCmdFlagsVar.ImagePromoteRegistryHost
		}

		if len(promoteCmdFlagsVar.ImagePromoteTags) > 0 {
			options.ImagePromoteTags = promoteCmdFlagsVar.ImagePromoteTags
		}

		imagesEngine, err = engine.NewImagesEngine(ctx, config.NumWorkers, config.TreePathFile, config.BuilderPathFile)
		if err != nil {
			return errors.New("(command::promoteHandler)", "Error creating new image engine", err)
		}

		err = imagesEngine.Promote(options)
		if err != nil {
			return errors.New("(command::promoteHandler)", fmt.Sprintf("Error promoting image '%s'", promoteCmdFlagsVar.ImageName), err)
		}

		return nil
	}
}
