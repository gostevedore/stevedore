package types

import (
	"github.com/apenella/go-common-utils/data"
)

// promoteOptions
type PromoteOptions struct {
	DryRun                        bool     `yaml:"dry_run"`
	EnableSemanticVersionTags     bool     `yaml:"enable_semantic_version_tags"`
	ImagePromoteName              string   `yaml:"image_promote_name"`
	ImagePromoteRegistryNamespace string   `yaml:"image_promote_registry_namespace"`
	ImagePromoteRegistryHost      string   `yaml:"image_promote_registry_host"`
	ImagePromoteTags              []string `yaml:"image_promote_tags"`
	RemovePromotedTags            bool     `yaml:"remove_promoted_tags"`
	ImageName                     string   `yaml:"image_name"`
	OutputPrefix                  string   `yaml:"output_prefix"`
	SemanticVersionTagsTemplates  []string `yaml:"semantic_version_tags_templates"`
}

func (o *PromoteOptions) String() string {
	str, err := data.ObjectToYamlString(o)
	if err != nil {
		return ""
	}

	return str
}
