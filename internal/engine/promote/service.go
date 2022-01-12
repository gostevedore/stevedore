package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/promote"
)

type Service struct {
	credentials   CredentialsStorer
	factory       PromoteFactorier
	semver        Semverser
	configuration *configuration.Configuration
}

func NewService(f PromoteFactorier, conf *configuration.Configuration, c CredentialsStorer, s Semverser) *Service {
	return &Service{
		credentials:   c,
		factory:       f,
		semver:        s,
		configuration: conf,
	}
}

// Promote an image
func (e *Service) Promote(ctx context.Context, options *ServiceOptions, promoteType string) error {

	var err error
	var sourceImage, targetImage *image.Image

	promoteOptions := &promote.PromoteOptions{}
	errContext := "(Service::Promote)"

	if e.factory == nil {
		return errors.New(errContext, "Promote factory has not been initialized")
	}

	if e.configuration == nil {
		return errors.New(errContext, "Configuration has not been initialized")
	}

	if e.semver == nil {
		return errors.New(errContext, "Semver has not been initialized")
	}

	if e.credentials == nil {
		return errors.New(errContext, "Credentials has not been initialized")
	}

	if options == nil {
		return errors.New(errContext, "Options are required on promote service")
	}

	if options.SourceImageName == "" {
		return errors.New(errContext, "Promote options requires an image source name defined")
	}

	promoteOptions.SourceImageName = options.SourceImageName
	sourceImage, err = image.Parse(options.SourceImageName)
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	targetImage = sourceImage

	pullAuth := e.getCredentials(sourceImage.RegistryHost)
	if pullAuth != nil {
		promoteOptions.PullAuthUsername = pullAuth.Username
		promoteOptions.PullAuthPassword = pullAuth.Password
	}

	if options.PromoteSourceImageTag {
		promoteOptions.TargetImageTags = append(promoteOptions.TargetImageTags, sourceImage.Version)
	}

	if options.TargetImageRegistryHost != "" {
		targetImage.RegistryHost = options.TargetImageRegistryHost
	}

	if options.TargetImageRegistryNamespace != "" {
		targetImage.RegistryNamespace = options.TargetImageRegistryNamespace
	}

	if options.TargetImageName != "" {
		targetImage.Name = options.TargetImageName
	}

	if len(options.TargetImageTags) > 0 {
		targetImage.Version = options.TargetImageTags[0]
		promoteOptions.TargetImageTags = append(promoteOptions.TargetImageTags, options.TargetImageTags[1:]...)
	}

	if !options.EnableSemanticVersionTags {
		options.EnableSemanticVersionTags = e.configuration.EnableSemanticVersionTags
	}

	if len(options.SemanticVersionTagsTemplates) == 0 {
		options.SemanticVersionTagsTemplates = e.configuration.SemanticVersionTagsTemplates
	}

	if options.EnableSemanticVersionTags {
		semVerTags, _ := e.semver.GenerateSemverList(options.TargetImageTags, options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			promoteOptions.TargetImageTags = append(promoteOptions.TargetImageTags, semVerTags...)
		}
	}

	promoteOptions.TargetImageName, err = targetImage.DockerNormalizedNamed()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	pushAuth := e.getCredentials(targetImage.RegistryHost)
	if pushAuth != nil {
		promoteOptions.PushAuthUsername = pushAuth.Username
		promoteOptions.PushAuthPassword = pushAuth.Password
	}

	promoteOptions.RemoteSourceImage = options.RemoteSourceImage
	promoteOptions.RemoveTargetImageTags = options.RemoveTargetImageTags

	promoter, err := e.factory.GetPromoter(promoteType)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = promoter.Promote(ctx, promoteOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

func (e *Service) getCredentials(registry string) *credentials.RegistryUserPassAuth {
	auth, _ := e.credentials.GetCredentials(registry)

	return auth
}
