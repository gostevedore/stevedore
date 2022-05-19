package driver

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/driver"
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
func (d *DefaultDriver) Build(ctx context.Context, i *image.Image, options *driver.BuildDriverOptions) error {
	fmt.Fprintln(d.write, fmt.Sprintf("%+v", *options))
	return nil
}

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/gostevedore/stevedore/internal/types"
// 	"github.com/gostevedore/stevedore/internal/ui/console"
// )

// const (
// 	//	BuilderName "default"
// 	DriverName = "default"
// )

// type DefaultDriver struct {
// 	Writer  io.Writer
// 	options *types.BuildOptions
// }

// func (b *DefaultDriver) Run(ctx context.Context) error {

// 	if b.Writer == nil {
// 		b.Writer = os.Stdout
// 	}

// 	fmt.Fprintln(b.Writer, fmt.Sprintf("%+v", *b.options))
// 	return nil
// }

// func NewDefaultDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {

// 	b := &DefaultDriver{
// 		options: o,
// 		Writer:  console.GetConsole(),
// 	}

// 	return b, nil
// }
