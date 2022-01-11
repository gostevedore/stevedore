package driver

import (
	"github.com/gostevedore/stevedore/internal/builders/builder"
)

// BuildDriverOptions options required by driver to build an image
type BuildDriverOptions struct {
	// BuilderName is the intermediate container or build stage container name
	BuilderName string `yaml:"builder_name"`
	// BuilderOptions are those options comming from builder
	BuilderOptions *builder.BuilderOptions `yaml:"builder_options"`
	// BuilderVarMappings are those variables name that will be automatically generated by builder and set to the driver for building the image
	BuilderVarMappings map[string]string `yaml:"builder_variables_mapping"`
	// ConnectionLocal indicates whether to use a local connection during the build when ansible-playbook driver is used
	ConnectionLocal bool `yaml:"connection_local"`
	// ImageFromName is the parent's image name
	ImageFromName string `yaml:"image_from_name"`
	// ImageFromRegistryNamespace is the parent's image namespace
	ImageFromRegistryNamespace string `yaml:"image_from_registry_namespace"`
	// ImageFromRegistryHost is the parent's image registry host
	ImageFromRegistryHost string `yaml:"image_from_registry_host"`
	// ImageFromVersion is the paren't image version
	ImageFromVersion string `yaml:"image_from_version"`
	// ImageName is the name of the image to be built
	ImageName string `yaml:"image_name"`
	// ImageVersion is the version of the image to be built
	ImageVersion string `yaml:"image_version"`
	// Lables is a list of labels to add to the image
	Labels map[string]string `yaml:"labels"`
	// OutputPrefix prefixes each output line
	OutputPrefix string `yaml:"output_prefix"`
	// PersistentVars is a persistent variables list to be sent to driver
	PersistentVars map[string]interface{} `yaml:"persistent_variables"`
	// RegistryNamespace is the namespace of the image to be built
	RegistryNamespace string `yaml:"image_registry_namespace"`
	// RegistryHost is the registry's host of the image to be built
	RegistryHost string `yaml:"image_registry_host"`
	// PullAuthUsername is the username to use for pulling the image
	PullAuthUsername string `yaml:"pull_auth_username"`
	// PullAuthPassword is the password to use for pulling the image
	PullAuthPassword string `yaml:"pull_auth_password"`
	// PullParentImage indicates whether to pull the parent image
	PullParentImage bool `yaml:"pull_parent_image"`
	// PushAuthUsername is the username to use for pushing the image
	PushAuthUsername string `yaml:"push_auth_username"`
	// PushAuthPassword is the password to use for pushing the image
	PushAuthPassword string `yaml:"push_auth_password"`
	// PushImageAfterBuild flag indicate whether to push the image to the registry once it has been built
	PushImageAfterBuild bool `yaml:"push_image_after_build"`
	// RemoveImageAfterBuild flag indicate whether to remove the image after build
	RemoveImageAfterBuild bool `yaml:"remove_image_after_build"`
	// Tags is a list of tags to generate
	Tags []string `yaml:"tags"`
	// Vars is a variables list to be sent to driver
	Vars map[string]interface{} `yaml:"variables"`
}
