package build

// HandlerOptions is the options for the build command
type HandlerOptions struct {
	// AnsibleConnectionLocal if is true ansible driver uses local connection
	AnsibleConnectionLocal bool
	// AnsibleInventoryPath is the path to the ansible inventory file ??
	AnsibleInventoryPath string
	// AnsibleLimit is the ansible limit ??
	AnsibleLimit string

	// BuildBuilderName is the name of the builder to use
	BuildBuilderName string
	// DryRun is true if the build should be a dry run
	DryRun bool
	// BuildOnCascade if is true the build should be cascaded
	BuildOnCascade bool
	// CascadeDepth is the number of levels to build when build on cascade is executed
	CascadeDepth int
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
	// NumWorkers is the number of workers to use
	NumWorkers int
	// PersistentVars is the list of persistent vars to use
	PersistentVars []string
	// PushImagesAfterBuild is true if the images should be pushed
	PushImagesAfterBuild bool
	// SemanticVersionTagsTemplates is the list of semantic version tags templates
	SemanticVersionTagsTemplates []string
	// Tags is the list of tags to generate
	Tags []string
	// Vars is the list of vars to use
	Vars []string
	// Versions is the list of versions to generate
	Versions []string
}
