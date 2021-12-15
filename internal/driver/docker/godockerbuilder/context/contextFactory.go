package context

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	gitcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git"
	pathcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/path"
)

// DockerBuildContextFactory is a factory for docker build context
type DockerBuildContextFactory struct {
	gitAuth GitAuthFactorier
}

// NewDockerBuildContextFactory returns a new DockerBuildContextFactory
func NewDockerBuildContextFactory(gitAuth GitAuthFactorier) *DockerBuildContextFactory {
	return &DockerBuildContextFactory{
		gitAuth: gitAuth,
	}
}

// GenerateDockerBuildContext returns the docker build context
func (f *DockerBuildContextFactory) GenerateDockerBuildContext(options *builder.DockerDriverContextOptions) (DockerBuildContexter, error) {

	errContext := "(DockerBuildContextFactory::GenerateDockerBuildContext)"

	if options == nil {
		return nil, errors.New(errContext, "Docker build context options are required to generate a build context")
	}

	// when path is defined, is returned a path context
	if options.Path != "" {
		pathBuildContext := pathcontext.NewPathBuildContext()
		pathBuildContext.WithPath(options.Path)
		return pathBuildContext, nil
	}

	// when git is defined, is returned a git context
	if options.Git != nil {
		if options.Git.Repository == "" {
			return nil, errors.New(errContext, "A repository must be specified on git build docker context")
		}

		gitBuildContext := gitcontext.NewGitBuildContext()
		gitBuildContext.WithRepository(options.Git.Repository)

		if options.Git.Reference != "" {
			gitBuildContext.WithReference(options.Git.Reference)
		}

		if options.Git.Path != "" {
			gitBuildContext.WithPath(options.Git.Path)
		}

		if options.Git.Auth != nil {

			if f.gitAuth == nil {
				return nil, errors.New(errContext, "Git auth generator is required to generate a git build context")
			}

			auth, err := f.gitAuth.GenerateAuthMethod(options.Git.Auth)
			if err != nil {
				return nil, errors.New(errContext, err.Error())
			}

			gitBuildContext.WithAuth(auth)
		}

		return gitBuildContext, nil
	}

	return nil, errors.New(errContext, "Unknown context type")
}
