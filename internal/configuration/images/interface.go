package images

import (
	"github.com/gostevedore/stevedore/internal/configuration/images/graph"
	configimage "github.com/gostevedore/stevedore/internal/configuration/images/image"
	domainimage "github.com/gostevedore/stevedore/internal/core/domain/image"
)

// TemplatesStorer interface
type ImagesGraphTemplatesStorer interface {
	AddImage(name, version string, image *configimage.Image) error
	Iterate() <-chan graph.GraphNoder
}

// ImagesStorer interfaces defines the storage of images
type ImagesStorer interface {
	Store(name string, version string, parent *domainimage.Image) error
	Find(name string, version string) (*domainimage.Image, error)
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
