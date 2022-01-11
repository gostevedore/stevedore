package store

import (
	"github.com/gostevedore/stevedore/internal/image"
)

// ImageStore is a store for images
type ImageStore struct{}

// NewImageStore returns a new instance of the ImageStore
func NewImageStore() *ImageStore {
	return &ImageStore{}
}

// All returns all the images asociated to the image name
func (m *ImageStore) All(name string) ([]*image.Image, error) {
	return nil, nil
}

// Find returns the image associated to the image name and version
func (m *ImageStore) Find(name string, version string) (*image.Image, error) {
	return nil, nil
}
