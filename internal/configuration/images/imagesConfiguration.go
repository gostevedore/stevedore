package images

import (
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	domainimage "github.com/gostevedore/stevedore/internal/images/image"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// const (
// 	ImageNodeNameSeparator = ":"
// )

// ImagesConfiguration
//
// Image structure
// image_tree:
// 	image_name:
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
	store         ImagesStorer

	// DEPRECATEDImagesTree is replaced by Images
	DEPRECATEDImagesTree map[string]map[string]*image.Image `yaml:"images_tree"`
	Images               map[string]map[string]*image.Image `yaml:"images"`
}

// NewImagesConfiguration method create a new ImagesConfiguration struct
func NewImagesConfiguration(fs afero.Fs, graph ImagesGraphTemplatesStorer, store ImagesStorer, compatibility Compatibilitier) *ImagesConfiguration {
	return &ImagesConfiguration{
		fs:            fs,
		compatibility: compatibility,
		graph:         graph,
		store:         store,

		DEPRECATEDImagesTree: make(map[string]map[string]*image.Image),
		Images:               make(map[string]map[string]*image.Image),
	}
}

// CheckCompatibility method ensures that ImagesConfiguration is compatible with current version
func (t *ImagesConfiguration) CheckCompatibility() error {

	if t.DEPRECATEDImagesTree != nil && len(t.DEPRECATEDImagesTree) > 0 {
		t.compatibility.AddDeprecated("'images_tree' is deprecated and will be removed on v0.12.0, please use 'images' instead")
	}

	return nil
}

// LoadImagesToStore method loads images defined on configuration to images store
func (t *ImagesConfiguration) LoadImagesToStore(path string) error {

	var err error
	errContext := "(images::LoadImagesToStore)"
	var nodeDomainImage, copyDomainImage *domainimage.Image
	_ = copyDomainImage

	err = t.LoadImagesConfiguration(path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	storedNodes := map[string]struct{}{}
	for node := range t.graph.Iterate() {

		// skip node if already stored
		_, stored := storedNodes[node.Name()]
		if stored {
			continue
		}

		name, version, err := graph.ParseNodeName(node)
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		nodeImage := node.Item().(*image.Image)
		nodeDomainImage, err = nodeImage.CreateDomainImage()
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		if node.Parents() == nil || len(node.Parents()) <= 0 {
			err = t.store.Store(name, version, nodeDomainImage)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
		} else {
			for _, parent := range node.Parents() {

				parentName, parentVersion, err := graph.ParseNodeName(parent.(graph.GraphNoder))
				parentDomainImage, err := t.store.Find(parentName, parentVersion)
				if err != nil {
					return errors.New(errContext, err.Error())
				}

				copyDomainImage, err = nodeDomainImage.Copy()
				if err != nil {
					return errors.New(errContext, err.Error())
				}

				copyDomainImage.Options(
					domainimage.WithParent(parentDomainImage),
				)

				parentDomainImage.AddChild(copyDomainImage)

				err = t.store.Store(name, version, copyDomainImage)
				if err != nil {
					return errors.New(errContext, err.Error())
				}

				storedNodes[node.Name()] = struct{}{}
			}
		}
	}

	return nil
}

// LoadImagesConfiguration method generate and return an ImagesConfiguration struct from a file
func (t *ImagesConfiguration) LoadImagesConfiguration(path string) error {

	var err error
	var isDir bool

	errContext := "(images::LoadImagesConfiguration)"

	isDir, err = afero.IsDir(t.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if isDir {
		return t.LoadImagesConfigurationFromDir(path)
	} else {
		return t.LoadImagesConfigurationFromFile(path)
	}
}

// LoadImagesConfigurationFromDir loads images tree from all files on directory
func (t *ImagesConfiguration) LoadImagesConfigurationFromDir(dir string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(images::LoadImagesConfigurationFromDir)"

	yamlFiles, err := afero.Glob(t.fs, dir+"/*.yaml")
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	ymlFiles, err := afero.Glob(t.fs, dir+"/*.yml")
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	files := append(yamlFiles, ymlFiles...)

	loadImagesConfigurationFromFile := func(path string) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			err = t.LoadImagesConfigurationFromFile(path)
			t.wg.Done()
		}()

		return func() error {
			<-c
			return err
		}
	}

	for _, file := range files {
		t.wg.Add(1)
		f := loadImagesConfigurationFromFile(file)
		errFuncs = append(errFuncs, f)
	}

	t.wg.Wait()

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
func (t *ImagesConfiguration) LoadImagesConfigurationFromFile(path string) error {

	var err error
	var fileData []byte

	errContext := "(images::LoadImagesConfigurationFromFile)"

	if t == nil {
		return errors.New(errContext, "Builders is nil")
	}

	imageTreeAux := NewImagesConfiguration(t.fs, t.graph, t.store, t.compatibility)

	fileData, err = afero.ReadFile(t.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	err = yaml.Unmarshal(fileData, imageTreeAux)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error loading images tree from file '%s'\nfound:\n%s", path, string(fileData)), err)
	}

	err = imageTreeAux.CheckCompatibility()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	for name, images := range imageTreeAux.Images {
		if !isAValidName(name) {
			return errors.New(errContext, fmt.Sprintf("Found an invalid image name '%s' defined in file '%s'", name, path))
		}

		for version, image := range images {
			if !isAValidVersion(version) {
				return errors.New(errContext, fmt.Sprintf("Found an invalid image version '%s' defined in file '%s'", version, path))
			}

			err = t.graph.AddImage(name, version, image)
			// err = t.AddImage(name, version, image)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
		}
	}

	// TO BE REMOVE on v0.12: is kept just for compatibility concerns
	for name, images := range imageTreeAux.DEPRECATEDImagesTree {
		if !isAValidName(name) {
			return errors.New(errContext, fmt.Sprintf("Found an invalid image name '%s' defined in file '%s'", name, path))
		}

		for version, image := range images {
			if !isAValidVersion(version) {
				return errors.New(errContext, fmt.Sprintf("Found an invalid image version '%s' defined in file '%s'", version, path))
			}

			err := t.graph.AddImage(name, version, image)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
		}
	}

	return nil
}

// isValidName method checks if a string is a valid image name
func isAValidName(name string) bool {
	if strings.IndexRune(name, ':') != -1 {
		return false
	}
	return true
}

// isValidVersion method checks if a string is a valid image version
func isAValidVersion(name string) bool {
	if strings.IndexRune(name, ':') != -1 {
		return false
	}
	return true
}
