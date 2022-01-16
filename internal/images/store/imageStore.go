package store

import (
	"github.com/gostevedore/stevedore/internal/images/image"
)

// ImageStore is a store for images
type ImageStore struct {
	//tree
	//index
}

// NewImageStore returns a new instance of the ImageStore
func NewImageStore() *ImageStore {
	return &ImageStore{}
}

// AddImage adds an image to the store
func (m *ImageStore) AddImage(image *image.Image) error {
	return nil
}

// All returns all the images asociated to the image name
func (m *ImageStore) All(name string) ([]*image.Image, error) {
	return nil, nil
}

// Find returns the image associated to the image name and version
func (m *ImageStore) Find(name string, version string) (*image.Image, error) {
	return nil, nil
}
