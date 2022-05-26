package promote

type Options struct {
	// DryRun is a flag to indicate if the promote should be a dry run
	DryRun bool
	// EnableSemanticVersionTags flag generate semantic versioning tags when is true
	EnableSemanticVersionTags bool
	// TargetImageName is the target image name
	TargetImageName string
	// TargetImageRegistryNamespace is the target namespace name
	TargetImageRegistryNamespace string
	// TargetImageRegistryHost is the target registry host
	TargetImageRegistryHost string
	// TargetImageTags list of extra tags for the target image
	TargetImageTags []string
	// PromoteSourceImageTag push source image to registry
	PromoteSourceImageTag bool
	// RemoveTargetImageTags flag removes all images from local host once the image is promoted
	RemoveTargetImageTags bool
	// RemoteSourceImage flag indicates to use an image from remote source
	RemoteSourceImage bool
	// SourceImageName is the source image name
	SourceImageName string
	// SemanticVersionTagsTemplates is a list of templates to use to generate semantic version tags
	SemanticVersionTagsTemplates []string
}
