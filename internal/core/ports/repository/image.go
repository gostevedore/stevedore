package repository

import (
	"context"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// BuildDriverer interface defines which methods are needed to build a docker image
type BuildDriverer interface {
	Build(context.Context, *image.Image, *image.BuildDriverOptions) error
}

// Promoter interface defines which methods are needed to promote a docker image
type Promoter interface {
	Promote(context.Context, *image.PromoteOptions) error
}

// Renderer interface defined which methods are needed to renderize an image
type Renderer interface {
	Render(name, version string, image *image.Image) (*image.Image, error)
}

// ImagesStorer interface defines which methods are needed to save and read images from an images store
type ImagesStorer interface {
	ImagesStorerWriter
	ImagesStorerReader
}

// ImagesStoreWriter interface defines which methods are needed to save images to an images store
type ImagesStorerWriter interface {
	Store(name string, version string, image *image.Image) error
}

// ImagesStorerReader interface defines which methods are needed to read images from an images store
type ImagesStorerReader interface {
	List() ([]*image.Image, error)
	FindByName(name string) ([]*image.Image, error)
	Find(name string, version string) ([]*image.Image, error)
	FindGuaranteed(imageName, imageVersion string) ([]*image.Image, error)
}

// ImagesSelector interface defines which methods are needed to select images from a list
type ImagesSelector interface {
	Select(images []*image.Image, operation string, item string) ([]*image.Image, error)
}

// ImagesOutputter interface defines which methods are needed to return images definitions
type ImagesOutputter interface {
	Output(list []*image.Image) error
}

// ImagesPlainPrinter is an interface for printing images content output.
type ImagesPlainPrinter interface {
	PrintTable(content [][]string) error
}
