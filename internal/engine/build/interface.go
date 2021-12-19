package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/image"
)

// BuildDriverer interface defines which methods are used to build a docker image
type BuildDriverer interface {
	Build(context.Context, *driver.BuildDriverOptions) error
}

// type Imager interface {
// 	GetItem() (*image.Image, error)
// 	GetParent() (Imager, error)
// 	GetChildren() ([]Imager, error)
// }

// ImagesStorer interfaces defines the storage of images
type ImagesStorer interface {
	All(string) ([]*image.Image, error)
	Find(string, string) (*image.Image, error)
}

// BuildersStorer interface defines the storage of builders
type BuildersStorer interface {
	Find(name string) (*builder.Builder, error)
}

type BuildCommandFactorier interface {
	New(BuildDriverer, *driver.BuildDriverOptions) BuildCommander
}

// BuildCommander interface defines the command to build a docker image
type BuildCommander interface {
	Execute(context.Context) error
}

type JobFactorier interface {
	New(BuildCommander) Jobber
}

// Jobber interface defines the job to build a docker image
type Jobber interface {
	Run(context.Context)
	Done() <-chan struct{}
	Err() <-chan error
	Close()
}

// Dispatcher is a dispatcher to build docker images
type Dispatcher interface {
	Enqueue(Jobber)
}
