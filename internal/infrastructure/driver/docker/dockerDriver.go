package docker

import (
	"context"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

const (
	// DriverName is the name for the driver
	DriverName = "docker"
)

// DockerDriver is a driver for Docker
type DockerDriver struct {
	driver        DockerDriverer
	referenceName repository.ImageReferenceNamer
	writer        io.Writer
}

// NewDockerDriver creates a new DockerDriver
func NewDockerDriver(driver DockerDriverer, ref repository.ImageReferenceNamer, writer io.Writer) (*DockerDriver, error) {

	errContext := "(dockerdriver::NewDockerDriver)"

	if driver == nil {
		return nil, errors.New(errContext, "To create a DockerDriver is expected a driver")
	}

	if ref == nil {
		return nil, errors.New(errContext, "To create a DockerDriver is expected a reference name")
	}

	if writer == nil {
		writer = os.Stdout
	}

	return &DockerDriver{
		driver:        driver,
		writer:        writer,
		referenceName: ref,
	}, nil
}

// Build performs the build. In case the build could not performed it returns an error
func (d *DockerDriver) Build(ctx context.Context, i *image.Image, options *image.BuildDriverOptions) error {

	var err error
	//var dockerBuildContext dockercontext.DockerBuildContexter
	var imageName string

	errContext := "(dockerdriver::Build)"

	if d.driver == nil {
		return errors.New(errContext, "To build an image is required a driver")
	}

	if d.referenceName == nil {
		return errors.New(errContext, "To build an image is required a reference name")
	}

	if i == nil {
		return errors.New(errContext, "To build an image is required a image")
	}

	if options == nil {
		return errors.New(errContext, "To build an image is required a build options")
	}

	if ctx == nil {
		return errors.New(errContext, "To build an image is required a golang context")
	}

	if i.Name == "" {
		return errors.New(errContext, "To build an image is required an image name")
	}

	imageName, err = d.referenceName.GenerateName(i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	d.driver.WithImageName(imageName)

	if options.BuilderOptions.Dockerfile != "" {
		d.driver.WithDockerfile(options.BuilderOptions.Dockerfile)
	}

	// add docker build arguments: Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
	if len(i.PersistentVars) > 0 {
		for varName, varValue := range i.PersistentVars {
			d.driver.AddBuildArgs(varName, varValue.(string))
		}
	}

	// add docker build arguments: Vars contains the variables defined by the user on execution time and has precedences over the default values
	if len(i.Vars) > 0 {
		for varName, varValue := range i.Vars {
			d.driver.AddBuildArgs(varName, varValue.(string))
		}
	}

	// add docker tags
	if len(i.Tags) > 0 {
		for _, tag := range i.Tags {
			imageTaggedAux, err := image.NewImage(i.Name, tag, i.RegistryHost, i.RegistryNamespace)
			if err != nil {
				return errors.New(errContext, "", err)
			}
			imageTaggedName, err := d.referenceName.GenerateName(imageTaggedAux)
			if err != nil {
				return errors.New(errContext, "", err)
			}
			d.driver.AddTags(imageTaggedName)
		}
	}

	if len(i.PersistentLabels) > 0 {
		for label, value := range i.PersistentLabels {
			d.driver.AddLabel(label, value)
		}
	}

	if len(i.Labels) > 0 {
		for label, value := range i.Labels {
			d.driver.AddLabel(label, value)
		}
	}

	// add docker build arguments
	if i.Parent != nil {
		parentFullyQualifiedName, err := d.referenceName.GenerateName(i.Parent)
		if err != nil {
			return errors.New(errContext, "", err)
		}
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromFullyQualifiedNameKey], parentFullyQualifiedName)
	}

	if i.Parent != nil && i.Parent.RegistryNamespace != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], i.Parent.RegistryNamespace)
	}

	if i.Parent != nil && i.Parent.Name != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], i.Parent.Name)
	}

	if i.Parent != nil && i.Parent.Version != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], i.Parent.Version)
	}

	// add docker build arguments: map de command flag options to build argurments
	if i.Parent != nil && i.Parent.RegistryHost != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], i.Parent.RegistryHost)
		d.driver.AddAuth(options.PullAuthUsername, options.PullAuthPassword, i.Parent.RegistryHost)
	}
	d.driver.AddAuth(options.PushAuthUsername, options.PushAuthPassword, i.RegistryHost)

	if options.PushImageAfterBuild {
		d.driver.WithPushAfterBuild()
		d.driver.AddPushAuth(options.PushAuthUsername, options.PushAuthPassword)
	}

	if options.PullParentImage {
		d.driver.WithPullParentImage()
	}

	if options.RemoveImageAfterBuild {
		d.driver.WithRemoveAfterPush()
	}

	dockerBuildContextList, err := options.BuilderOptions.GetContext()
	if err != nil {
		return errors.New(errContext, "Docker building context has not been defined on build options", err)
	}

	if dockerBuildContextList == nil || len(dockerBuildContextList) == 0 {
		return errors.New(errContext, "Docker building context list is empty")
	}

	err = d.driver.AddBuildContext(dockerBuildContextList...)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	responseOutputPrefix := options.OutputPrefix
	if responseOutputPrefix == "" {
		responseOutputPrefix = imageName
	}

	d.driver.WithResponse(d.writer, responseOutputPrefix)
	d.driver.WithUseNormalizedNamed()

	err = d.driver.Run(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
