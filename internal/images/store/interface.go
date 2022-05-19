package store

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// GraphTemplateStorer is the interface for the graph template store
type GraphTemplateStorer interface{}

// ImageRenderer is the interface for the image renderer
type ImageRenderer interface {
	Render(name, version string, image *image.Image) (*image.Image, error)
}
