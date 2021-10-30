package promote

type HandlerOptions struct {
	DryRun                           bool
	EnableSemanticVersionTags        bool
	SourceImageName                  string
	TargetImageName                  string
	TargetImageRegistryNamespace     string
	TargetImageRegistryHost          string
	TargetImageTags                  []string
	RemoveTargetImageTags            bool
	DEPRECATED_RemoveTargetImageTags bool
	SemanticVersionTagsTemplates     []string
	PromoteSourceImageTag            bool
	RemoteSourceImage                bool
}
