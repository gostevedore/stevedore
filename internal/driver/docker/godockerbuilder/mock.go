package godockerbuilder

import (
	"context"

	"github.com/apenella/go-docker-builder/pkg/types"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
	"github.com/stretchr/testify/mock"
)

// MockDockerDriver is a mocked docker driver
type MockDockerDriver struct {
	mock.Mock
}

// NewMockDockerDriver creates a new mock docker driver
func NewMockDockerDriver() *MockDockerDriver {
	return &MockDockerDriver{}
}

// WithDockerfile is a mocked method
func (d *MockDockerDriver) WithDockerfile(dockerfile string) {
	d.Mock.Called(dockerfile)
}

// WithImageName is a mocked method
func (d *MockDockerDriver) WithImageName(image string) {
	d.Mock.Called(image)
}

// WithPushAfterBuild is a mocked method
func (d *MockDockerDriver) WithPushAfterBuild() {
	d.Mock.Called()
}

// WithResponse is a mocked method
func (d *MockDockerDriver) WithResponse(response types.Responser) {
	d.Mock.Called(response)
}

// WithUseNormalizedNamed is a mocked method
func (d *MockDockerDriver) WithUseNormalizedNamed() {
	d.Mock.Called()
}

// WithRemoveAfterPush is a mocked method
func (d *MockDockerDriver) WithRemoveAfterPush() {
	d.Mock.Called()
}

// AddAuth is a mocked method
func (d *MockDockerDriver) AddAuth(username string, password string, registry string) error {
	args := d.Mock.Called(username, password, registry)
	return args.Error(0)
}

// AddPushAuth is a mocked method
func (d *MockDockerDriver) AddPushAuth(username string, password string) error {
	args := d.Mock.Called(username, password)
	return args.Error(0)
}

// AddBuildArgs is a mocked method
func (d *MockDockerDriver) AddBuildArgs(arg string, value string) error {
	args := d.Mock.Called(arg, value)
	return args.Error(0)
}

// AddBuildContext is a mocked method
func (d *MockDockerDriver) AddBuildContext(context ...*buildcontext.DockerBuildContextOptions) error {
	args := d.Mock.Called(context)
	return args.Error(0)
}

// AddLabel is a mocked method
func (d *MockDockerDriver) AddLabel(label string, value string) error {
	args := d.Mock.Called(label, value)
	return args.Error(0)
}

// AddTags is a mocked method
func (d *MockDockerDriver) AddTags(tags ...string) error {
	args := d.Mock.Called(tags)
	return args.Error(0)
}

// Run is a mocked method
func (d *MockDockerDriver) Run(ctx context.Context) error {
	args := d.Mock.Called(ctx)
	return args.Error(0)
}
