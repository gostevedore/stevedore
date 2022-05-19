package driver

import (
	"context"

	"github.com/gostevedore/stevedore/internal/images/image"
)

// BuildDriverer interface defines which methods are used to build a docker image
type BuildDriverer interface {
	//Run(context.Context) error
	Build(context.Context, *image.Image, *BuildDriverOptions) error
}
