package driver

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// DryRunDriver is a driver that just simulates the build process
type DryRunDriver struct {
	write io.Writer
}

// NewDryRunDriver creates a new DryRunDriver
func NewDryRunDriver(w io.Writer) *DryRunDriver {

	if w == nil {
		w = os.Stdout
	}

	return &DryRunDriver{
		write: w,
	}
}

// Build simulate a new image build
func (d *DryRunDriver) Build(ctx context.Context, i *image.Image, options *image.BuildDriverOptions) error {
	fmt.Fprintln(d.write, fmt.Sprintf("%+v", *options))
	return nil
}
