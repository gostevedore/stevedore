package command

import (
	"context"

	"github.com/gostevedore/stevedore/internal/driver"
)

// BuildDriverer interface defines which methods are used to build a docker image
type BuildDriverer interface {
	Build(context.Context, *driver.BuildDriverOptions) error
}
