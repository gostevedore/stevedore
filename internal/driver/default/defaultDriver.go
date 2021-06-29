package defaultdriver

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

const (
	//	BuilderName "default"
	DriverName = "default"
)

type DefaultDriver struct {
	Writer  io.Writer
	options *types.BuildOptions
}

func (b *DefaultDriver) Run(ctx context.Context) error {

	if b.Writer == nil {
		b.Writer = os.Stdout
	}

	fmt.Fprintln(b.Writer, fmt.Sprintf("%+v", *b.options))
	return nil
}

func NewDefaultDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {

	b := &DefaultDriver{
		options: o,
		Writer:  console.GetConsole(),
	}

	return b, nil
}
