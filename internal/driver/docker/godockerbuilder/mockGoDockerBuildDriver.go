package godockerbuilder

import (
	"context"
	"io"

	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
	"github.com/stretchr/testify/mock"
)

// MockGoDockerBuildDriver is a mocked docker driver
type MockGoDockerBuildDriver struct {
	mock.Mock
}

// NewMockGoDockerBuildDriver creates a new mock docker driver
func NewMockGoDockerBuildDriver() *MockGoDockerBuildDriver {
	return &MockGoDockerBuildDriver{}
}

// WithDockerfile is a mocked method
func (d *MockGoDockerBuildDriver) WithDockerfile(dockerfile string) {
	d.Mock.Called(dockerfile)
}

// WithImageName is a mocked method
func (d *MockGoDockerBuildDriver) WithImageName(image string) {
	d.Mock.Called(image)
}

// WithPullParentImage is a mocked method
func (d *MockGoDockerBuildDriver) WithPullParentImage() {
	d.Mock.Called()
}

// WithPushAfterBuild is a mocked method
func (d *MockGoDockerBuildDriver) WithPushAfterBuild() {
	d.Mock.Called()
}

// WithResponse is a mocked method
func (d *MockGoDockerBuildDriver) WithResponse(w io.Writer, prefix string) {
	d.Mock.Called(w, prefix)
}

// WithUseNormalizedNamed is a mocked method
func (d *MockGoDockerBuildDriver) WithUseNormalizedNamed() {
	d.Mock.Called()
}

// WithRemoveAfterPush is a mocked method
func (d *MockGoDockerBuildDriver) WithRemoveAfterPush() {
	d.Mock.Called()
}

// AddAuth is a mocked method
func (d *MockGoDockerBuildDriver) AddAuth(username string, password string, registry string) error {
	args := d.Mock.Called(username, password, registry)
	return args.Error(0)
}

// AddPushAuth is a mocked method
func (d *MockGoDockerBuildDriver) AddPushAuth(username string, password string) error {
	args := d.Mock.Called(username, password)
	return args.Error(0)
}

// AddBuildArgs is a mocked method
func (d *MockGoDockerBuildDriver) AddBuildArgs(arg string, value string) error {
	args := d.Mock.Called(arg, value)
	return args.Error(0)
}

// AddBuildContext is a mocked method
func (d *MockGoDockerBuildDriver) AddBuildContext(context ...*buildcontext.DockerBuildContextOptions) error {
	args := d.Mock.Called(context)
	return args.Error(0)
}

// AddLabel is a mocked method
func (d *MockGoDockerBuildDriver) AddLabel(label string, value string) error {
	args := d.Mock.Called(label, value)
	return args.Error(0)
}

// AddTags is a mocked method
func (d *MockGoDockerBuildDriver) AddTags(tags ...string) error {
	args := d.Mock.Called(tags)
	return args.Error(0)
}

// Run is a mocked method
func (d *MockGoDockerBuildDriver) Run(ctx context.Context) error {
	args := d.Mock.Called(ctx)
	return args.Error(0)
}
