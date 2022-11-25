package image

import (
	errors "github.com/apenella/go-common-utils/error"
	domainimage "github.com/gostevedore/stevedore/internal/core/domain/image"
	"gopkg.in/yaml.v3"
)

const (
	InlineBuilder = "<in-line>"
)

// Image is the domain definition of a docker image
type Image struct {
	Builder           interface{}            `yaml:"builder"`
	Children          map[string][]string    `yaml:"children"`
	Labels            map[string]string      `yaml:"labels"`
	Name              string                 `yaml:"name"`
	Parents           map[string][]string    `yaml:"parents"`
	PersistentLabels  map[string]string      `yaml:"persistent_labels"`
	PersistentVars    map[string]interface{} `yaml:"persistent_vars"`
	RegistryHost      string                 `yaml:"registry"`
	RegistryNamespace string                 `yaml:"namespace"`
	Tags              []string               `yaml:"tags"`
	Vars              map[string]interface{} `yaml:"vars"`
	Version           string                 `yaml:"version"`
}

// Copy method return a copy of the instanced Image
func (i *Image) Copy() (*Image, error) {

	if i == nil {
		return nil, errors.New("(image::Image::Copy)", "Image is nil")
	}

	copiedImage := *i

	if i.Children != nil {
		copiedImage.Children = map[string][]string{}
		for keyVar, keyValue := range i.Children {
			copiedImage.Children[keyVar] = append([]string{}, keyValue...)
		}
	}

	copiedImage.PersistentLabels = map[string]string{}
	for keyVar, keyValue := range i.PersistentLabels {
		copiedImage.PersistentLabels[keyVar] = keyValue
	}

	copiedImage.PersistentVars = map[string]interface{}{}
	for keyVar, keyValue := range i.PersistentVars {
		copiedImage.PersistentVars[keyVar] = keyValue
	}

	copiedImage.Labels = map[string]string{}
	for keyVar, keyValue := range i.Labels {
		copiedImage.Labels[keyVar] = keyValue
	}

	if i.Parents != nil {
		copiedImage.Parents = map[string][]string{}
		for keyParent, keyValue := range i.Parents {
			copiedImage.Parents[keyParent] = append([]string{}, keyValue...)
		}
	}

	copiedImage.Tags = append([]string{}, i.Tags...)

	copiedImage.Vars = map[string]interface{}{}
	for keyVar, keyValue := range i.Vars {
		copiedImage.Vars[keyVar] = keyValue
	}

	return &copiedImage, nil
}

// CreateDomainImage creates a domain image from the image
func (i *Image) CreateDomainImage() (*domainimage.Image, error) {

	// NewImage from domain could not be used because in that stage some values are still a template
	image := &domainimage.Image{
		Name:              i.Name,
		Version:           i.Version,
		RegistryHost:      i.RegistryHost,
		RegistryNamespace: i.RegistryNamespace,
	}

	image.Options(
		domainimage.WithBuilder(i.Builder),
		domainimage.WithPersistentLabels(i.PersistentLabels),
		domainimage.WithPersistentVars(i.PersistentVars),
		domainimage.WithLabels(i.Labels),
		domainimage.WithTags(i.Tags...),
		domainimage.WithVars(i.Vars),
	)

	return image, nil
}

// CheckCompatibility checks that image compatibility
func (i *Image) CheckCompatibility(compabilitiy Compatibilitier) {}

// YAMLMarshal marshals the image to YAML
func (i *Image) YAMLMarshal() ([]byte, error) {
	errContext := "(image::YAMLMarshal)"

	marshaled, err := yaml.Marshal(i)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return marshaled, nil
}

// YAMLUnmarshal unmarshals the image from a YAML string
func (i *Image) YAMLUnmarshal(in []byte) error {
	errContext := "(image::YAMLUnmarshal)"

	err := yaml.Unmarshal(in, i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
