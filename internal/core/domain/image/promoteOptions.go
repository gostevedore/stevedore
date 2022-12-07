package image

import (
	data "github.com/apenella/go-common-utils/data"
)

// PromoteOptions
type PromoteOptions struct {
	// TargetImageName is the target image name
	TargetImageName string `yaml:"target_image_name"`
	// TargetImageTags list of extra tags for the target image
	TargetImageTags []string `yaml:"target_image_tags"`
	// RemoveTargetImageTags flag removes all images from local host once the image is promoted
	RemoveTargetImageTags bool `yaml:"remove_local_images_after_push"`
	// RemoteSourceImage flag use an image from remote source
	RemoteSourceImage bool `yaml:"use_image_from_remote_source"`
	// SourceImageName is the source image name
	SourceImageName string `yaml:"source_image_name"`
	// PullAuthUsername
	PullAuthUsername string `yaml:"pull_auth_username"`
	// PullAuthPassword
	PullAuthPassword string
	// PushAuthUsername
	PushAuthUsername string `yaml:"push_auth_username"`
	// PushAuthPassword
	PushAuthPassword string
}

// String TODO
func (o *PromoteOptions) String() string {
	str, err := data.ObjectToYamlString(o)
	if err != nil {
		return ""
	}

	return str
}
