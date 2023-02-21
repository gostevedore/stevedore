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
	fmt.Fprintln(p.writer, fmt.Sprintf(" pull_auth_username: %s", options.PullAuthUsername))
	fmt.Fprintln(p.writer, fmt.Sprintf(" push_auth_username: %s", options.PushAuthUsername))
	fmt.Fprintln(p.writer, fmt.Sprintf(" remove_local_images_after_push: %t", options.RemoveTargetImageTags))
	fmt.Fprintln(p.writer, fmt.Sprintf(" source_image_name: %s", options.SourceImageName))
	fmt.Fprintln(p.writer, fmt.Sprintf(" target_image_name: %s", options.TargetImageName))
	if len(options.TargetImageTags) > 0 {
		fmt.Fprintln(p.writer, " target_image_tags:")
		for _, tag := range options.TargetImageTags {
			fmt.Fprintf(p.writer, fmt.Sprintf("  - %s\n", tag))
		}
	}
	fmt.Fprintln(p.writer, fmt.Sprintf(" use_image_from_remote_source: %t", options.RemoteSourceImage))
	fmt.Fprintln(p.writer)

	return nil
}
