package images

import (
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	configimage "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/image"
)

// TemplatesStorer interface
type ImagesGraphTemplatesStorer interface {
	AddImage(name, version string, image *configimage.Image) error
	Iterate() <-chan graph.GraphNoder
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
