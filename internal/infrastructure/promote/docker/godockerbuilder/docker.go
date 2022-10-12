package godockerbuilder

import (
	"context"
	"io"

	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/copy"
	"github.com/apenella/go-docker-builder/pkg/response"
)

// DockerCopier
type DockerCopy struct {
	cmd *copy.DockerImageCopyCmd
}

func NewDockerCopy(cmd *copy.DockerImageCopyCmd) *DockerCopy {
	return &DockerCopy{
		cmd: cmd,
	}
}

func (c *DockerCopy) WithSourceImage(source string) {
	c.cmd = c.cmd.WithSourceImage(source)
}

// WithTags
func (c *DockerCopy) WithTags(tags []string) {
	c.cmd = c.cmd.WithTags(tags)
}

// WithTargetImage
func (c *DockerCopy) WithTargetImage(target string) {
	c.cmd = c.cmd.WithTargetImage(target)
}

// WithRemoteSource
func (c *DockerCopy) WithRemoteSource() {
	c.cmd = c.cmd.WithRemoteSource()
}

// WithRemoveAfterPush
func (c *DockerCopy) WithRemoveAfterPush() {
	c.cmd = c.cmd.WithRemoveAfterPush()
}

// WithResponse
func (c *DockerCopy) WithResponse(w io.Writer, prefix string) {
	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend(prefix),
		),
		response.WithWriter(w),
	)

	c.cmd = c.cmd.WithResponse(res)
}

// WithUseNormalizedNamed
func (c *DockerCopy) WithUseNormalizedNamed() {
	c.cmd = c.cmd.WithUseNormalizedNamed()
}

// AddAuth
func (c *DockerCopy) AddAuth(username string, password string) error {
	return c.cmd.AddAuth(username, password)
}

// AddPullAuth
func (c *DockerCopy) AddPullAuth(username string, password string) error {
	return c.cmd.AddPullAuth(username, password)
}

// AddPushAuth
func (c *DockerCopy) AddPushAuth(username string, password string) error {
	return c.cmd.AddPushAuth(username, password)
}

// Run
func (c *DockerCopy) Run(ctx context.Context) error {
	return c.cmd.Run(ctx)
}
