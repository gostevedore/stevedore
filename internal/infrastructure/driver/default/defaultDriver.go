package driver

import (
	"context"
	"io"
	"os"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/dryrun"
)

// DefaultDriver is a driver that just simulates the build process
type DefaultDriver struct {
	write io.Writer
}

// NewDefaultDriver creates a new DefaultDriver
func NewDefaultDriver(w io.Writer) *DefaultDriver {

	if w == nil {
		w = os.Stdout
	}

	return &DefaultDriver{
		write: w,
	}
}

// Build simulate a new image build
func (d *DefaultDriver) Build(ctx context.Context, i *image.Image, options *image.BuildDriverOptions) error {
	driver := dryrun.NewDryRunDriver(d.write)

	return driver.Build(ctx, i, options)
}
