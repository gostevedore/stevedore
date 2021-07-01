package dockerdriver

import (
	"context"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/build"
	dockercontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/docker/client"
	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/credentials"
	drivercommon "github.com/gostevedore/stevedore/internal/driver/common"
	builddockercontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

const (
	DriverName = "docker"

	ImageTagVersionSeparator   = ":"
	ImageTagNamespaceSeparator = "/"
	ImageTagRegistrySeparator  = "/"

	builderConfOptionsContextKey = "context"
)

func NewDockerDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {

	var err error
	var imageName string
	var dockerCli *client.Client
	var dockerBuildContext dockercontext.DockerBuildContexter
	var auth *credentials.RegistryUserPassAuth

	if o == nil {
		return nil, errors.New("(build::NewDockerDriver)", "Build options are nil")
	}

	if ctx == nil {
		return nil, errors.New("(build::NewDockerDriver)", "Context is nil")
	}

	if o.ImageName == "" {
		return nil, errors.New("(build::NewDockerDriver)", "Image name is not set")
	}

	// gave a name to image
	imageNameURL := &image.ImageURL{
		Name: o.ImageName,
	}

	if o.ImageVersion != "" {
		imageNameURL.Tag = drivercommon.SanitizeTag(o.ImageVersion)
	}

	if o.RegistryNamespace != "" {
		imageNameURL.Namespace = o.RegistryNamespace
	}

	if o.RegistryHost != "" {
		imageNameURL.Registry = o.RegistryHost
	}

	imageName, err = imageNameURL.URL()
	if err != nil {
		return nil, errors.New("(build::NewDockerDriver)", "Image url for image '"+o.ImageName+"' could not be created", err)
	}

	// create docker sdk client
	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, errors.New("(images::NewDockerDriver)", "Docker client could not be created", err)
	}
	dockerCli.NegotiateAPIVersion(ctx)

	// generate output prefix to include into response
	if o.OutputPrefix == "" {
		o.OutputPrefix = imageNameURL.Name
		if o.ImageVersion != "" {
			o.OutputPrefix = strings.Join([]string{o.OutputPrefix, imageNameURL.Tag}, ":")
		}
	}

	// create docker build cmd instance
	dockerBuilder := &build.DockerBuildCmd{
		Cli:       dockerCli,
		ImageName: imageName,
		Response: response.NewDefaultResponse(
			response.WithTransformers(
				transformer.Prepend(o.OutputPrefix),
			),
			response.WithWriter(console.GetConsole()),
		),
	}

	// add docker build context to cmd instance
	builderConfOptions := o.BuilderOptions

	context, exists := builderConfOptions[builderConfOptionsContextKey]
	if !exists {
		return nil, errors.New("(build::NewDockerDriver)", "Docker building context has not been defined on build options")
	}
	dockerBuildContext, err = builddockercontext.GenerateDockerBuildContext(context)
	if err != nil {
		return nil, errors.New("(build::NewDockerDriver)", "Docker build context could not be extracted", err)
	}

	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		return nil, errors.New("(images::NewDockerDriver)", "Error adding docker build context", err)
	}

	// include dockerfile location
	if o.Dockerfile != "" {
		dockerBuilder.ImageBuildOptions.Dockerfile = o.Dockerfile
	} else {
		dockerfile, exists := builderConfOptions["dockerfile"]
		if exists {
			dockerBuilder.ImageBuildOptions.Dockerfile = dockerfile.(string)
		}
	}

	// add docker tags
	if len(o.Tags) > 0 {
		for _, tag := range o.Tags {

			tag = drivercommon.SanitizeTag(tag)

			imageURL := &image.ImageURL{
				Name: o.ImageName,
				Tag:  tag,
			}

			if o.RegistryNamespace != "" {
				imageURL.Namespace = o.RegistryNamespace
			}
			if o.RegistryHost != "" {
				imageURL.Registry = o.RegistryHost
			}

			url, err := imageURL.URL()
			if err != nil {
				return nil, errors.New("(build::NewDockerDriver)", "Image url for image '"+o.ImageName+"' and tag '"+tag+"' could not be created", err)
			}
			dockerBuilder.AddTags(url)
		}
	}

	// add docker build arguments: Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
	if len(o.PersistentVars) > 0 {
		for varName, varValue := range o.PersistentVars {
			dockerBuilder.AddBuildArgs(varName, varValue.(string))
		}
	}

	// add docker build arguments: Vars contains the variables defined by the user on execution time and has precedences over the default values
	if len(o.Vars) > 0 {
		for varName, varValue := range o.Vars {
			dockerBuilder.AddBuildArgs(varName, varValue.(string))
		}
	}

	// add docker build arguments
	if o.ImageFromRegistryNamespace != "" {
		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], o.ImageFromRegistryNamespace)
	}

	if o.ImageFromName != "" {
		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], o.ImageFromName)
	}

	if o.ImageFromVersion != "" {
		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], o.ImageFromVersion)
	}

	// add docker build arguments: map de command flag options to build argurments
	if o.ImageFromRegistryHost != "" {
		dockerBuilder.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], o.ImageFromRegistryHost)

		auth, err = credentials.AchieveCredential(o.ImageFromRegistryHost)
		if err == nil {
			user := auth.Username
			pass := auth.Password
			dockerBuilder.AddAuth(user, pass, o.ImageFromRegistryHost)
		}
	}

	// set whether to push automatically images after build is done
	if o.PushImages {
		dockerBuilder.PushAfterBuild = o.PushImages

		auth, err = credentials.AchieveCredential(o.RegistryHost)
		if err == nil {
			user := auth.Username
			pass := auth.Password
			dockerBuilder.AddPushAuth(user, pass)

			// add auth to build options when it not already set
			_, added := dockerBuilder.ImageBuildOptions.AuthConfigs[o.RegistryHost]
			if !added {
				dockerBuilder.AddAuth(user, pass, o.RegistryHost)
			}
		}
	}

	return dockerBuilder, nil
}
