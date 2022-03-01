package store

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
)

const (
	// ImageWildcardVersion is the wildcard version
	ImageWildcardVersion = "*"
)

// ImageStore is a store for images
type ImageStore struct {
	render ImageRenderer
	//tree
	//index
	imageNameVersionIndex map[string]map[string]*image.Image
	imageTagsIndex        map[string]map[string]*image.Image
	imageWildcardIndex    map[string]*image.Image
}

// NewImageStore returns a new instance of the ImageStore
func NewImageStore(render ImageRenderer) *ImageStore {
	return &ImageStore{
		render:                render,
		imageNameVersionIndex: make(map[string]map[string]*image.Image),
	}
}

// AddImage adds an image to the store
func (s *ImageStore) AddImage(name string, version string, i *image.Image) error {
	var err error
	errContext := "(store::AddImage)"

	if version == ImageWildcardVersion {
		err = s.storeWildcardImage(name, i)
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		return nil
	}

	// render the image
	err = s.render.Render(name, version, i)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	// store the image
	err = s.storeImageNameVersion(name, version, i)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if i.Version != version {
		err = s.storeImageTags(name, version, i)
		if err != nil {
			return errors.New(errContext, err.Error())
		}
	}

	return nil
}

func (s *ImageStore) storeWildcardImage(name string, i *image.Image) error {

	errContext := "(store::storeWildcardImage)"

	if s.imageWildcardIndex == nil {
		s.imageWildcardIndex = make(map[string]*image.Image)
	}

	_, exist := s.imageWildcardIndex[name]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Image %s already exists on wildcard images index", name))
	}

	s.imageWildcardIndex[name] = i

	return nil

}

func (s *ImageStore) storeImageNameVersion(name string, version string, i *image.Image) error {

	errContext := "(store::storeImageNameVersion)"

	if s.imageNameVersionIndex == nil {
		s.imageNameVersionIndex = make(map[string]map[string]*image.Image)
	}

	if s.imageNameVersionIndex[name] == nil {
		s.imageNameVersionIndex[name] = make(map[string]*image.Image)
	}

	_, exist := s.imageNameVersionIndex[name][version]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Image %s:%s already exists", name, version))
	}

	s.imageNameVersionIndex[name][version] = i

	return nil
}

func (s *ImageStore) storeImageTags(name string, version string, i *image.Image) error {

	errContext := "(store::storeImageNameVersion)"

	if s.imageNameVersionIndex == nil {
		s.imageNameVersionIndex = make(map[string]map[string]*image.Image)
	}

	if s.imageNameVersionIndex[name] == nil {
		s.imageNameVersionIndex[name] = make(map[string]*image.Image)
	}

	_, exist := s.imageNameVersionIndex[name][version]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Image %s:%s already exists", name, version))
	}

	s.imageNameVersionIndex[name][version] = i

	return nil
}

// All returns all the images asociated to the image name
func (s *ImageStore) All(name string) ([]*image.Image, error) {

	// rerturn all images associated to an image name

	return nil, nil
}

// Find returns the image associated to the image name and version
func (s *ImageStore) Find(name string, version string) (*image.Image, error) {

	// version is *
	//  return the image associated to the image name and version

	// lookup names index
	// return image associated to the image name and version

	// lookup tags index
	//  return image associated to the image name and version

	// lookup wildcard
	//  generate the image based on the wildcard version and return it

	return nil, nil
}
