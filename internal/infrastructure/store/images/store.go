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

// Store is a store for images
type Store struct {

	// imageNameDefinitionVersionList is the list of versions defined for an image name
	imageNameDefinitionVersionList map[string]map[string]struct{}
	// imageNameRenderedVersionsList is the list of rendered versions for an image name and version
	imageNameVersionRenderedVersionsList map[string]map[string]map[string]struct{}
	// imagesIndex images referenced by its name and rendered version
	imagesIndex map[string]map[string]*image.Image

	imageWildcardIndex map[string]*image.Image
	render             repository.Renderer
	store              []*image.Image

	mutex sync.Mutex
}

// NewStore returns a new instance of the Store
func NewStore(render repository.Renderer) *Store {
	return &Store{
		imageNameDefinitionVersionList:       make(map[string]map[string]struct{}),
		imageNameVersionRenderedVersionsList: make(map[string]map[string]map[string]struct{}),
		imagesIndex:                          make(map[string]map[string]*image.Image),
		imageWildcardIndex:                   make(map[string]*image.Image),

		render: render,
		store:  []*image.Image{},
	}
}

// Store adds an image to the store
func (s *Store) Store(name string, version string, i *image.Image) error {
	var err error
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

	if version == image.ImageWildcardVersionSymbol {
		err = s.storeWildcardImage(name, i)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		return nil
	}

	// Commented because store does not render images anymore it is done by imagesConfiguration
	//
	// // render the image
	// renderedImage, err = s.render.Render(name, version, i)
	// if err != nil {
	// 	return errors.New(errContext, "", err)
	// }

	err = s.addImageNameDefinitionVersionList(name, version)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = s.addImageNameVersionRenderedVersionsList(name, version, i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = s.addImageToIndex(name, i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	s.store = append(s.store, i)

	return nil
}

func (s *Store) addImageNameDefinitionVersionList(name, version string) error {

	errContext := "(store::images::Store::addImageNameDefinitionVersionList)"

	if name == "" {
		return errors.New(errContext, "Image name must be provided to add to add a definition version into list")
	}

	if version == "" {
		return errors.New(errContext, "Image version must be provided to add to add a definition version into list")
	}

	var exist bool
	if s.imageNameDefinitionVersionList == nil {
		s.imageNameDefinitionVersionList = make(map[string]map[string]struct{})
	}

	_, exist = s.imageNameDefinitionVersionList[name]
	if !exist {
		s.imageNameDefinitionVersionList[name] = map[string]struct{}{
			version: {},
		}
	} else {
		s.imageNameDefinitionVersionList[name][version] = struct{}{}
	}

	return nil
}

func (s *Store) addImageNameVersionRenderedVersionsList(name, version string, i *image.Image) error {
	var exist bool

	errContext := "(store::images::Store::addImageNameVersionRenderedVersionsList)"

	if name == "" {
		return errors.New(errContext, "Image name must be provided to add to add a rendered version into list")
	}

	if version == "" {
		return errors.New(errContext, "Image version must be provided to add to add a rendered version into list")
	}

	if i == nil {
		return errors.New(errContext, "Image must be provided to add to add a rendered version into list")
	}

	if s.imageNameVersionRenderedVersionsList == nil {
		s.imageNameVersionRenderedVersionsList = make(map[string]map[string]map[string]struct{})
	}

	_, exist = s.imageNameVersionRenderedVersionsList[name]
	if !exist {
		s.imageNameVersionRenderedVersionsList[name] = map[string]map[string]struct{}{
			version: {
				i.Version: struct{}{},
			},
		}
	} else {
		_, exist = s.imageNameVersionRenderedVersionsList[name][version]
		if !exist {
			s.imageNameVersionRenderedVersionsList[name][version] = map[string]struct{}{
				i.Version: {},
			}
		} else {

			s.imageNameVersionRenderedVersionsList[name][version][i.Version] = struct{}{}
		}
	}

	return nil
}

func (s *Store) addImageToIndex(name string, i *image.Image) error {

	var exist bool
	errContext := "(store::images::Store::addImageToIndex)"

	if name == "" {
		return errors.New(errContext, "Image name must be provided to add to add image to index")
	}

	if i == nil {
		return errors.New(errContext, "Image must be provided to add to add to add image to index")
	}

	if s.imagesIndex == nil {
		s.imagesIndex = make(map[string]map[string]*image.Image)
	}

	_, exist = s.imagesIndex[name]
	if !exist {
		s.imagesIndex[name] = map[string]*image.Image{
			i.Version: i,
		}
	} else {
		s.imagesIndex[name][i.Version] = i
	}

	for _, tag := range i.Tags {
		s.imagesIndex[name][tag] = i
	}

	return nil
}

// wildcard images are those images that have * on its version. Wildcard images are used to generate a default image definition, and accepts any version value
//
// storeWildcardImage stores the image in the store
func (s *Store) storeWildcardImage(name string, i *image.Image) error {

	errContext := "(store::images::Store::storeWildcardImage)"

	if name == "" {
		return errors.New(errContext, "Image name must be provided to store a wildcard image")
	}

	if i == nil {
		return errors.New(errContext, fmt.Sprintf("Image must be provided to store '%s' wildcard image", name))
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

	errContext := "(store::images::Store::List)"

	if s.store == nil {
		return nil, errors.New(errContext, "To list images, store must be initialized")
	}

	sort.Sort(images.SortedImages(s.store))
	return s.store, nil
}

// FindByName returns all the images associated to the image name
func (s *Store) FindByName(name string) ([]*image.Image, error) {

	errContext := "(store::images::Store::FindByName)"

	if s.imageNameDefinitionVersionList == nil || s.imageNameVersionRenderedVersionsList == nil || s.imagesIndex == nil {
		return nil, errors.New(errContext, "To find images by name into images store, list structures must be initialized")
	}

	list := []*image.Image{}

	//	definitioVersionList,_ := s.imageNameDefinitionVersionList[name]
	for definitionVersion := range s.imageNameDefinitionVersionList[name] {
		for renderedVersion := range s.imageNameVersionRenderedVersionsList[name][definitionVersion] {
			i := s.imagesIndex[name][renderedVersion]
			list = append(list, i)
		}
	}

	sort.Sort(images.SortedImages(list))

	return list, nil
}

// Find returns the image associated to the image name and version
func (s *Store) Find(name string, version string) ([]*image.Image, error) {

	var i *image.Image
	var exist bool

	errContext := "(store::images::Store::Find)"
	_ = errContext

	if s.imageNameVersionRenderedVersionsList == nil || s.imagesIndex == nil || s.imageWildcardIndex == nil {
		return nil, errors.New(errContext, "To find images into images store, list structures must be initialized")
	}

	list := []*image.Image{}

	//  return the image associated to the image name and version
	if version == image.ImageWildcardVersionSymbol {
		i, _ = s.imageWildcardIndex[name]
		return append(list, i), nil
	}

	i, exist = s.imagesIndex[name][version]
	if exist {
		return append(list, i), nil
	}

	renderedVersionList, _ := s.imageNameVersionRenderedVersionsList[name][version]
	for renderedVersion := range renderedVersionList {
		i := s.imagesIndex[name][renderedVersion]
		list = append(list, i)
	}

	sort.Sort(images.SortedImages(list))

	return list, nil
}

// FindGuaranteed returns the image associated to the image name and version. In case of a wildcard image, it generates the image. Otherwise, it returns a nil image and an error
func (s *Store) FindGuaranteed(findName, findVersion, imageName, imageVersion string) ([]*image.Image, error) {

	var err error
	errContext := "(store::images::Store::FindGuaranteed)"

	var list []*image.Image
	var image, imageWildcard *image.Image

	list, err = s.Find(findName, findVersion)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if len(list) > 0 {
		return list, nil
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

	sort.Sort(images.SortedImages(list))

	return append(list, image), nil
}

func (s *Store) FindWildcardImage(name string) (*image.Image, error) {

	errContext := "(store::images::Store::FindWildcardImage)"

	if s.imageWildcardIndex == nil {
		return nil, errors.New(errContext, "To find a wildcard image, Wildcard index must be initialized")
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
	errContext := "(store::images::Store::GenerateImageFromWildcard)"

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
