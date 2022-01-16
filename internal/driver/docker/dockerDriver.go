package dockerdriver

import (
	"context"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/varsmap"
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/images/image"
)

const (
	// DriverName is the name for the driver
	DriverName = "docker"
)

// DockerDriver is a driver for Docker
type DockerDriver struct {
	driver DockerDriverer
	writer io.Writer
}

// NewDockerDriver creates a new DockerDriver
func NewDockerDriver(driver DockerDriverer, writer io.Writer) (*DockerDriver, error) {

	errContext := "(dockerdriver::NewDockerDriver)"

	if driver == nil {
		return nil, errors.New(errContext, "To create a DockerDriver is expected a driver")
	}

	if writer == nil {
		writer = os.Stdout
	}

	return &DockerDriver{
		driver: driver,
		writer: writer,
	}, nil
}

// Build performs the build. In case the build could not performed it returns an error
func (d *DockerDriver) Build(ctx context.Context, options *driver.BuildDriverOptions) error {

	var err error
	//var dockerBuildContext dockercontext.DockerBuildContexter
	var imageName string

	errContext := "(dockerdriver::Build)"

	if d.driver == nil {
		return errors.New(errContext, "To build an image is required a driver")
	}

	if options == nil {
		return errors.New(errContext, "To build an image is required a build options")
	}

	if ctx == nil {
		return errors.New(errContext, "To build an image is required a golang context")
	}

	if options.ImageName == "" {
		return errors.New(errContext, "To build an image is required an image name")
	}

	imageAux, err := image.NewImage(options.ImageName, options.ImageVersion, options.RegistryHost, options.RegistryNamespace)
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	imageName, err = imageAux.DockerNormalizedNamed()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	d.driver.WithImageName(imageName)

	// TO REMOVE
	// add docker build context to cmd instance
	// builderConfOptions := options.BuilderOptions

	if options.BuilderOptions.Dockerfile != "" {
		d.driver.WithDockerfile(options.BuilderOptions.Dockerfile)
	}

	// add docker build arguments: Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
	if len(options.PersistentVars) > 0 {
		for varName, varValue := range options.PersistentVars {
			d.driver.AddBuildArgs(varName, varValue.(string))
		}
	}

	// add docker build arguments: Vars contains the variables defined by the user on execution time and has precedences over the default values
	if len(options.Vars) > 0 {
		for varName, varValue := range options.Vars {
			d.driver.AddBuildArgs(varName, varValue.(string))
		}
	}

	// add docker tags
	if len(options.Tags) > 0 {
		for _, tag := range options.Tags {

			imageTaggedAux, err := image.NewImage(options.ImageName, tag, options.RegistryHost, options.RegistryNamespace)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
			imageTaggedName, err := imageTaggedAux.DockerNormalizedNamed()
			if err != nil {
				return errors.New(errContext, err.Error())
			}

			d.driver.AddTags(imageTaggedName)
		}
	}

	if len(options.Labels) > 0 {
		for label, value := range options.Labels {
			d.driver.AddLabel(label, value)
		}
	}

	// add docker build arguments
	if options.ImageFromRegistryNamespace != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], options.ImageFromRegistryNamespace)
	}

	if options.ImageFromName != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], options.ImageFromName)
	}

	if options.ImageFromVersion != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], options.ImageFromVersion)
	}

	// add docker build arguments: map de command flag options to build argurments
	if options.ImageFromRegistryHost != "" {
		d.driver.AddBuildArgs(options.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], options.ImageFromRegistryHost)

		if options.PullAuthUsername != "" && options.PullAuthPassword != "" {
			d.driver.AddAuth(options.PullAuthUsername, options.PullAuthPassword, options.ImageFromRegistryHost)
		}
	}

	if options.PushAuthUsername != "" && options.PushAuthPassword != "" {
		d.driver.AddAuth(options.PushAuthUsername, options.PushAuthPassword, options.RegistryHost)
	}

	if options.PushImageAfterBuild {
		d.driver.WithPushAfterBuild()

		if options.PushAuthUsername != "" && options.PushAuthPassword != "" {
			d.driver.AddPushAuth(options.PushAuthUsername, options.PushAuthPassword)
		}
	}

	if options.PullParentImage {
		d.driver.WithPullParentImage()
	}

	if options.RemoveImageAfterBuild {
		d.driver.WithRemoveAfterPush()
	}

	if options.BuilderOptions.Context == nil || len(options.BuilderOptions.Context) == 0 {
		return errors.New(errContext, "Docker building context has not been defined on build options")
	}

	err = d.driver.AddBuildContext(options.BuilderOptions.Context...)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	responseOutputPrefix := options.OutputPrefix
	if responseOutputPrefix == "" {
		responseOutputPrefix = imageName
	}

	d.driver.WithResponse(d.writer, responseOutputPrefix)
	d.driver.WithUseNormalizedNamed()

	err = d.driver.Run(ctx)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
