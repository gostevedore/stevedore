package types

// promoteOptions
type PromoteOptions struct {
	DryRun                        bool
	EnableSemanticVersionTags     bool
	ImagePromoteName              string
	ImagePromoteRegistryNamespace string
	ImagePromoteRegistryHost      string
	ImagePromoteTags              []string
	RemovePromotedTags            bool
	ImageName                     string
	OutputPrefix                  string
	SemanticVersionTagsTemplate   []string
}
