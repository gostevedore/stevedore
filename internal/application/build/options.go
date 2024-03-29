package build

// Options
type Options struct {
	// AnsibleConnectionLocal is the local connection to use on ansible driver
	AnsibleConnectionLocal bool
	// AnsibleIntermediateContainerName is the name of an intermediate container that can be used during ansible build process
	AnsibleIntermediateContainerName string
	// AnsibleInventoryPath is the path to the ansible inventory file ??
	AnsibleInventoryPath string
	// AnsibleLimit is the ansible limit ??
	AnsibleLimit string
	// EnableSemanticVersionTags is a flag to enable semantic version tags
	EnableSemanticVersionTags bool
	// ImageFromName is the parent's image name
	ImageFromName string `yaml:"image_from_name"`
	// ImageFromRegistryHost is the parent's image registry host
	ImageFromRegistryHost string `yaml:"image_from_registry_host"`
	// ImageFromRegistryNamespace is the parent's image namespace
	ImageFromRegistryNamespace string `yaml:"image_from_namespace"`
	// ImageFromVersion is the paren't image version
	ImageFromVersion string `yaml:"image_from_version"`
	// ImageName is the name of the image to build
	ImageName string
	// ImageRegistryHost is the registry's host of the image to be built
	ImageRegistryHost string `yaml:"image_registry_host"`
	// ImageRegistryNamespace is the namespace of the image to be built
	ImageRegistryNamespace string `yaml:"image_namespace"`
	// ImageVersions is a list of versions to build
	ImageVersions []string
	// Lables is a list of labels to add to the image
	Labels map[string]string
	// PersistentLabels is a persistent labels list to be sent to driver
	PersistentLabels map[string]string `yaml:"persistent_labels"`
	// PersistentVars is a persistent variables list to be sent to driver
	PersistentVars map[string]interface{} `yaml:"persistent_variables"`
	// PullParentImage flag indicate whether to pull the parent image before building
	PullParentImage bool `yaml:"pull_parent_image"`
	// PushImageAfterBuild flag indicate whether to push the image to the registry once it has been built
	PushImageAfterBuild bool `yaml:"push_image_after_build"`
	// RemoveImagesAfterPush flag indicate whether to remove the image after build
	RemoveImagesAfterPush bool
	// SemanticVersionTagsTemplate are the semantic version tags templates to generate automatically
	SemanticVersionTagsTemplates []string `yaml:"semantic_version_tags_template"`
	// Tags is a list of tags to generate
	Tags []string `yaml:"tags"`
	// Vars is a variables list to be sent to driver
	Vars map[string]interface{} `yaml:"variables"`
}

// Copy returns a copy of the Options
func (o *Options) Copy() *Options {
	copy := &Options{}

	copy.ImageName = o.ImageName

	copy.ImageVersions = []string{}
	for _, version := range o.ImageVersions {
		copy.ImageVersions = append(copy.ImageVersions, version)
	}
	copy.ImageRegistryNamespace = o.ImageRegistryNamespace
	copy.ImageRegistryHost = o.ImageRegistryHost
	copy.EnableSemanticVersionTags = o.EnableSemanticVersionTags
	copy.ImageFromName = o.ImageFromName
	copy.ImageFromRegistryNamespace = o.ImageFromRegistryNamespace
	copy.ImageFromRegistryHost = o.ImageFromRegistryHost
	copy.ImageFromVersion = o.ImageFromVersion

	copy.PushImageAfterBuild = o.PushImageAfterBuild
	copy.RemoveImagesAfterPush = o.RemoveImagesAfterPush
	copy.AnsibleConnectionLocal = o.AnsibleConnectionLocal

	copy.PersistentVars = map[string]interface{}{}
	for name, value := range o.PersistentVars {
		copy.PersistentVars[name] = value
	}

	copy.SemanticVersionTagsTemplates = []string{}
	for _, semanticVersionTagTemplate := range o.SemanticVersionTagsTemplates {
		copy.SemanticVersionTagsTemplates = append(copy.SemanticVersionTagsTemplates, semanticVersionTagTemplate)
	}

	copy.Tags = []string{}
	for _, tag := range o.Tags {
		copy.Tags = append(copy.Tags, tag)
	}

	copy.Vars = map[string]interface{}{}
	for name, value := range o.Vars {
		copy.Vars[name] = value
	}

	copy.Labels = map[string]string{}
	for name, value := range o.Labels {
		copy.Labels[name] = value
	}

	return copy
}
