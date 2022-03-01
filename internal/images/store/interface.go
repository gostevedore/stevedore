package store

import (
	"github.com/gostevedore/stevedore/internal/images/image"
)

// GraphTemplateStorer is the interface for the graph template store
type GraphTemplateStorer interface{}

type ImageRenderer interface {
	Render(name, version string, image *image.Image) error
}
