package godockerbuilder

import (
	"context"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/response"
	contextoptions "github.com/gostevedore/stevedore/internal/driver/docker/context"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context"
)

// GoDockerBuildDriver is a driver for building docker images
type GoDockerBuildDriver struct {
	docker         DockerBuilder
	contextFactory *buildcontext.DockerBuildContextFactory
}

// NewGoDockerBuildDriver creates a new GoDockerBuildDriver
func NewGoDockerBuildDriver(contextFactory *buildcontext.DockerBuildContextFactory) *GoDockerBuildDriver {
	return &GoDockerBuildDriver{
		docker:         &build.DockerBuildCmd{},
		contextFactory: contextFactory,
	}
}

// WithDockerfile sets dockerfile to use
func (d *GoDockerBuildDriver) WithDockerfile(dockerfile string) {
	d.docker = d.docker.WithDockerfile(dockerfile)
}

// WithImageName sets the image name
func (d *GoDockerBuildDriver) WithImageName(image string) {
	d.docker = d.docker.WithImageName(image)
}

// WithPullParentImage sets if the image should be pushed after build
func (d *GoDockerBuildDriver) WithPullParentImage() {
	d.docker = d.docker.WithPullParentImage()
}

// WithPushAfterBuild sets if the image should be pushed after build
func (d *GoDockerBuildDriver) WithPushAfterBuild() {
	d.docker = d.docker.WithPushAfterBuild()
}

// WithUseNormalizedNamed sets if image name should be normalized
func (d *GoDockerBuildDriver) WithUseNormalizedNamed() {
	d.docker = d.docker.WithUseNormalizedNamed()
}

// WithRemoveAfterPush sets if the image should be removed after push
func (d *GoDockerBuildDriver) WithRemoveAfterPush() {
	d.docker = d.docker.WithRemoveAfterPush()
}

// WithResponse sets the responser to use
func (d *GoDockerBuildDriver) WithResponse(w io.Writer, prefix string) {

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend(prefix),
		),
		response.WithWriter(w),
	)

	d.docker = d.docker.WithResponse(res)
}

// AddAuth defines the authentication to use for an specific registry
func (d *GoDockerBuildDriver) AddAuth(username string, password string, registry string) error {
	return d.docker.AddAuth(username, password, registry)
}

// AddPushAuth defines the authentication to use for an specific registry
func (d *GoDockerBuildDriver) AddPushAuth(username string, password string) error {
	return d.docker.AddPushAuth(username, password)
}

// AddBuildArgs append new build args
func (d *GoDockerBuildDriver) AddBuildArgs(arg string, value string) error {
	return d.docker.AddBuildArgs(arg, value)
}

// AddBuildContext sets those docker build contexts required to build an image. It supports to use several context which are merged before to start the image build
func (d *GoDockerBuildDriver) AddBuildContext(options ...*contextoptions.DockerBuildContextOptions) error {

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
func (d *GoDockerBuildDriver) AddLabel(label string, value string) error {
	return d.docker.AddLabel(label, value)
}

// AddTags adds tags to the image
func (d *GoDockerBuildDriver) AddTags(tags ...string) error {
	return d.docker.AddTags(tags...)
}

// Run starts the build
func (d *GoDockerBuildDriver) Run(ctx context.Context) error {
	return d.docker.Run(ctx)
}
