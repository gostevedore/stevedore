package godockerbuilder

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/types"
	contextoptions "github.com/gostevedore/stevedore/internal/driver/docker/context"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context"
)

// GoDockerDriver is a driver for building docker images
type GoDockerDriver struct {
	docker         DockerBuilder
	contextFactory *buildcontext.DockerBuildContextFactory
}

// NewGoDockerDriver creates a new GoDockerDriver
func NewGoDockerDriver(contextFactory *buildcontext.DockerBuildContextFactory) *GoDockerDriver {
	return &GoDockerDriver{
		docker:         &build.DockerBuildCmd{},
		contextFactory: contextFactory,
	}
}

// WithDockerfile sets dockerfile to use
func (d *GoDockerDriver) WithDockerfile(dockerfile string) {
	d.docker = d.docker.WithDockerfile(dockerfile)
}

// WithImageName sets the image name
func (d *GoDockerDriver) WithImageName(image string) {
	d.docker = d.docker.WithImageName(image)
}

// WitPushAfterBuild sets if the image should be pushed after build
func (d *GoDockerDriver) WithPushAfterBuild() {
	d.docker = d.docker.WithPushAfterBuild()
}

// WithResponse sets responser to manage the response
func (d *GoDockerDriver) WithResponse(response types.Responser) {
	d.docker = d.docker.WithResponse(response)
}

// WithUseNormalizedNamed sets if image name should be normalized
func (d *GoDockerDriver) WithUseNormalizedNamed() {
	d.docker = d.docker.WithUseNormalizedNamed()
}

// WithRemoveAfterPush sets if the image should be removed after push
func (d *GoDockerDriver) WithRemoveAfterPush() {
	d.docker = d.docker.WithRemoveAfterPush()
}

// AddAuth defines the authentication to use for an specific registry
func (d *GoDockerDriver) AddAuth(username string, password string, registry string) error {
	return d.docker.AddAuth(username, password, registry)
}

// AddPushAuth defines the authentication to use for an specific registry
func (d *GoDockerDriver) AddPushAuth(username string, password string) error {
	return d.docker.AddPushAuth(username, password)
}

// AddBuildArgs append new build args
func (d *GoDockerDriver) AddBuildArgs(arg string, value string) error {
	return d.docker.AddBuildArgs(arg, value)
}

// AddBuildContext sets those docker build contexts required to build an image. It supports to use several context which are merged before to start the image build
func (d *GoDockerDriver) AddBuildContext(options ...*contextoptions.DockerBuildContextOptions) error {

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

	return d.docker.AddBuildContext(buildContextList...)
}

// AddLabel adds a label to the image
func (d *GoDockerDriver) AddLabel(label string, value string) error {
	return d.docker.AddLabel(label, value)
}

// AddTags adds tags to the image
func (d *GoDockerDriver) AddTags(tags ...string) error {
	return d.docker.AddTags(tags...)
}

// Run starts the build
func (d *GoDockerDriver) Run(ctx context.Context) error {
	return d.docker.Run(ctx)
}
