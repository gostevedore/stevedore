package store

import (
	"github.com/gostevedore/stevedore/internal/images/image"
)

// GraphTemplateStorer is the interface for the graph template store
type GraphTemplateStorer interface{}

type ImageSerializer interface {
	YAMLMarshal() ([]byte, error)
	YAMLUnmarshal([]byte) error
}

type ImageRenderer interface {
	Render(name, version string, parent *image.Image, image ImageSerializer) error
}
