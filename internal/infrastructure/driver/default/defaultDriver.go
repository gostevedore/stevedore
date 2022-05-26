package driver

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
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
	fmt.Fprintln(d.write, fmt.Sprintf("%+v", *options))
	return nil
}
