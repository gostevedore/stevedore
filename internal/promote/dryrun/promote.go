package dryrunpromote

import (
	"context"
	"fmt"
	"io"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/promote"
)

type DryRunPromote struct {
	writer io.Writer
}

func NewDryRunPromote(cmd promote.DockerCopier, w io.Writer) *DryRunPromote {
	return &DryRunPromote{
		writer: w,
	}
}

func (p *DryRunPromote) Promote(ctx context.Context, options *image.PromoteOptions) error {
	fmt.Fprint(p.writer, options.String())
	return nil
}
