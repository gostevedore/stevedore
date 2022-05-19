package context

import (
	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	gitauth "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git/auth"
)

// DockerBuildContexter defines a docker build context entity
type DockerBuildContexter interface {
	GenerateContextFilesystem() (*filesystem.ContextFilesystem, error)
}

// GitAuthFactorier is an interface for git authentication
type GitAuthFactorier interface {
	GenerateAuthMethod(*builder.DockerDriverGitContextAuthOptions) (gitauth.GitAuther, error)
}
