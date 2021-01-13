package promote

import (
	"context"
	"fmt"
	dockerpromoter "stevedore/internal/promote/promoter/docker"
	dryrunpromoter "stevedore/internal/promote/promoter/dryrun"
	"stevedore/internal/types"

	errors "github.com/apenella/go-common-utils/error"
)

func Promote(ctx context.Context, options *types.PromoteOptions) error {
	var err error
	if options.DryRun {
		err = dryrunpromoter.Promote(ctx, options)
	} else {
		err = dockerpromoter.Promote(ctx, options)
	}

	if err != nil {
		return errors.New("(promote::Promote)", fmt.Sprintf("Error promoting '%s' image", options.ImageName), err)
	}

	return nil
}

// func Promote(ctx *context.Context, options *types.PromoteOptions) error {

// 	var err error
// 	var promoteImageURL *image.ImageURL
// 	var sourceImageName, promoteImageName string
// 	// var normalizedPromoteImageName reference.Named
// 	var credentials *credentials.RegistryUserPassAuth
// 	var tags []string

// 	if options.ImageName == "" {
// 		return errors.New("(promote::Promote) Image name must be defined on promote options")
// 	}

// 	promoteImageURL, err = image.Parse(options.ImageName)
// 	sourceImageName, err = promoteImageURL.URL()
// 	if err != nil {
// 		return errors.New("(promote::Promote) Error when achiving image URL")
// 	}

// 	if options.ImagePromoteRegistryHost != "" {
// 		promoteImageURL.Registry = options.ImagePromoteRegistryHost
// 	}
// 	if options.ImagePromoteRegistryNamespace != "" {
// 		promoteImageURL.Namespace = options.ImagePromoteRegistryNamespace
// 	}
// 	if options.ImagePromoteName != "" {
// 		promoteImageURL.Name = options.ImagePromoteName
// 	}

// 	if promoteImageURL.Registry != "" {
// 		credentials, err = ctx.Credentials.AchieveCredential(promoteImageURL.Registry)
// 		if err != nil {
// 			credentials = nil
// 		}
// 	}

// 	// when no tags are defined use the source image tag
// 	if len(options.ImagePromoteTags) == 0 {
// 		tags = append(options.ImagePromoteTags, promoteImageURL.Tag)
// 	} else {
// 		tags = options.ImagePromoteTags
// 	}

// 	for _, tag := range tags {
// 		promoteImageURL.Tag = tag
// 		promoteImageName, err = promoteImageURL.URL()
// 		if err != nil {
// 			return errors.New("(promote::Promote) Error when achiving image URL")
// 		}

// 		err = promoteImage(ctx, options, sourceImageName, promoteImageName, credentials)
// 		if err != nil {
// 			return errors.New("(promote::Promote) Error promoting " + sourceImageName + " to " + promoteImageName)
// 		}
// 	}

// 	return nil
// }

// func promoteImage(ctx *context.Context, options *types.PromoteOptions, src, dest string, credentials *credentials.RegistryUserPassAuth) error {

// 	var err error
// 	var dockerCli *client.Client
// 	var normalizedPromoteImageName reference.Named

// 	normalizedPromoteImageName, err = reference.ParseNormalizedNamed(dest)
// 	if err != nil {
// 		return errors.New("(promote::Promote) Error normalizing image name '" + dest + "'. " + err.Error())
// 	}

// 	if options.OutputPrefix == "" {
// 		options.OutputPrefix = normalizedPromoteImageName.String()
// 	}

// 	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		return errors.New("(promote::Promote) Error on docker client creation. " + err.Error())
// 	}

// 	dockerPushOptions := &dockerpush.DockerPushOptions{
// 		ImageName: normalizedPromoteImageName.String(),
// 	}

// 	if credentials != nil {
// 		user := credentials.Username
// 		pass := credentials.Password
// 		dockerPushOptions.AddAuth(user, pass)

// 		// add auth to build options when it not already set
// 		if dockerPushOptions.RegistryAuth == nil {
// 			dockerPushOptions.AddAuth(user, pass)
// 		}
// 	}

// 	err = dockerCli.ImageTag(ctx.Ctx, src, normalizedPromoteImageName.String())
// 	if err != nil {
// 		return errors.New("(promote::Promote) Error When tagging '" + src + "' to '" + normalizedPromoteImageName.String() + "'. " + err.Error())
// 	}

// 	dockerPusher := &dockerpush.DockerPushCmd{
// 		Writer:            ctx.Writer,
// 		Cli:               dockerCli,
// 		Context:           ctx.Ctx,
// 		DockerPushOptions: dockerPushOptions,
// 		ExecPrefix:        options.OutputPrefix,
// 	}

// 	err = dockerPusher.Run()
// 	if err != nil {
// 		return errors.New("(promote::Promote) Error pushing '" + normalizedPromoteImageName.String() + "'. " + err.Error())
// 	}

// 	if options.RemovePromotedTags {
// 		dockerCli.ImageRemove(ctx.Ctx, normalizedPromoteImageName.String(), dockertypes.ImageRemoveOptions{
// 			Force:         true,
// 			PruneChildren: true,
// 		})
// 	}

// 	return nil
// }
