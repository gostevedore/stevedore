package godockerbuilder

import (
	"context"

	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/types"
)

// DockerBuilder defines a docker driver
type DockerBuilder interface {
	WithDockerfile(dockerfile string) *build.DockerBuildCmd
	WithImageName(image string) *build.DockerBuildCmd
	WithPushAfterBuild() *build.DockerBuildCmd
	WithResponse(response types.Responser) *build.DockerBuildCmd
	WithUseNormalizedNamed() *build.DockerBuildCmd
	WithRemoveAfterPush() *build.DockerBuildCmd
	AddAuth(username string, password string, registry string) error
	AddPushAuth(username string, password string) error
	AddBuildArgs(arg string, value string) error
	AddBuildContext(contexts ...godockerbuilderbuildcontext.DockerBuildContexter) error
	AddLabel(label string, value string) error
	AddTags(tags ...string) error
	Run(ctx context.Context) error
}
