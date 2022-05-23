package dockerpromote

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/promote"
)

const (
	DockerImageFilterReference = "reference"
)

type DockerPromete struct {
	cmd promote.DockerCopier
	//	logger Logger
	writer io.Writer
}

func NewDockerPromote(cmd promote.DockerCopier, w io.Writer) *DockerPromete {

	if w == nil {
		w = os.Stdout
	}

	return &DockerPromete{
		cmd:    cmd,
		writer: w,
	}
}

func (p *DockerPromete) Promote(ctx context.Context, options *image.PromoteOptions) error {

	var err error

	contextError := "(docker::Promote)"

	if p.cmd == nil {
		return errors.New(contextError, "Command to copy docker images must be initialized before promote an image to docker registry")
	}

	if p.writer == nil {
		return errors.New(contextError, "Writer must be initialized before promote an image to docker registry")
	}

	if options == nil {
		return errors.New(contextError, "Image could not be promoted because options must be defined")
	}

	if options.SourceImageName == "" {
		return errors.New(contextError, "Image could not be promoted because source image name must be defined on promote options")
	}

	if options.TargetImageName == "" {
		return errors.New(contextError, "Image could not be promoted because target image name must be defined on promote options")
	}

	if options.RemoteSourceImage {
		if options.PullAuthUsername != "" && options.PullAuthPassword != "" {
			err = p.cmd.AddPullAuth(options.PullAuthUsername, options.PullAuthPassword)
			if err != nil {
				return errors.New(contextError, fmt.Sprintf("Image '%s' could not be promoted because is not possible to achieve pull credentials", options.SourceImageName), err)
			}
		}

		p.cmd.WithRemoteSource()
	}

	if options.PushAuthUsername != "" && options.PushAuthPassword != "" {
		err = p.cmd.AddPushAuth(options.PushAuthUsername, options.PushAuthPassword)
		if err != nil {
			return errors.New(contextError, fmt.Sprintf("Image '%s' could not be promoted because is not possible to achieve push credentials", options.SourceImageName), err)
		}
	}

	if options.RemoveTargetImageTags {
		p.cmd.WithRemoveAfterPush()
	}

	p.cmd.WithSourceImage(options.SourceImageName)
	p.cmd.WithTargetImage(options.TargetImageName)
	p.cmd.WithTags(options.TargetImageTags)
	p.cmd.WithUseNormalizedNamed()
	p.cmd.WithResponse(p.writer, options.TargetImageName)

	err = p.cmd.Run(ctx)
	if err != nil {
		return errors.New(contextError, fmt.Sprintf("Image '%s' could not be promoted", options.SourceImageName), err)
	}

	return nil
}
