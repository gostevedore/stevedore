package dockerpromoter

import (
	"context"
	"fmt"
	"stevedore/internal/credentials"
	"stevedore/internal/image"
	"stevedore/internal/types"
	"stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	dockerpush "github.com/apenella/go-docker-builder/pkg/push"
	"github.com/docker/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

	for _, tag := range tags {
		promoteImageURL.Tag = tag
		promoteImageName, err = promoteImageURL.URL()
		if err != nil {
			return errors.New("(promote::Promote)", "Error when achiving image URL", err)
		}

		err = promoteWorker(ctx, options, sourceImageName, promoteImageName, auth)
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
		return errors.New("(promote::Promote)", "Error normalizing image name '"+dest+"'", err)
	}

	if options.OutputPrefix == "" {
		options.OutputPrefix = normalizedPromoteImageName.String()
	}

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("(promote::Promote)", "Error on docker client creation", err)
	}
	dockerCli.NegotiateAPIVersion(ctx)

	dockerPushOptions := &dockerpush.DockerPushOptions{
		ImageName: normalizedPromoteImageName.String(),
	}

	if credentials != nil {
		user := credentials.Username
		pass := credentials.Password
		dockerPushOptions.AddAuth(user, pass)

		// add auth to build options when it not already set
		if dockerPushOptions.RegistryAuth == nil {
			dockerPushOptions.AddAuth(user, pass)
		}
	}

	err = dockerCli.ImageTag(ctx, src, normalizedPromoteImageName.String())
	if err != nil {
		return errors.New("(promote::Promote)", "Error tagging '"+src+"' to '"+normalizedPromoteImageName.String()+"'", err)
	}

	dockerPusher := &dockerpush.DockerPushCmd{
		Writer:            console.GetConsole(),
		Cli:               dockerCli,
		Context:           ctx,
		DockerPushOptions: dockerPushOptions,
		ExecPrefix:        options.OutputPrefix,
	}

	err = dockerPusher.Run()
	if err != nil {
		return errors.New("(promote::Promote)", "Error pushing '"+normalizedPromoteImageName.String()+"'", err)
	}

	if options.RemovePromotedTags {
		dockerCli.ImageRemove(ctx, normalizedPromoteImageName.String(), dockertypes.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
	}

	return nil
}
