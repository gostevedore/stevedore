package dockerdriver

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/driver"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
	"github.com/gostevedore/stevedore/internal/image"
	"gopkg.in/yaml.v2"
)

const (
	DriverName = "docker"

	ImageTagVersionSeparator   = ":"
	ImageTagNamespaceSeparator = "/"
	ImageTagRegistrySeparator  = "/"

	builderConfOptionsContextKey    = "context"
	builderConfOptionsDockerfileKey = "dockerfile"
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

	imageNameURL := &image.ImageURL{
		Name: options.ImageName,
	}

	if options.ImageVersion != "" {
		imageNameURL.Tag = options.ImageVersion
	}

	if options.RegistryNamespace != "" {
		imageNameURL.Namespace = options.RegistryNamespace
	}

	if options.RegistryHost != "" {
		imageNameURL.Registry = options.RegistryHost
	}

	imageName, err = imageNameURL.URL()
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	d.driver.WithImageName(imageName)

	// add docker build context to cmd instance
	builderConfOptions := options.BuilderOptions

	dockerfile, exists := builderConfOptions[builderConfOptionsDockerfileKey]
	if exists {
		d.driver.WithDockerfile(dockerfile.(string))
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

			imageURL := &image.ImageURL{
				Name: options.ImageName,
				Tag:  tag,
			}

			if options.RegistryNamespace != "" {
				imageURL.Namespace = options.RegistryNamespace
			}
			if options.RegistryHost != "" {
				imageURL.Registry = options.RegistryHost
			}

			url, err := imageURL.URL()
			if err != nil {
				return errors.New(errContext, err.Error())
			}
			d.driver.AddTags(url)
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

	if options.PushImages {
		d.driver.WithPushAfterBuild()

		if options.PushAuthUsername != "" && options.PushAuthPassword != "" {
			d.driver.AddPushAuth(options.PushAuthUsername, options.PushAuthPassword)
		}
	}

	if options.PullParentImage {
		d.driver.WithPullParentImage()
	}

	builderConfOptionsRaw, exists := builderConfOptions[builderConfOptionsContextKey]

	if !exists {
		return errors.New(errContext, "Docker building context has not been defined on build options")
	}

	dockerBuildContextOptionsList := []*buildcontext.DockerBuildContextOptions{}

	DockerBuildContextOptionsRawList := struct {
		Context []*buildcontext.DockerBuildContextOptions
	}{}

	builderConfOptionsRaw = fmt.Sprintf("%s:%s", builderConfOptionsContextKey, builderConfOptionsRaw.(string))

	err = yaml.Unmarshal([]byte(builderConfOptionsRaw.(string)), &DockerBuildContextOptionsRawList)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Docker build context options are not properly configured\n found:\n%s\n", builderConfOptionsRaw.(string)), err)
	}

	if len(DockerBuildContextOptionsRawList.Context) == 0 {
		return errors.New(errContext, fmt.Sprintf("There is no docker build context definition found on:\n%s", builderConfOptionsRaw.(string)))
	}

	for _, dockerBuildContextOptions := range DockerBuildContextOptionsRawList.Context {
		dockerBuildContextOptionsList = append(dockerBuildContextOptionsList, dockerBuildContextOptions)
	}

	err = d.driver.AddBuildContext(dockerBuildContextOptionsList...)
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

// func X_NewDockerDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {

// 	var err error
// 	var imageName string
// 	var dockerCli *client.Client
// 	var dockerBuildContext dockercontext.DockerBuildContexter
// 	var auth *credentials.RegistryUserPassAuth

// 	if o == nil {
// 		return nil, errors.New("(build::NewDockerDriver)", "Build options are nil")
// 	}

// 	if ctx == nil {
// 		return nil, errors.New("(build::NewDockerDriver)", "Context is nil")
// 	}

// 	if o.ImageName == "" {
// 		return nil, errors.New("(build::NewDockerDriver)", "Image name is not set")
// 	}

// 	// gave a name to image
// 	imageNameURL := &image.ImageURL{
// 		Name: o.ImageName,
// 	}

// 	if o.ImageVersion != "" {
// 		imageNameURL.Tag = drivercommon.SanitizeTag(o.ImageVersion)
// 	}

// 	if o.RegistryNamespace != "" {
// 		imageNameURL.Namespace = o.RegistryNamespace
// 	}

// 	if o.RegistryHost != "" {
// 		imageNameURL.Registry = o.RegistryHost
// 	}

// 	imageName, err = imageNameURL.URL()
// 	if err != nil {
// 		return nil, errors.New("(build::NewDockerDriver)", "Image url for image '"+o.ImageName+"' could not be created", err)
// 	}

// 	// create docker sdk client
// 	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		return nil, errors.New("(images::NewDockerDriver)", "Docker client could not be created", err)
// 	}
// 	dockerCli.NegotiateAPIVersion(ctx)

// 	// generate output prefix to include into response
// 	if o.OutputPrefix == "" {
// 		o.OutputPrefix = imageNameURL.Name
// 		if o.ImageVersion != "" {
// 			o.OutputPrefix = strings.Join([]string{o.OutputPrefix, imageNameURL.Tag}, ":")
// 		}
// 	}

// 	// create docker build cmd instance
// 	dockerBuilder := &build.DockerBuildCmd{
// 		Cli:       dockerCli,
// 		ImageName: imageName,
// 		Response: response.NewDefaultResponse(
// 			response.WithTransformers(
// 				transformer.Prepend(o.OutputPrefix),
// 			),
// 			response.WithWriter(console.GetConsole()),
// 		),
// 	}

// 	// add docker build context to cmd instance
// 	builderConfOptions := o.BuilderOptions

// 	context, exists := builderConfOptions[builderConfOptionsContextKey]
// 	if !exists {
// 		return nil, errors.New("(build::NewDockerDriver)", "Docker building context has not been defined on build options")
// 	}
// 	dockerBuildContext, err = builddockercontext.GenerateDockerBuildContext(context)
// 	if err != nil {
// 		return nil, errors.New("(build::NewDockerDriver)", "Docker build context could not be extracted", err)
// 	}

// 	err = dockerBuilder.AddBuildContext(dockerBuildContext)
// 	if err != nil {
// 		return nil, errors.New("(images::NewDockerDriver)", "Error adding docker build context", err)
// 	}

// 	// include dockerfile location
// 	if o.Dockerfile != "" {
// 		dockerBuilder.ImageBuildOptions.Dockerfile = o.Dockerfile
// 	} else {
// 		dockerfile, exists := builderConfOptions["dockerfile"]
// 		if exists {
// 			dockerBuilder.ImageBuildOptions.Dockerfile = dockerfile.(string)
// 		}
// 	}

// 	// add docker tags
// 	if len(o.Tags) > 0 {
// 		for _, tag := range o.Tags {

// 			tag = drivercommon.SanitizeTag(tag)

// 			imageURL := &image.ImageURL{
// 				Name: o.ImageName,
// 				Tag:  tag,
// 			}

// 			if o.RegistryNamespace != "" {
// 				imageURL.Namespace = o.RegistryNamespace
// 			}
// 			if o.RegistryHost != "" {
// 				imageURL.Registry = o.RegistryHost
// 			}

// 			url, err := imageURL.URL()
// 			if err != nil {
// 				return nil, errors.New("(build::NewDockerDriver)", "Image url for image '"+o.ImageName+"' and tag '"+tag+"' could not be created", err)
// 			}
// 			dockerBuilder.AddTags(url)
// 		}
// 	}

// 	// add docker build arguments: Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
// 	if len(o.PersistentVars) > 0 {
// 		for varName, varValue := range o.PersistentVars {
// 			dockerBuilder.AddBuildArgs(varName, varValue.(string))
// 		}
// 	}

// 	// add docker build arguments: Vars contains the variables defined by the user on execution time and has precedences over the default values
// 	if len(o.Vars) > 0 {
// 		for varName, varValue := range o.Vars {
// 			dockerBuilder.AddBuildArgs(varName, varValue.(string))
// 		}
// 	}

// 	// add docker build arguments
// 	if o.ImageFromRegistryNamespace != "" {
// 		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], o.ImageFromRegistryNamespace)
// 	}

// 	if o.ImageFromName != "" {
// 		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], o.ImageFromName)
// 	}

// 	if o.ImageFromVersion != "" {
// 		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], o.ImageFromVersion)
// 	}

// 	// add docker build arguments: map de command flag options to build argurments
// 	if o.ImageFromRegistryHost != "" {
// 		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], o.ImageFromRegistryHost)

// 		auth, err = credentials.AchieveCredential(o.ImageFromRegistryHost)
// 		if err == nil {
// 			user := auth.Username
// 			pass := auth.Password
// 			dockerBuilder.AddAuth(user, pass, o.ImageFromRegistryHost)
// 		}
// 	}

// 	// set whether to push automatically images after build is done
// 	if o.PushImages {
// 		dockerBuilder.PushAfterBuild = o.PushImages

// 		auth, err = credentials.AchieveCredential(o.RegistryHost)
// 		if err == nil {
// 			user := auth.Username
// 			pass := auth.Password
// 			dockerBuilder.AddPushAuth(user, pass)

// 			// add auth to build options when it not already set
// 			_, added := dockerBuilder.ImageBuildOptions.AuthConfigs[o.RegistryHost]
// 			if !added {
// 				dockerBuilder.AddAuth(user, pass, o.RegistryHost)
// 			}
// 		}
// 	}

// 	return dockerBuilder, nil
// }
