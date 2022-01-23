package store

import (
	"github.com/gostevedore/stevedore/internal/images/image"
)

// ImageStore is a store for images
type ImageStore struct {
	// graphTemplate
	template GraphTemplateStorer
	//tree
	//index
}

// NewImageStore returns a new instance of the ImageStore
func NewImageStore(template GraphTemplateStorer) *ImageStore {
	return &ImageStore{
		template: template,
	}
}

// AddImage adds an image to the store
func (s *ImageStore) AddImage(image *image.Image) error {
	return nil
}

// All returns all the images asociated to the image name
func (s *ImageStore) All(name string) ([]*image.Image, error) {
	return nil, nil
}

// Find returns the image associated to the image name and version
func (s *ImageStore) Find(name string, version string) (*image.Image, error) {
	return nil, nil
}
