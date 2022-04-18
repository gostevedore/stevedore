package build

// HandlerOptions is the options for the build command
type HandlerOptions struct {
	// AnsibleConnectionLocal if is true ansible driver uses local connection
	AnsibleConnectionLocal bool
	// AnsibleIntermediateContainerName is the name of an intermediate container that can be used during ansible build process
	AnsibleIntermediateContainerName string
	// AnsibleInventoryPath is the path to the ansible inventory file ??
	AnsibleInventoryPath string
	// AnsibleLimit is the ansible limit ??
	AnsibleLimit string
	// BuildOnCascade if is true the build should be cascaded: ???
	BuildOnCascade bool
	// CascadeDepth is the number of levels to build when build on cascade is executed: ???
	CascadeDepth int

	// // Concurrency is the number of images builds that can be excuted at the same time
	// Concurrency int
	// // Debug if is true debug mode is enabled: ???
	// Debug bool

	// DryRun if is true, the build is executed in dry run mode
	DryRun bool
	// EnableSemanticVersionTags if is true semantic version tags are generated
	EnableSemanticVersionTags bool
	// ImageFromName is the name of the image to use as source
	ImageFromName string
	// ImageFromRegistryHost is the host of the registry to use as source
	ImageFromRegistryHost string
	// ImageFromRegistryNamespace is the namespace of the registry to use as source
	ImageFromRegistryNamespace string
	// ImageFromVersion is the version of the image to use as source
	ImageFromVersion string
	// ImageName is the name of the image to build
	ImageName string
	// ImageRegistryHost is the host of the registry to use
	ImageRegistryHost string
	// ImageRegistryNamespace is the namespace of the registry to use
	ImageRegistryNamespace string
	// Labels is the list of labes to assign to the image
	Labels []string
	// PersistentVars is the list of persistent vars to use
	PersistentVars []string
	// PullParentImage if is true the parent image is pull
	PullParentImage bool
	// PushImagesAfterBuild if is true the image is pushed after build
	PushImagesAfterBuild bool
	// RemoveImagesAfterPush if is true the images are removed from local after push
	RemoveImagesAfterPush bool
	// SemanticVersionTagsTemplates is the list of semantic version tags templates
	SemanticVersionTagsTemplates []string
	// Tags is the list of tags to generate
	Tags []string
	// Vars is the list of vars to use
	Vars []string
	// Versions is the list of versions to generate
	Versions []string
}
