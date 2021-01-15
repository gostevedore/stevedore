package context

import (
	gitcontext "github.com/gostevedore/stevedore/internal/driver/docker/context/git"

	errors "github.com/apenella/go-common-utils/error"
	dockercontext "github.com/apenella/go-docker-builder/pkg/build/context"
	dockercontextgit "github.com/apenella/go-docker-builder/pkg/build/context/git"
	dockercontextpath "github.com/apenella/go-docker-builder/pkg/build/context/path"
	"gopkg.in/yaml.v2"
)

const (
	PathContextType uint8 = iota
	GitContextType
	UnknownContextType
)

// DockerBuildContext
type DockerBuildContext struct {
	Path string                 `yaml:"path"`
	Git  *gitcontext.GitContext `yaml:"git"`
}

func GenerateDockerBuildContext(context interface{}) (dockercontext.DockerBuildContexter, error) {

	var buildContext dockercontext.DockerBuildContexter
	dockerBuildContext := &DockerBuildContext{}

	byteContext, err := yaml.Marshal(context)
	if err != nil {
		return nil, errors.New("(build::docker::context::GenerateDockerBuildContext)", "Docker build context could not be marshalled", err)
	}

	err = yaml.Unmarshal(byteContext, dockerBuildContext)
	if err != nil {
		return nil, errors.New("(build::docker::context::GenerateDockerBuildContext)", "Docker build context could not be unmarshalled to DockerBuildContext", err)
	}

	switch dockerBuildContext.GetContextType() {
	case PathContextType:

		buildContext = &dockercontextpath.PathBuildContext{
			Path: dockerBuildContext.Path,
		}

	case GitContextType:

		if dockerBuildContext.Git.Repository == "" {
			return nil, errors.New("(build::docker::context::GenerateDockerBuildContext)", "A repository must be specified on git build docker context")
		}

		buildContext = &dockercontextgit.GitBuildContext{
			Repository: dockerBuildContext.Git.Repository,
		}
		if dockerBuildContext.Git.Reference != "" {
			buildContext.(*dockercontextgit.GitBuildContext).Reference = dockerBuildContext.Git.Reference
		}

	default:
		return nil, errors.New("(build::docker::context::GenerateDockerBuildContext)", "Unknown context type")
	}

	return buildContext, nil
}

func (c *DockerBuildContext) GetContextType() uint8 {

	if c.IsPathContext() {
		return PathContextType
	}

	if c.IsGitContext() {
		return GitContextType
	}

	return UnknownContextType
}

func (c *DockerBuildContext) IsPathContext() bool {
	return c.Path != ""
}

func (c *DockerBuildContext) IsGitContext() bool {
	return c.Git != nil
}
