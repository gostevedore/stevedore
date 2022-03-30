package promote

import (
	"context"
	"fmt"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/copy"
	dockerclient "github.com/docker/docker/client"
	"github.com/gostevedore/stevedore/internal/command"
	handler "github.com/gostevedore/stevedore/internal/command/promote/handler"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	service "github.com/gostevedore/stevedore/internal/engine/promote"
	repofactory "github.com/gostevedore/stevedore/internal/promote"
	repodocker "github.com/gostevedore/stevedore/internal/promote/docker"
	repodockercopy "github.com/gostevedore/stevedore/internal/promote/docker/promoter"
	repodryrun "github.com/gostevedore/stevedore/internal/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/spf13/cobra"
)

var promoteHandler HandlerPromoter

// NewCommand returns a new command to promote images
func NewCommand(ctx context.Context, conf *configuration.Configuration) *command.StevedoreCommand {

	handlerOptions := &handler.HandlerOptions{}

	promoteCmd := &cobra.Command{
		Use:     "promote",
		Aliases: []string{"publish", "copy"},
		Short:   "Stevedore command to promote, publish or copy images to a docker registry or namespace",
		Long:    "Stevedore command to promote, publish or copy images to a docker registry or namespace",
		Example: "stevedore promote ubuntu:impish --romote-image-registry myregistry.example.com --promote-image-namespace mynamespace",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			errContext := "(promote::PreRunE)"

			dockerClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
			if err != nil {
				return errors.New(errContext, err.Error())
			}

			copyCmd := copy.NewDockerImageCopyCmd(dockerClient)
			copyCmdFacade := repodockercopy.NewDockerCopy(copyCmd)
			promoteRepoDocker := repodocker.NewDockerPromote(copyCmdFacade, os.Stdout)
			promoteRepoDryRun := repodryrun.NewDryRunPromote(copyCmdFacade, os.Stdout)
			promoteRepoFactory := repofactory.NewPromoteFactory()
			err = promoteRepoFactory.Register("docker", promoteRepoDocker)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
			err = promoteRepoFactory.Register("dry-run", promoteRepoDryRun)
			if err != nil {
				return errors.New(errContext, err.Error())
			}

			credentialsStore := credentials.NewCredentialsStore()
			err = credentialsStore.LoadCredentials(conf.DockerCredentialsDir)
			if err != nil {
				return errors.New(errContext, err.Error())
			}

			semverGenerator := semver.NewSemVerGenerator()

			promoteService := service.NewService(promoteRepoFactory, conf, credentialsStore, semverGenerator)

			promoteHandler = handler.NewHandler(promoteService)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			errContext := "(promote::RunE)"
			if cmd.Flags().NArg() == 0 {
				return errors.New(errContext, "Source images name must be provided")
			} else {
				handlerOptions.SourceImageName = cmd.Flags().Arg(0)
				if cmd.Flags().NArg() > 1 {
					args := cmd.Flags().Args()
					fmt.Println("Arguments to be ignored:", args[1:])
				}
			}

			return runeHandler(ctx, promoteHandler, handlerOptions)
		},
	}

	promoteCmd.Flags().BoolVarP(&handlerOptions.EnableSemanticVersionTags, "enable-semver-tags", "S", false, "Generate extra tags based on semantic version tree, when main version is semver 2.0.0 compliance")
	promoteCmd.Flags().StringSliceVarP(&handlerOptions.SemanticVersionTagsTemplates, "semver-tags-template", "T", []string{}, "List templates to generate tags following semantic version expression")
	promoteCmd.Flags().BoolVarP(&handlerOptions.DryRun, "dry-run", "D", false, "Dry run promotion")
	promoteCmd.Flags().StringVarP(&handlerOptions.TargetImageName, "promote-image-name", "i", "", "Target image name")
	promoteCmd.Flags().StringVarP(&handlerOptions.TargetImageRegistryNamespace, "promote-image-namespace", "n", "", "Target image registry mamespace")
	promoteCmd.Flags().StringVarP(&handlerOptions.TargetImageRegistryHost, "promote-image-registry", "r", "", "Target image registry host")
	promoteCmd.Flags().StringSliceVarP(&handlerOptions.TargetImageTags, "promote-image-tag", "t", []string{}, "Target image tag")
	promoteCmd.Flags().BoolVar(&handlerOptions.RemoveTargetImageTags, "remove-local-images-after-push", false, "Remove source image tags")

	promoteCmd.Flags().BoolVar(&handlerOptions.DEPRECATED_RemoveTargetImageTags, "remove-promote-tags", false, "[DEPRECATED] use remove-local-images-after-push. Remove source image tags")

	promoteCmd.Flags().BoolVarP(&handlerOptions.PromoteSourceImageTag, "promote-source-tags", "s", false, "Promote source image. It must be used when a promote image tag is defined and source image needs to be promoted")
	promoteCmd.Flags().BoolVarP(&handlerOptions.RemoteSourceImage, "remote-source-image", "R", false, "Use as a source image an image stored on a Docker registry")

	// Transitorial flags
	if handlerOptions.DEPRECATED_RemoveTargetImageTags && !handlerOptions.RemoveTargetImageTags {
		handlerOptions.RemoveTargetImageTags = handlerOptions.DEPRECATED_RemoveTargetImageTags
		console.Warn("[DEPRECATED FLAG] use `remove-local-images-after-push` instead of `remove-promote-tags`")
	}

	command := &command.StevedoreCommand{
		Command: promoteCmd,
	}

	return command
}

func runeHandler(ctx context.Context, handler HandlerPromoter, options *handler.HandlerOptions) error {

	errContext := "(promote::runeHandler)"

	err := handler.Handler(ctx, options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
