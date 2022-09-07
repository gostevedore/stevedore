package repository

import (
	"context"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// BuildDriverer interface defines which methods are used to build a docker image
type BuildDriverer interface {
	Build(context.Context, *image.Image, *image.BuildDriverOptions) error
}

// Promoter
type Promoter interface {
	Promote(context.Context, *image.PromoteOptions) error
}

// Renderer is the interface for the image renderer
type Renderer interface {
	Render(name, version string, image *image.Image) (*image.Image, error)
}

// ImagesStorer interfaces defines the storage of images
type ImagesStorer interface {
	ImagesStorerWriter
	ImagesStorerReader
}

type ImagesStorerWriter interface {
	Store(name string, version string, image *image.Image) error
}

// ImagesStorer interfaces defines the storage of images
type ImagesStorerReader interface {
	List() ([]*image.Image, error)
	FindByName(name string) ([]*image.Image, error)
	Find(name string, version string) ([]*image.Image, error)
	FindGuaranteed(findName, findVersion, imageName, imageVersion string) ([]*image.Image, error)
}
