package docker

import (
	"context"
	"io"
)

// // Promoter
// type Promoter interface {
// 	Promote(context.Context, *image.PromoteOptions) error
// }

// DockerCopier
type DockerCopier interface {
	DockerCopyConfigurer
	DockerCopyAuther
	Run(context.Context) error
}

// DockerCopyAuther
type DockerCopyAuther interface {
	AddAuth(string, string) error
	AddPullAuth(string, string) error
	AddPushAuth(string, string) error
}

// DockerCopyConfigurer
type DockerCopyConfigurer interface {
	WithSourceImage(source string)
	WithTags(tags []string)
	WithTargetImage(target string)
	WithRemoteSource()
	WithRemoveAfterPush()
	WithResponse(io.Writer, string)
	WithUseNormalizedNamed()
}
