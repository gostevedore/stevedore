package images

import (
	"fmt"
	"sort"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/filters/images"
)

const (
	// ImageWildcardVersionSymbol is the wildcard version
	ImageWildcardVersionSymbol = "*"
)

// Store is a store for images
type Store struct {
	render                repository.Renderer
	store                 []*image.Image
	imageNameVersionIndex map[string]map[string]*image.Image
	imageWildcardIndex    map[string]*image.Image

	mutex sync.Mutex
}

// NewStore returns a new instance of the Store
func NewStore(render repository.Renderer) *Store {
	return &Store{
		render:                render,
		store:                 []*image.Image{},
		imageNameVersionIndex: make(map[string]map[string]*image.Image),
	}
}

// Store adds an image to the store
func (s *Store) Store(name string, version string, i *image.Image) error {
	var err error
	var renderedImage *image.Image
	errContext := "(store::Store)"

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
			return errors.New(errContext, "", err)
		}

		return nil
	}

	// render the image
	renderedImage, err = s.render.Render(name, version, i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	// store the image
	err = s.storeImage(name, version, renderedImage)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if renderedImage.Version != version {
		err = s.storeImage(name, renderedImage.Version, renderedImage)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	}

	for _, tag := range renderedImage.Tags {

		if version != tag && renderedImage.Version != tag {
			err = s.storeImage(name, tag, renderedImage)
			if err != nil {
				return errors.New(errContext, "", err)
			}
		}
	}

	s.store = append(s.store, renderedImage)

	return nil
}

// storeImage stores the image in the store
func (s *Store) storeImage(name string, version string, i *image.Image) error {

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
func (s *Store) storeWildcardImage(name string, i *image.Image) error {

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
func (s *Store) List() ([]*image.Image, error) {

	errContext := "(store::List)"

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	sort.Sort(images.SortedImages(s.store))
	return s.store, nil
}

// FindByName returns all the images asociated to the image name
func (s *Store) FindByName(name string) ([]*image.Image, error) {

	errContext := "(store::FindByName)"
	list := []*image.Image{}

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	listOfVersion, _ := s.imageNameVersionIndex[name]
	for _, i := range listOfVersion {
		list = append(list, i)
	}

	sort.Sort(images.SortedImages(list))
	return list, nil
}

// Find returns the image associated to the image name and version
func (s *Store) Find(name string, version string) (*image.Image, error) {
	errContext := "(store::Find)"

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	//  return the image associated to the image name and version
	if version == ImageWildcardVersionSymbol {
		i, _ := s.imageWildcardIndex[name]
		return i, nil
	}

	i, exist := s.imageNameVersionIndex[name][version]
	if !exist {
		return nil, nil
	}

	return i, nil
}

// FindGuaranteed returns the image associated to the image name and version. In case of a wildcard image, it generates the image. Otherwise, it returns a nil image and an error
func (s *Store) FindGuaranteed(findName, findVersion, imageName, imageVersion string) (*image.Image, error) {

	var err error
	errContext := "(store::FindGuaranteed)"
	var image, imageWildcard *image.Image

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	image, err = s.Find(findName, findVersion)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if image != nil {
		return image, nil
	}

	imageWildcard, err = s.FindWildcardImage(findName)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if imageWildcard == nil {
		return nil, errors.New(errContext, fmt.Sprintf("Image '%s:%s' does not exist on the store", findName, findVersion))
	}

	image, err = s.GenerateImageFromWildcard(imageWildcard, imageName, imageVersion)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return image, nil
}

func (s *Store) FindWildcardImage(name string) (*image.Image, error) {

	errContext := "(store::FindWildcardImage)"

	if s.store == nil {
		return nil, errors.New(errContext, "Store has not been initialized")
	}

	if s.imageWildcardIndex == nil {
		return nil, errors.New(errContext, "Wildcard index has not been initialized")
	}

	i, exist := s.imageWildcardIndex[name]
	if !exist {
		return nil, nil
	}

	return i, nil
}

func (s *Store) GenerateImageFromWildcard(i *image.Image, name string, version string) (*image.Image, error) {

	var err error
	var parent, parentWildcard, renderedImage, imageToRender *image.Image
	errContext := "(store::GenerateImageFromWildcard)"

	if i == nil {
		return nil, errors.New(errContext, "Provided wildcard image is nil")
	}

	imageToRender, err = i.Copy()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	parent = i.Parent

	// ensure that parent is properly rended when it is also a wildcard image
	if parent != nil {
		parentWildcard, err = s.FindWildcardImage(parent.Name)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		if parentWildcard != nil {
			parent, err = s.GenerateImageFromWildcard(parentWildcard, parent.Name, version)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}
		}
		imageToRender.Options(image.WithParent(parent))
	}

	renderedImage, err = s.render.Render(name, version, imageToRender)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return renderedImage, nil
}
