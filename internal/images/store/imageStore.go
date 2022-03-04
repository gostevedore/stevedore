package store

import (
	"fmt"
	"sort"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/images/store/filter"
)

const (
	// ImageWildcardVersionSymbol is the wildcard version
	ImageWildcardVersionSymbol = "*"
)

// ImageStore is a store for images
type ImageStore struct {
	render ImageRenderer
	//tree
	//index
	store                 []*image.Image
	imageNameVersionIndex map[string]map[string]*image.Image
	imageWildcardIndex    map[string]*image.Image

	mutex sync.Mutex
}

// NewImageStore returns a new instance of the ImageStore
func NewImageStore(render ImageRenderer) *ImageStore {
	return &ImageStore{
		render:                render,
		store:                 []*image.Image{},
		imageNameVersionIndex: make(map[string]map[string]*image.Image),
	}
}

// AddImage adds an image to the store
func (s *ImageStore) AddImage(name string, version string, i *image.Image) error {
	var err error
	var renderedImage *image.Image
	errContext := "(store::AddImage)"

	if s.render == nil {
		return errors.New(errContext, "To add an image to the store an image render is required")
	}

	if name == "" {
		return errors.New(errContext, "To add an image to the store a name is required")
	}

	if version == "" {
		return errors.New(errContext, "To add an image to the store a version is required")
	}

	if i == nil {
		return errors.New(errContext, "To add an image to the store an image is required")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if version == ImageWildcardVersionSymbol {
		err = s.storeWildcardImage(name, i)
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		return nil
	}

	// render the image
	renderedImage, err = s.render.Render(name, version, i)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	// store the image
	err = s.storeImage(name, version, renderedImage)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if renderedImage.Version != version {
		err = s.storeImage(name, renderedImage.Version, renderedImage)
		if err != nil {
			return errors.New(errContext, err.Error())
		}
	}

	for _, tag := range i.Tags {
		err = s.storeImage(name, tag, renderedImage)
		if err != nil {
			return errors.New(errContext, err.Error())
		}
	}

	s.store = append(s.store, renderedImage)

	return nil
}

// storeImage stores the image in the store
func (s *ImageStore) storeImage(name string, version string, i *image.Image) error {

	errContext := "(store::storeImage)"

	if i == nil {
		return errors.New(errContext, fmt.Sprintf("Provided image for '%s:%s' is nil", name, version))
	}

	if s.imageNameVersionIndex == nil {
		s.imageNameVersionIndex = make(map[string]map[string]*image.Image)
	}

	if s.imageNameVersionIndex[name] == nil {
		s.imageNameVersionIndex[name] = make(map[string]*image.Image)
	}

	_, exist := s.imageNameVersionIndex[name][version]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Image '%s:%s' already exists", name, version))
	}

	s.imageNameVersionIndex[name][version] = i

	return nil
}

//  wildcard images are those images that have * on its version. Wildcard images are used to generate a default image definition, and accepts any version value
// storeWildcardImage stores the image in the store
func (s *ImageStore) storeWildcardImage(name string, i *image.Image) error {

	errContext := "(store::storeWildcardImage)"

	if i == nil {
		return errors.New(errContext, fmt.Sprintf("Provided wildcard image for '%s' is nil", name))
	}

	if s.imageWildcardIndex == nil {
		s.imageWildcardIndex = make(map[string]*image.Image)
	}

	_, exist := s.imageWildcardIndex[name]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Image '%s' already exists on wildcard images index", name))
	}

	s.imageWildcardIndex[name] = i

	return nil
}

// List returns a sorted list of all images
func (s *ImageStore) List() ([]*image.Image, error) {

	errContext := "(store::List)"

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	sort.Sort(filter.SortedImages(s.store))
	return s.store, nil
}

// FindByName returns all the images asociated to the image name
func (s *ImageStore) FindByName(name string) ([]*image.Image, error) {

	errContext := "(store::FindByName)"
	list := []*image.Image{}

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	listOfVersion, _ := s.imageNameVersionIndex[name]
	for _, i := range listOfVersion {
		list = append(list, i)
	}

	sort.Sort(filter.SortedImages(list))
	return list, nil
}

// Find returns the image associated to the image name and version
func (s *ImageStore) Find(name string, version string) (*image.Image, error) {
	var err error
	errContext := "(store::Find)"

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	// version is *
	//  return the image associated to the image name and version
	if version == ImageWildcardVersionSymbol {
		i, _ := s.imageWildcardIndex[name]
		return i, nil
	}

	// lookup names index
	// return image associated to the image name and version

	// lookup tags index
	//  return image associated to the image name and version
	i, exist := s.imageNameVersionIndex[name][version]
	if !exist {
		i, err = s.GenerateImageFromWildcard(name, version)
		if err != nil {
			return nil, errors.New(errContext, err.Error())
		}
	}

	// lookup wildcard
	//  generate the image based on the wildcard version and return it

	return i, nil
}

func (s *ImageStore) GenerateImageFromWildcard(name string, version string) (*image.Image, error) {

	var err error
	var parent, aux *image.Image
	errContext := "(store::GenerateImageFromWildcard)"

	i, exists := s.imageWildcardIndex[name]
	if !exists {
		return nil, nil
	}

	parent = i.Parent

	if i.Parent != nil {
		if i.Parent.Version == ImageWildcardVersionSymbol {
			// next level
			parent, err = s.GenerateImageFromWildcard(i.Parent.Name, version)
			if err != nil {
				return nil, errors.New(errContext, err.Error())
			}

			err = s.AddImage(i.Parent.Name, version, parent)
			if err != nil {
				return nil, errors.New(errContext, err.Error())
			}
		}
	}

	aux, err = i.Copy()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}
	aux.Options(image.WithParent(parent))

	aux, err = s.render.Render(name, version, aux)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return aux, nil
}
