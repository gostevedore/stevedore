package dockerdriver

import (
	"context"
	"stevedore/internal/build/varsmap"
	"stevedore/internal/credentials"
	drivercommon "stevedore/internal/driver/common"
	builddockercontext "stevedore/internal/driver/docker/context"
	"stevedore/internal/image"
	"stevedore/internal/types"
	"stevedore/internal/ui/console"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/build"
	dockercontext "github.com/apenella/go-docker-builder/pkg/build/context"
	dockerpush "github.com/apenella/go-docker-builder/pkg/push"
	"github.com/docker/docker/client"
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

	builderConfOptions := o.BuilderOptions

	context, exists := builderConfOptions[builderConfOptionsContextKey]
	if !exists {
		return nil, errors.New("(build::NewDockerDriver)", "Docker building context has not been defined on build options")
	}
	dockerBuildContext, err = builddockercontext.GenerateDockerBuildContext(context)
	if err != nil {
		return nil, errors.New("(build::NewDockerDriver)", "Docker build context could not be extracted", err)
	}

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

	if o.OutputPrefix == "" {
		o.OutputPrefix = imageNameURL.Name
		if o.ImageVersion != "" {
			o.OutputPrefix = strings.Join([]string{o.OutputPrefix, imageNameURL.Tag}, ":")
		}
	}

	imageName, err = imageNameURL.URL()
	if err != nil {
		return nil, errors.New("(build::NewDockerDriver)", "Image url for image '"+o.ImageName+"' could not be created", err)
	}

	dockerBuildOptions := &build.DockerBuildOptions{
		ImageName:          imageName,
		BuildArgs:          map[string]*string{},
		DockerBuildContext: dockerBuildContext,
	}

	if o.Dockerfile != "" {
		dockerBuildOptions.Dockerfile = o.Dockerfile
	} else {
		dockerfile, exists := builderConfOptions["dockerfile"]
		if exists {
			dockerBuildOptions.Dockerfile = dockerfile.(string)
		}
	}

	if o.PushImages {
		dockerBuildOptions.PushAfterBuild = o.PushImages
	}

	// Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
	if len(o.PersistentVars) > 0 {
		for varName, varValue := range o.PersistentVars {
			dockerBuildOptions.AddBuildArgs(varName, varValue.(string))
		}
	}

	// Vars contains the variables defined by the user on execution time and has precedences over the default values
	if len(o.Vars) > 0 {
		for varName, varValue := range o.Vars {
			dockerBuildOptions.AddBuildArgs(varName, varValue.(string))
		}
	}

	// map de command flag options to build argurments
	if o.ImageFromRegistryHost != "" {
		dockerBuildOptions.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], o.ImageFromRegistryHost)

		auth, err = credentials.AchieveCredential(o.ImageFromRegistryHost)
		if err == nil {
			user := auth.Username
			pass := auth.Password
			dockerBuildOptions.AddAuth(user, pass, o.ImageFromRegistryHost)
		}
	}

	if o.ImageFromRegistryNamespace != "" {
		dockerBuildOptions.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], o.ImageFromRegistryNamespace)
	}

	if o.ImageFromName != "" {
		dockerBuildOptions.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], o.ImageFromName)
	}

	if o.ImageFromVersion != "" {
		dockerBuildOptions.AddBuildArgs(o.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], o.ImageFromVersion)
	}

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
			dockerBuildOptions.AddTags(url)
		}
	}

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, errors.New("(images::NewDockerDriver)", "Docker client could not be created", err)
	}
	dockerCli.NegotiateAPIVersion(ctx)

	dockerPushOptions := &dockerpush.DockerPushOptions{
		ImageName: imageName,
	}

	auth, err = credentials.AchieveCredential(o.RegistryHost)
	if err == nil {
		user := auth.Username
		pass := auth.Password
		dockerPushOptions.AddAuth(user, pass)

		// add auth to build options when it not already set
		if dockerBuildOptions.Auth == nil {
			dockerBuildOptions.AddAuth(user, pass, o.RegistryHost)
		}
	}

	dockerBuilder := &build.DockerBuildCmd{
		Writer:             console.GetConsole(),
		Context:            ctx,
		Cli:                dockerCli,
		DockerBuildOptions: dockerBuildOptions,
		DockerPushOptions:  dockerPushOptions,
		ExecPrefix:         o.OutputPrefix,
	}

	return dockerBuilder, nil
}
