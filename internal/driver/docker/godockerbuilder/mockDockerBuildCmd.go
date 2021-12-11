package godockerbuilder

import (
	"context"

	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/types"
	"github.com/stretchr/testify/mock"
)

// MockDockerBuildCmd is a mock for DockerBuildCmd "github.com/apenella/go-docker-builder/pkg/build"
type MockDockerBuildCmd struct {
	mock.Mock
}

// WithDockerfile is a mock for WithDockerfile "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithDockerfile(dockerfile string) *build.DockerBuildCmd {
	d.Called(dockerfile)
	return &build.DockerBuildCmd{}
}

// WithImageName is a mock for WithImageName "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithImageName(image string) *build.DockerBuildCmd {
	d.Called(image)
	return &build.DockerBuildCmd{}
}

// WithPullParentImage is a mock for WithPullParentImage "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithPullParentImage() *build.DockerBuildCmd {
	d.Called()
	return &build.DockerBuildCmd{}
}

// WithPushAfterBuild is a mock for WithPushAfterBuild "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithPushAfterBuild() *build.DockerBuildCmd {
	d.Called()
	return &build.DockerBuildCmd{}
}

// WithResponse is a mock for WithResponse "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithResponse(response types.Responser) *build.DockerBuildCmd {
	d.Called(response)
	return &build.DockerBuildCmd{}
}

// WithUseNormalizedNamed is a mock for WithUseNormalizedNamed "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithUseNormalizedNamed() *build.DockerBuildCmd {
	d.Called()
	return &build.DockerBuildCmd{}
}

// WithRemoveAfterPush is a mock for WithRemoveAfterPush "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) WithRemoveAfterPush() *build.DockerBuildCmd {
	d.Called()
	return &build.DockerBuildCmd{}
}

// AddAuth is a mock for AddAuth "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) AddAuth(username string, password string, registry string) error {
	args := d.Called(username, password, registry)
	return args.Error(0)
}

// AddPushAuth is a mock for AddPushAuth "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) AddPushAuth(username string, password string) error {
	args := d.Called(username, password)
	return args.Error(0)
}

// AddBuildArgs is a mock for AddBuildArgs "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) AddBuildArgs(arg string, value string) error {
	args := d.Called(arg, value)
	return args.Error(0)
}

// AddBuildContext is a mock for AddBuildContext "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) AddBuildContext(contexts ...godockerbuilderbuildcontext.DockerBuildContexter) error {
	args := d.Called(contexts)
	return args.Error(0)
}

// AddLabel is a mock for AddLabel "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) AddLabel(label string, value string) error {
	args := d.Called(label, value)
	return args.Error(0)
}

// AddTags is a mock for AddTags "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) AddTags(tags ...string) error {
	args := d.Called(tags)
	return args.Error(0)
}

// Run is a mock for Run "github.com/apenella/go-docker-builder/pkg/build"
func (d *MockDockerBuildCmd) Run(ctx context.Context) error {
	args := d.Called(ctx)
	return args.Error(0)
}
