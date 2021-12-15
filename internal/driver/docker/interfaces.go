package dockerdriver

import (
	"context"
	"io"

	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/gostevedore/stevedore/internal/builders/builder"
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
	WithDockerfile(string)
	WithImageName(string)
	WithPullParentImage()
	WithPushAfterBuild()
	WithResponse(io.Writer, string)
	WithUseNormalizedNamed()
	WithRemoveAfterPush()
	AddAuth(string, string, string) error
	AddPushAuth(string, string) error
	AddBuildArgs(string, string) error
	AddBuildContext(...*builder.DockerDriverContextOptions) error
	AddLabel(string, string) error
	AddTags(...string) error
	Run(context.Context) error
}
