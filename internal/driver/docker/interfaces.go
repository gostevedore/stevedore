package dockerdriver

import (
	"context"

	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/apenella/go-docker-builder/pkg/types"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
)

// DockerBuildContexter defines a docker build context
type DockerBuildContexter interface {
	GenerateContextFilesystem() (*filesystem.ContextFilesystem, error)
}

// DockerBuildContextFactorier defines a docker build context factory
type DockerBuildContextFactorier interface {
	GenerateDockerBuildContext(context interface{}) (DockerBuildContexter, error)
}

// DockerDriverer defines a docker driver
type DockerDriverer interface {
	WithDockerfile(dockerfile string)
	WithImageName(image string)
	WithPushAfterBuild()
	WithResponse(response types.Responser)
	WithUseNormalizedNamed()
	WithRemoveAfterPush()
	AddAuth(username string, password string, registry string) error
	AddPushAuth(username string, password string) error
	AddBuildArgs(arg string, value string) error
	AddBuildContext(context ...*buildcontext.DockerBuildContextOptions) error
	AddLabel(label string, value string) error
	AddTags(tags ...string) error
	Run(ctx context.Context) error
}
