package images

import (
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	domainimage "github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/image"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// const (
// 	ImageNodeNameSeparator = ":"
// )

// ImagesConfiguration
//
// Image structure
// image_tree:
//
//	image_name:
//		image_tag1:
//			<Image>
//		image_tag2:
//			<Image>
type ImagesConfiguration struct {
	compatibility Compatibilitier
	graph         ImagesGraphTemplatesStorer
	fs            afero.Fs
	mutex         sync.RWMutex
	wg            sync.WaitGroup
	store         repository.ImagesStorer
	render        repository.Renderer

	// DEPRECATEDImagesTree is replaced by Images
	DEPRECATEDImagesTree map[string]map[string]*image.Image `yaml:"images_tree"`
	Images               map[string]map[string]*image.Image `yaml:"images"`
}

// NewImagesConfiguration method create a new ImagesConfiguration struct
func NewImagesConfiguration(fs afero.Fs, graph ImagesGraphTemplatesStorer, store repository.ImagesStorer, render repository.Renderer, compatibility Compatibilitier) *ImagesConfiguration {
	return &ImagesConfiguration{
		fs:            fs,
		compatibility: compatibility,
		graph:         graph,
		store:         store,
		render:        render,

		DEPRECATEDImagesTree: make(map[string]map[string]*image.Image),
		Images:               make(map[string]map[string]*image.Image),
	}
}

// CheckCompatibility method ensures that ImagesConfiguration is compatible with current version
func (c *ImagesConfiguration) CheckCompatibility() error {

	if c.DEPRECATEDImagesTree != nil && len(c.DEPRECATEDImagesTree) > 0 {
		c.compatibility.AddDeprecated("'images_tree' is deprecated and will be removed on v0.12.0, please use 'images' instead")
	}

	return nil
}

// LoadImagesToStore method loads images defined on configuration to images store
func (c *ImagesConfiguration) LoadImagesToStore(path string) error {

	var err error
	errContext := "(images::LoadImagesToStore)"

	err = c.LoadImagesConfiguration(path)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	storedNodes := map[string]struct{}{}
	pendingNodes := map[string]map[string]struct{}{}

	for node := range c.graph.Iterate() {

		// skip node when it is already stored
		_, stored := storedNodes[node.Name()]
		if stored {
			continue
		}

		err = c.storeNodeImages(node, storedNodes, pendingNodes)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	}

	if len(pendingNodes) != 0 {
		return errors.New(errContext, fmt.Sprintf("There are orphan references to images that have not been defined\n%+v", pendingNodes))
	}

	return nil
}

// storeImage stores image to images store
func (c *ImagesConfiguration) storeNodeImages(node graph.GraphNoder, storedNodes map[string]struct{}, pendingNodes map[string]map[string]struct{}) error {
	var err error
	var nodeDomainImage *domainimage.Image
	errContext := "(images::storeImage)"

	name, version, err := graph.ParseNodeName(node)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if node.Item() == nil {
		return errors.New(errContext, fmt.Sprintf("Definition for the image '%s' has not been found", node.Name()))
	}

	nodeImage := node.Item().(*image.Image)

	nodeDomainImage, err = nodeImage.CreateDomainImage()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if node.Parents() == nil || len(node.Parents()) <= 0 {

		imageToStore, err := c.renderImage(name, version, nodeDomainImage)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		err = c.store.Store(name, version, imageToStore)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	} else {
		for _, parent := range node.Parents() {
			// skip node when it is already stored
			_, stored := storedNodes[strings.Join([]string{node.Name(), parent.Name()}, ":")]
			if stored {
				continue
			}

			copyDomainImage, err := nodeDomainImage.Copy()
			if err != nil {
				return errors.New(errContext, "", err)
			}

			parentName, parentVersion, err := graph.ParseNodeName(parent.(graph.GraphNoder))
			parentDomainImageList, err := c.store.Find(parentName, parentVersion)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			// when then parent is not already created, the current node is skipped
			// Node is referred by all its parents then when the latest parent is created, the node will be created as well
			if len(parentDomainImageList) == 0 || parentDomainImageList == nil {
				pendingNodes[parentName] = map[string]struct{}{parentVersion: {}}
				continue
			}

			// it assign the first item as parent
			parentDomainImage := parentDomainImageList[0]
			copyDomainImage.Options(
				domainimage.WithParent(parentDomainImage),
			)

			err = c.propagatePersistentAttributes(copyDomainImage)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			imageToStore, err := c.renderImage(name, version, copyDomainImage)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			parentDomainImage.AddChild(imageToStore)

			err = c.store.Store(name, version, imageToStore)
			if err != nil {
				return errors.New(errContext, "", err)
			}

			storedNodes[strings.Join([]string{node.Name(), parent.Name()}, ":")] = struct{}{}
		}
	}

	pendingNodeName, exists := pendingNodes[name]
	if exists {
		_, exists := pendingNodeName[version]
		if exists {
			delete(pendingNodeName, version)
		}
		if len(pendingNodeName) == 0 {
			delete(pendingNodes, name)
		}
	}

	return nil
}

func (c *ImagesConfiguration) propagatePersistentAttributes(i *domainimage.Image) error {

	errContext := "(images::propagatePersistentAttributes)"
	if i == nil {
		return errors.New(errContext, "To propagate persistent attributs, an image must be provided")
	}

	if i.Parent == nil {
		return nil
	}

	if i.Parent.PersistentLabels != nil {
		for k, v := range i.Parent.PersistentLabels {
			i.PersistentLabels[k] = v
		}
	}

	if i.Parent.PersistentVars != nil {
		for k, v := range i.Parent.PersistentVars {
			i.PersistentVars[k] = v
		}
	}

	return nil
}

// renderImage return a renderized image base on the input image. The input image is not rendered when the version value is ImageWildcardSymbol
func (c *ImagesConfiguration) renderImage(name, version string, i *domainimage.Image) (*domainimage.Image, error) {

	var err error
	var renderedImage *domainimage.Image
	//var normalizedImage *domainimage.Image
	errContext := "(images::renderImage)"

	if version == domainimage.ImageWildcardVersionSymbol {
		err = i.Sanetize()
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		return i, nil
	}

	if name == "" {
		return nil, errors.New(errContext, "Image name must be provided to render an image")
	}

	if version == "" {
		return nil, errors.New(errContext, "Image version must be provided to render an image")
	}

	if c.render == nil {
		return nil, errors.New(errContext, fmt.Sprintf("Image '%s:%s' could not be rendered because renderer must by provided", name, version))
	}

	renderedImage, err = c.render.Render(name, version, i)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	// normalizedImage, err = domainimage.NewImage(renderedImage.Name, renderedImage.Version, renderedImage.RegistryHost, renderedImage.RegistryNamespace)
	// if err != nil {
	// 	return nil, errors.New(errContext, "", err)
	// }

	// renderedImage.Name = normalizedImage.Name
	// renderedImage.Version = normalizedImage.Version
	// renderedImage.RegistryHost = normalizedImage.RegistryHost
	// renderedImage.RegistryNamespace = normalizedImage.RegistryNamespace

	err = renderedImage.Sanetize()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return renderedImage, nil
}

// LoadImagesConfiguration method generate and return an ImagesConfiguration struct from a file
func (c *ImagesConfiguration) LoadImagesConfiguration(path string) error {

	var err error
	var isDir bool

	errContext := "(images::LoadImagesConfiguration)"

	isDir, err = afero.IsDir(c.fs, path)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if isDir {
		return c.LoadImagesConfigurationFromDir(path)
	} else {
		return c.LoadImagesConfigurationFromFile(path)
	}
}

// LoadImagesConfigurationFromDir loads images tree from all files on directory
func (c *ImagesConfiguration) LoadImagesConfigurationFromDir(dir string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(images::LoadImagesConfigurationFromDir)"

	yamlFiles, err := afero.Glob(c.fs, dir+"/*.yaml")
	if err != nil {
		return errors.New(errContext, "", err)
	}

	ymlFiles, err := afero.Glob(c.fs, dir+"/*.yml")
	if err != nil {
		return errors.New(errContext, "", err)
	}
	files := append(yamlFiles, ymlFiles...)

	loadImagesConfigurationFromFile := func(path string) func() error {
		var err error

		ch := make(chan struct{}, 1)
		go func() {
			defer close(ch)
			err = c.LoadImagesConfigurationFromFile(path)
			c.wg.Done()
		}()

		return func() error {
			<-ch
			return err
		}
	}

	for _, file := range files {
		c.wg.Add(1)
		f := loadImagesConfigurationFromFile(file)
		errFuncs = append(errFuncs, f)
	}

	c.wg.Wait()

	errMsg := ""
	for _, f := range errFuncs {
		err = f()
		if err != nil {
			errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errContext, errMsg)
	}

	return nil
}

// LoadImagesConfigurationFromFile loads images tree from file
func (c *ImagesConfiguration) LoadImagesConfigurationFromFile(path string) error {

	var err error
	var fileData []byte

	errContext := "(images::LoadImagesConfigurationFromFile)"

	// TODO check if it is required
	if c == nil {
		return errors.New(errContext, "ImagesConfiguration is nil")
	}

	imageTreeAux := NewImagesConfiguration(c.fs, c.graph, c.store, c.render, c.compatibility)

	fileData, err = afero.ReadFile(c.fs, path)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = yaml.Unmarshal(fileData, imageTreeAux)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error loading images tree from file '%s'\nfound:\n%s", path, string(fileData)), err)
	}

	err = imageTreeAux.CheckCompatibility()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	for name, images := range imageTreeAux.Images {

		if !isAValidName(name) {
			return errors.New(errContext, fmt.Sprintf("Found an invalid image name '%s' defined in file '%s'", name, path))
		}

		for version, imageDefinition := range images {
			if !isAValidVersion(version) {
				return errors.New(errContext, fmt.Sprintf("Found an invalid image version '%s' defined in file '%s'", version, path))
			}

			// When the imageDefinition is not found, is set an empty values. It allow you to define images from DockerHub without an explicit image definition
			if imageDefinition == nil {
				imageDefinition = &image.Image{}
			}

			if imageDefinition.Name == "" {
				imageDefinition.Name = name
			}

			if imageDefinition.Version == "" {
				imageDefinition.Version = version
			}

			err = c.graph.AddImage(name, version, imageDefinition)
			// err = t.AddImage(name, version, image)
			if err != nil {
				return errors.New(errContext, "", err)
			}
		}
	}

	// TO BE REMOVE on v0.12: is kept just for compatibility concerns
	for name, images := range imageTreeAux.DEPRECATEDImagesTree {

		if !isAValidName(name) {
			return errors.New(errContext, fmt.Sprintf("Found an invalid image name '%s' defined in file '%s'", name, path))
		}

		for version, imageDefinition := range images {
			if !isAValidVersion(version) {
				return errors.New(errContext, fmt.Sprintf("Found an invalid image version '%s' defined in file '%s'", version, path))
			}

			// When the imageDefinition is not found, is set an empty values. It allow you to define images from DockerHub without an explicit image definition
			if imageDefinition == nil {
				imageDefinition = &image.Image{}
			}

			if imageDefinition.Name == "" {
				imageDefinition.Name = name
			}

			if imageDefinition.Version == "" {
				imageDefinition.Version = version
			}

			err = c.graph.AddImage(name, version, imageDefinition)
			if err != nil {
				return errors.New(errContext, "", err)
			}
		}
	}

	return nil
}

// isValidName method checks if a string is a valid image name
func isAValidName(name string) bool {

	if name == "" {
		return false
	}

	if strings.IndexRune(name, ':') != -1 {
		return false
	}
	return true
}

// isValidVersion method checks if a string is a valid image version
func isAValidVersion(version string) bool {
	if strings.IndexRune(version, ':') != -1 {
		return false
	}
	return true
}
