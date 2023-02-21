package godockerbuilder

import (
	"context"

	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/types"
)

// DockerBuilder defines a docker driver
type DockerBuilder interface {
	WithDockerfile(string) *build.DockerBuildCmd
	WithImageName(string) *build.DockerBuildCmd
	WithPullParentImage() *build.DockerBuildCmd
	WithPushAfterBuild() *build.DockerBuildCmd
	WithResponse(types.Responser) *build.DockerBuildCmd
	WithUseNormalizedNamed() *build.DockerBuildCmd
	WithRemoveAfterPush() *build.DockerBuildCmd
	AddAuth(string, string, string) error
	AddPushAuth(string, string) error
	AddBuildArgs(string, string) error
	AddBuildContext(...godockerbuilderbuildcontext.DockerBuildContexter) error
	AddLabel(string, string) error
	AddTags(...string) error
	Run(context.Context) error
}
