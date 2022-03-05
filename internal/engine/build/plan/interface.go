package plan

import (
	"github.com/gostevedore/stevedore/internal/images/image"
)

// ImagesStorer interfaces defines the storage of images
type ImagesStorer interface {
	List() ([]*image.Image, error)
	FindByName(name string) ([]*image.Image, error)
	Find(string, string) (*image.Image, error)
}
