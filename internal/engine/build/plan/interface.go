package plan

import (
	"github.com/gostevedore/stevedore/internal/image"
)

// ImagesStorer interfaces defines the storage of images
type ImagesStorer interface {
	All(string) ([]*image.Image, error)
	Find(string, string) (*image.Image, error)
}
