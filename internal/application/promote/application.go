package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// OptionsFunc is a function used to configure the application
type OptionsFunc func(*Application)

// Application is the application used to promote images
type Application struct {
	credentials CredentialsStorer
	factory     PromoteFactorier
	semver      Semverser
}

// NewApplication creates an application
func NewApplication(options ...OptionsFunc) *Application {
	app := &Application{}
	app.Options(options...)

	return app
}

// WitCredentials sets credentials for the service
func WithCredentials(c CredentialsStorer) OptionsFunc {
	return func(a *Application) {
		a.credentials = c
	}
}

// WithPromoteFactory sets the factory used to create the promoter
func WithPromoteFactory(f PromoteFactorier) OptionsFunc {
	return func(a *Application) {
		a.factory = f
	}
}

// WithSemver sets the semver component for the service
func WithSemver(sv Semverser) OptionsFunc {
	return func(a *Application) {
		a.semver = sv
	}
}

// Options configure the service
func (a *Application) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Promote an image
func (a *Application) Promote(ctx context.Context, options *Options) error {

	var err error
	var sourceImage, targetImage *image.Image

	promoteOptions := &image.PromoteOptions{}
	errContext := "(Service::Promote)"

	if a.factory == nil {
		return errors.New(errContext, "Promote factory has not been initialized")
	}

	if a.semver == nil {
		return errors.New(errContext, "Semver has not been initialized")
	}

	if a.credentials == nil {
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
		return errors.New(errContext, "", err)
	}
	targetImage, err = sourceImage.Copy()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	pullAuth, err := a.getCredentials(sourceImage.RegistryHost)
	if err != nil {
		return errors.New(errContext, "", err)
	}
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

	if options.EnableSemanticVersionTags {
		semVerTags, _ := a.semver.GenerateSemverList(options.TargetImageTags, options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			promoteOptions.TargetImageTags = append(promoteOptions.TargetImageTags, semVerTags...)
		}
	}

	promoteOptions.TargetImageName, err = targetImage.DockerNormalizedNamed()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	pushAuth, err := a.getCredentials(targetImage.RegistryHost)
	if err != nil {
		return errors.New(errContext, "", err)
	}
	if pushAuth != nil {
		promoteOptions.PushAuthUsername = pushAuth.Username
		promoteOptions.PushAuthPassword = pushAuth.Password
	}

	promoteOptions.RemoteSourceImage = options.RemoteSourceImage
	promoteOptions.RemoveTargetImageTags = options.RemoveTargetImageTags

	promoter, err := a.getPromoter(options)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = promoter.Promote(ctx, promoteOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (a *Application) getCredentials(registry string) (*credentials.UserPasswordAuth, error) {
	errContext := "(Service::getCredentials)"

	if a.credentials == nil {
		return nil, errors.New(errContext, "Credentials has not been initialized")
	}

	auth, _ := a.credentials.Get(registry)

	return auth, nil
}

func (a *Application) getPromoter(options *Options) (repository.Promoter, error) {

	errContext := "(Handler::getPromoter)"

	if a.factory == nil {
		return nil, errors.New(errContext, "Promote factory has not been initialized")
	}

	promoteDriver := "docker"
	if options.DryRun {
		promoteDriver = "dry-run"
	}
	promoter, err := a.factory.Get(promoteDriver)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return promoter, nil
}
