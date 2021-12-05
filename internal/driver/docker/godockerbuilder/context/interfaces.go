package context

import (
	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
	gitauth "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git/auth"
)

// DockerBuildContexter defines a docker build context entity
type DockerBuildContexter interface {
	GenerateContextFilesystem() (*filesystem.ContextFilesystem, error)
}

// GitAuthFactorier is an interface for git authentication
type GitAuthFactorier interface {
	GenerateAuthMethod(*buildcontext.GitContextAuthOptions) (gitauth.GitAuther, error)
}
