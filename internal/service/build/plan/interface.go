package plan

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// ImagesStorer interfaces defines the storage of images
type ImagesStorer interface {
	List() ([]*image.Image, error)
	FindByName(name string) ([]*image.Image, error)
	Find(name string, version string) (*image.Image, error)
	FindGuaranteed(findName, findVersion, imageName, imageVersion string) (*image.Image, error)
}

// Planner interfaces defines the storage of images
type Planner interface {
	Plan(name string, versions []string) ([]*Step, error)
}
