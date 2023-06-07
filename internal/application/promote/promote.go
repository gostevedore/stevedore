package promote

import (
	"context"
	"fmt"
	"sort"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	authmethodbasic "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/types/list"
)

// OptionsFunc is a function used to configure the application
type OptionsFunc func(*Application)

// Application is the application used to promote images
type Application struct {
	credentials    repository.AuthFactorier
	factory        PromoteFactorier
	referenceNamer repository.ImageReferenceNamer
	semver         Semverser
}

// NewApplication creates an application
func NewApplication(options ...OptionsFunc) *Application {
	app := &Application{}
	app.Options(options...)

	return app
}

// WitCredentials sets credentials for the application
func WithCredentials(c repository.AuthFactorier) OptionsFunc {
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

// WithSemver sets the semver component for the application
func WithSemver(sv Semverser) OptionsFunc {
	return func(a *Application) {
		a.semver = sv
	}
}

// WithReferenceNamer sets the reference namer component for the application
func WithReferenceNamer(ref repository.ImageReferenceNamer) OptionsFunc {
	return func(a *Application) {
		a.referenceNamer = ref
	}
}

// Options configure the application
func (a *Application) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(a)
	}
}

// Promote method carries out the application tasks
func (a *Application) Promote(ctx context.Context, options *Options) error {

	var err error
	var sourceImage, targetImage *image.Image
	var referenceName string

	promoteOptions := &image.PromoteOptions{}
	errContext := "(application::promote::Promote)"

	if a.factory == nil {
		return errors.New(errContext, "Promote application requires promote factory")
	}

	if a.semver == nil {
		return errors.New(errContext, "Promote application requires semver")
	}

	if a.referenceNamer == nil {
		return errors.New(errContext, "Promote application requires a image reference namer")
	}

	if a.credentials == nil {
		return errors.New(errContext, "Promote application requires credentials factory")
	}

	if options == nil {
		return errors.New(errContext, "Promote application requires options")
	}

	if options.SourceImageName == "" {
		return errors.New(errContext, "Promote application options requires an image source name defined")
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

	auth, err := a.getCredentials(sourceImage.RegistryHost)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if auth != nil {
		pullAuth, isBasicAuth := auth.(*authmethodbasic.BasicAuthMethod)
		if !isBasicAuth {
			return errors.New(errContext, fmt.Sprintf("Invalid credentials method for '%s'. Found '%s' when is expected basic auth method", sourceImage.RegistryHost, auth.Name()))
		}

		promoteOptions.PullAuthUsername = pullAuth.Username
		promoteOptions.PullAuthPassword = pullAuth.Password
	}

	if options.TargetImageRegistryHost != image.UndefinedStringValue {
		targetImage.RegistryHost = options.TargetImageRegistryHost
	}

	if options.TargetImageRegistryNamespace != image.UndefinedStringValue {
		targetImage.RegistryNamespace = options.TargetImageRegistryNamespace
	}

	if options.TargetImageName != image.UndefinedStringValue {
		targetImage.Name = options.TargetImageName
	}

	auxTargetImageTagsMap := map[string]struct{}{}

	if options.PromoteSourceImageTag {
		auxTargetImageTagsMap[sourceImage.Version] = struct{}{}
	}

	if len(options.TargetImageTags) > 0 {
		targetImage.Version = options.TargetImageTags[0]

		for _, tag := range options.TargetImageTags[1:] {
			auxTargetImageTagsMap[tag] = struct{}{}
		}
	}

	if options.EnableSemanticVersionTags {
		auxTargetImageTagsMap[targetImage.Version] = struct{}{}
		auxTargetImageSemVerList := []string{}
		for tag := range auxTargetImageTagsMap {
			auxTargetImageSemVerList = append(auxTargetImageSemVerList, tag)
		}

		semVerTags, _ := a.semver.GenerateSemverList(auxTargetImageSemVerList, options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			for _, tag := range semVerTags {
				auxTargetImageTagsMap[tag] = struct{}{}
			}
		}
	}

	referenceName, err = a.referenceNamer.GenerateName(targetImage)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error generating target image reference name for '%s'", promoteOptions.SourceImageName), err)
	}
	promoteOptions.TargetImageName = referenceName

	auxTargetImageTagsList := []string{}
	for tag := range auxTargetImageTagsMap {
		if tag != targetImage.Version {
			auxTargetImageTagsList = append(auxTargetImageTagsList, tag)
		}
	}
	promoteOptions.TargetImageTags, err = a.generateReferenceNameList(targetImage, auxTargetImageTagsList)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	sort.Sort(list.SortedStringList(promoteOptions.TargetImageTags))

	// Registry host must be defined explicitly to achive the host credentials
	if targetImage.RegistryHost != "" {
		auth, err = a.getCredentials(targetImage.RegistryHost)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if auth != nil {
			pushAuth, isBasicAuth := auth.(*authmethodbasic.BasicAuthMethod)
			if !isBasicAuth {
				return errors.New(errContext, fmt.Sprintf("Invalid credentials method for '%s'. Found '%s' when is expected basic auth method", targetImage.RegistryHost, auth.Name()))
			}

			promoteOptions.PushAuthUsername = pushAuth.Username
			promoteOptions.PushAuthPassword = pushAuth.Password
		}
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

// generateReferenceNameList return a list of reference names
func (a *Application) generateReferenceNameList(i *image.Image, tags []string) ([]string, error) {

	errContext := "(application::promote::generateReferenceNameList)"
	list := []string{}

	for _, tag := range tags {
		if i.Version != tag {
			var auxReferenceName string
			var err error
			var auxImage *image.Image

			auxImage, err = i.Copy()
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}
			auxImage.Version = tag
			auxReferenceName, err = a.referenceNamer.GenerateName(auxImage)
			list = append(list, auxReferenceName)
		}
	}

	return list, nil
}

func (a *Application) getCredentials(registry string) (repository.AuthMethodReader, error) {
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
