package dockerpromoter

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	dockercopy "github.com/apenella/go-docker-builder/pkg/copy"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/client"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

const (
	DockerImageFilterReference = "reference"
)

func Promote(ctx context.Context, options *types.PromoteOptions) error {

	var err error
	var promoteImageURL *image.ImageURL
	var sourceImageName, promoteImageName string
	// var normalizedPromoteImageName reference.Named
	var auth *credentials.RegistryUserPassAuth
	var tags []string

	if options.ImageName == "" {
		return errors.New("(promote::Promote)", "Image name must be defined on promote options")
	}

	promoteImageURL, err = image.Parse(options.ImageName)
	if err != nil {
		return errors.New("(promote::Promote)", fmt.Sprintf("Error when parsing '%s'", options.ImageName), err)
	}

	sourceImageName, err = promoteImageURL.URL()
	if err != nil {
		return errors.New("(promote::Promote)", "Error when achiving image URL", err)
	}

	if options.ImagePromoteRegistryHost != "" {
		promoteImageURL.Registry = options.ImagePromoteRegistryHost
	}
	if options.ImagePromoteRegistryNamespace != "" {
		promoteImageURL.Namespace = options.ImagePromoteRegistryNamespace
	}
	if options.ImagePromoteName != "" {
		promoteImageURL.Name = options.ImagePromoteName
	}

	if promoteImageURL.Registry != "" {
		auth, err = credentials.AchieveCredential(promoteImageURL.Registry)
		if err != nil {
			auth = nil
		}
	}

	// when no tags are defined use the source image tag
	if len(options.ImagePromoteTags) == 0 {
		tags = append(options.ImagePromoteTags, promoteImageURL.Tag)
	} else {
		tags = options.ImagePromoteTags
	}

	// pushSource controls when the source image is going to be pushed to registry and forces it to be pushed such the last image
	// It also solve an issue when is defined a --promote-image-tag with same value as the source tag and --remove-promote-tags flag is enabled. That flags combination removes the source image and the upcomming image tags won't be pushed because the source image was already removed
	pushSource := false
	for _, tag := range tags {
		promoteImageURL.Tag = tag
		promoteImageName, err = promoteImageURL.URL()
		if err != nil {
			return errors.New("(promote::Promote)", "Error when achiving image URL", err)
		}

		if sourceImageName != promoteImageName {
			err = promoteWorker(ctx, options, sourceImageName, promoteImageName, auth)
			if err != nil {
				return errors.New("(promote::Promote) ", fmt.Sprintf("Error promoting '%s' to '%s'", sourceImageName, promoteImageName), err)
			}
		} else {
			pushSource = true
		}
	}

	if pushSource {
		err = promoteWorker(ctx, options, sourceImageName, sourceImageName, auth)
		if err != nil {
			return errors.New("(promote::Promote) ", fmt.Sprintf("Error promoting '%s' to '%s'", sourceImageName, promoteImageName), err)
		}
	}

	return nil
}

func promoteWorker(ctx context.Context, options *types.PromoteOptions, src, dest string, credentials *credentials.RegistryUserPassAuth) error {

	var err error
	var dockerCli *client.Client
	var normalizedPromoteImageName reference.Named

	normalizedPromoteImageName, err = reference.ParseNormalizedNamed(dest)
	if err != nil {
		return errors.New("(promote::Promote)", fmt.Sprintf("Error normalizing image name '%s'", dest), err)
	}

	if options.OutputPrefix == "" {
		options.OutputPrefix = normalizedPromoteImageName.String()
	}

	// create docker sdk client
	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("(promote::Promote)", "Error on docker client creation", err)
	}
	dockerCli.NegotiateAPIVersion(ctx)

	// create docker push instance
	dockerCopier := &dockercopy.DockerImageCopyCmd{
		Cli: dockerCli,
		Response: response.NewDefaultResponse(
			response.WithTransformers(
				transformer.Prepend(options.OutputPrefix),
			),
			response.WithWriter(console.GetConsole()),
		),
		SourceImage:     src,
		TargetImage:     normalizedPromoteImageName.String(),
		RemoveAfterPush: options.RemovePromotedTags,
		RemoteSource:    false,
	}

	// it just work with local images when remote source is accepted credentials must be updated
	if credentials != nil {
		user := credentials.Username
		pass := credentials.Password
		dockerCopier.AddAuth(user, pass)
	}

	err = dockerCopier.Run(ctx)
	if err != nil {
		return errors.New("(promote::Promote)", fmt.Sprintf("Error pushing '%s'", normalizedPromoteImageName.String()), err)
	}

	return nil
}
