package godockerbuilder

import (
	"context"
	"io"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	buildcontext "github.com/gostevedore/stevedore/internal/infrastructure/driver/docker/godockerbuilder/context"
)

// GoDockerBuildDriver is a driver for building docker images
type GoDockerBuildDriver struct {
	cmd            DockerBuilder
	contextFactory *buildcontext.DockerBuildContextFactory

	addBuildArgsMutex sync.Mutex
	addLabelMutex     sync.Mutex
	addTagsMutex      sync.Mutex
}

// NewGoDockerBuildDriver creates a new GoDockerBuildDriver
func NewGoDockerBuildDriver(cmd *build.DockerBuildCmd, contextFactory *buildcontext.DockerBuildContextFactory) *GoDockerBuildDriver {
	return &GoDockerBuildDriver{
		cmd:            cmd,
		contextFactory: contextFactory,
	}
}

// WithDockerfile sets dockerfile to use
func (d *GoDockerBuildDriver) WithDockerfile(dockerfile string) {
	d.cmd = d.cmd.WithDockerfile(dockerfile)
}

// WithImageName sets the image name
func (d *GoDockerBuildDriver) WithImageName(image string) {
	d.cmd = d.cmd.WithImageName(image)
}

// WithPullParentImage sets if the image should be pushed after build
func (d *GoDockerBuildDriver) WithPullParentImage() {
	d.cmd = d.cmd.WithPullParentImage()
}

// WithPushAfterBuild sets if the image should be pushed after build
func (d *GoDockerBuildDriver) WithPushAfterBuild() {
	d.cmd = d.cmd.WithPushAfterBuild()
}

// WithUseNormalizedNamed sets if image name should be normalized
func (d *GoDockerBuildDriver) WithUseNormalizedNamed() {
	d.cmd = d.cmd.WithUseNormalizedNamed()
}

// WithRemoveAfterPush sets if the image should be removed after push
func (d *GoDockerBuildDriver) WithRemoveAfterPush() {
	d.cmd = d.cmd.WithRemoveAfterPush()
}

// WithResponse sets the responser to use
func (d *GoDockerBuildDriver) WithResponse(w io.Writer, prefix string) {

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend(prefix),
		),
		response.WithWriter(w),
	)

	d.cmd = d.cmd.WithResponse(res)
}

// AddAuth defines the authentication to use for an specific registry
func (d *GoDockerBuildDriver) AddAuth(username string, password string, registry string) error {
	return d.cmd.AddAuth(username, password, registry)
}

// AddPushAuth defines the authentication to use for an specific registry
func (d *GoDockerBuildDriver) AddPushAuth(username string, password string) error {
	return d.cmd.AddPushAuth(username, password)
}

// AddBuildArgs append new build args
func (d *GoDockerBuildDriver) AddBuildArgs(arg string, value string) error {
	d.addBuildArgsMutex.Lock()
	defer d.addBuildArgsMutex.Unlock()

	return d.cmd.AddBuildArgs(arg, value)
}

// AddBuildContext sets those docker build contexts required to build an image. It supports to use several context which are merged before to start the image build
func (d *GoDockerBuildDriver) AddBuildContext(options ...*builder.DockerDriverContextOptions) error {

	errContext := "(godockerbuilder::AddBuildContext)"

	if options == nil {
		return errors.New(errContext, "Docker build context options are missing")
	}

	if len(options) == 0 {
		return errors.New(errContext, "No Docker build context is defined")
	}

	buildContextList := []godockerbuilderbuildcontext.DockerBuildContexter{}

	for _, context := range options {
		c, err := d.contextFactory.GenerateDockerBuildContext(context)
		if err != nil {
			return err
		}
		buildContextList = append(buildContextList, c)
	}

	return d.cmd.AddBuildContext(buildContextList...)
}

// AddLabel adds a label to the image
func (d *GoDockerBuildDriver) AddLabel(label string, value string) error {
	d.addLabelMutex.Lock()
	defer d.addLabelMutex.Unlock()

	return d.cmd.AddLabel(label, value)
}

// AddTags adds tags to the image
func (d *GoDockerBuildDriver) AddTags(tags ...string) error {
	d.addTagsMutex.Lock()
	defer d.addTagsMutex.Unlock()

	return d.cmd.AddTags(tags...)
}

// Run starts the build
func (d *GoDockerBuildDriver) Run(ctx context.Context) error {
	return d.cmd.Run(ctx)
}
