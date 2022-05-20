package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/promote"
)

// OptionsFunc is a function used to configure the service
type OptionsFunc func(*Service)

// Service is the service used to promote images
type Service struct {
	credentials CredentialsStorer
	factory     PromoteFactorier
	semver      Semverser
}

// NewService creates a new service
func NewService(options ...OptionsFunc) *Service {
	service := &Service{}
	service.Options(options...)

	return service
}

// WitCredentials sets credentials for the service
func WithCredentials(c CredentialsStorer) OptionsFunc {
	return func(s *Service) {
		s.credentials = c
	}
}

// WithPromoteFactory sets the factory used to create the promoter
func WithPromoteFactory(f PromoteFactorier) OptionsFunc {
	return func(s *Service) {
		s.factory = f
	}
}

// WithSemver sets the semver component for the service
func WithSemver(sv Semverser) OptionsFunc {
	return func(s *Service) {
		s.semver = sv
	}
}

// Options configure the service
func (s *Service) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

// Promote an image
func (s *Service) Promote(ctx context.Context, options *ServiceOptions) error {

	var err error
	var sourceImage, targetImage *image.Image

	promoteOptions := &promote.PromoteOptions{}
	errContext := "(Service::Promote)"

	if s.factory == nil {
		return errors.New(errContext, "Promote factory has not been initialized")
	}

	if s.semver == nil {
		return errors.New(errContext, "Semver has not been initialized")
	}

	if s.credentials == nil {
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

	pullAuth, err := s.getCredentials(sourceImage.RegistryHost)
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
		semVerTags, _ := s.semver.GenerateSemverList(options.TargetImageTags, options.SemanticVersionTagsTemplates)
		if len(semVerTags) > 0 {
			promoteOptions.TargetImageTags = append(promoteOptions.TargetImageTags, semVerTags...)
		}
	}

	promoteOptions.TargetImageName, err = targetImage.DockerNormalizedNamed()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	pushAuth, err := s.getCredentials(targetImage.RegistryHost)
	if err != nil {
		return errors.New(errContext, "", err)
	}
	if pushAuth != nil {
		promoteOptions.PushAuthUsername = pushAuth.Username
		promoteOptions.PushAuthPassword = pushAuth.Password
	}

	promoteOptions.RemoteSourceImage = options.RemoteSourceImage
	promoteOptions.RemoveTargetImageTags = options.RemoveTargetImageTags

	promoter, err := s.getPromoter(options)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = promoter.Promote(ctx, promoteOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (e *Service) getCredentials(registry string) (*credentials.UserPasswordAuth, error) {
	errContext := "(Service::getCredentials)"

	if e.credentials == nil {
		return nil, errors.New(errContext, "Credentials has not been initialized")
	}

	auth, _ := e.credentials.Get(registry)

	return auth, nil
}

func (e *Service) getPromoter(options *ServiceOptions) (promote.Promoter, error) {

	errContext := "(Handler::getPromoter)"

	if e.factory == nil {
		return nil, errors.New(errContext, "Promote factory has not been initialized")
	}

	promoteDriver := "docker"
	if options.DryRun {
		promoteDriver = "dry-run"
	}
	promoter, err := e.factory.Get(promoteDriver)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return promoter, nil
}
