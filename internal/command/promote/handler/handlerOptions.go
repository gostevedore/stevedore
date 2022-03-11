package promote

type HandlerOptions struct {
	// DryRun is a flag to indicate if the promote should be a dry run
	DryRun bool
	// EnableSemanticVersionTags is a flag to indicate whether to generate semantic version tags
	EnableSemanticVersionTags bool
	// SourceImageName is the name of the image to promote
	SourceImageName string
	// TargetImageName is the name of the image to promote to
	TargetImageName string
	// TargetImageRegistryNamespace is the namespace of the registry to use as target
	TargetImageRegistryNamespace string
	// TargetImageRegistryHost is the host of the registry to use as target
	TargetImageRegistryHost string
	// TargetImageTags is the list of tags to generate
	TargetImageTags []string
	// RemoveTargerImageTags is a flag to indicate whether to remove from local generated image tags
	RemoveTargetImageTags bool
	// DEPRECATED_RemoveTargetImageTags is a flag to indicate whether to remove from local generated image tags
	DEPRECATED_RemoveTargetImageTags bool
	// SemanticVersionTagsTemplates is the list of semantic version tags templates
	SemanticVersionTagsTemplates []string
	// PromoteSourceImageTag is the tag to promote
	PromoteSourceImageTag bool
	// RemoteSourceImage is the flag to indicate whether to promote from remote source image
	RemoteSourceImage bool
}
