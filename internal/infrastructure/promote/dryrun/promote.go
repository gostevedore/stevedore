package dryrun

import (
	"context"
	"fmt"
	"io"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

type DryRunPromote struct {
	writer io.Writer
}

func NewDryRunPromote(w io.Writer) *DryRunPromote {
	return &DryRunPromote{
		writer: w,
	}
}

func (p *DryRunPromote) Promote(ctx context.Context, options *image.PromoteOptions) error {

	fmt.Fprintln(p.writer)
	fmt.Fprintf(p.writer, " pull_auth_username: %s\n", options.PullAuthUsername)
	fmt.Fprintf(p.writer, " push_auth_username: %s\n", options.PushAuthUsername)
	fmt.Fprintf(p.writer, " remove_local_images_after_push: %t\n", options.RemoveTargetImageTags)
	fmt.Fprintf(p.writer, " source_image_name: %s\n", options.SourceImageName)
	fmt.Fprintf(p.writer, " target_image_name: %s\n", options.TargetImageName)
	if len(options.TargetImageTags) > 0 {
		fmt.Fprintln(p.writer, " target_image_tags:")
		for _, tag := range options.TargetImageTags {
			fmt.Fprintf(p.writer, "  - %s\n", tag)
		}
	}
	fmt.Fprintf(p.writer, " use_image_from_remote_source: %t\n", options.RemoteSourceImage)
	fmt.Fprintln(p.writer)

	return nil
}
