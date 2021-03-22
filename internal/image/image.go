package image

import (
	"fmt"

	"github.com/gostevedore/stevedore/internal/ui/console"

	common "github.com/apenella/go-common-utils/data"
	errors "github.com/apenella/go-common-utils/error"
)

const (
	InlineBuilder = "<in-line>"
)

// Image is the domain definition of a docker image
type Image struct {
	Builder        interface{}            `yaml:"builder"`
	Children       map[string][]string    `yaml:"children"`
	Childs         map[string][]string    `yaml:"childs"`
	Name           string                 `yaml:"name"`
	Namespace      string                 `yaml:"namespace"`
	PersistentVars map[string]interface{} `yaml:"persistent_vars"`
	Registry       string                 `yaml:"registry"`
	Tags           []string               `yaml:"tags"`
	Type           string                 `yaml:"type"`
	Vars           map[string]interface{} `yaml:"vars"`
	Version        string                 `yaml:"version"`
	// Parents        map[string][]string    `yaml:"parents"`
}

// LoadImage
func LoadImage(file string) (*Image, error) {
	image := &Image{}
	err := common.LoadYAMLFile(file, image)
	if err != nil {
		return nil, errors.New("(images::LoadImage)", "Images file could not be load", err)
	}

	return image, nil
}

// Copy method return a copy of the instanced Image
func (i *Image) Copy() (*Image, error) {

	if i == nil {
		return nil, errors.New("(image::Image::Copy)", "Image is nil")
	}

	copiedImage := *i
	copiedImage.Tags = append([]string{}, i.Tags...)

	copiedImage.PersistentVars = map[string]interface{}{}
	for keyVar, keyValue := range i.PersistentVars {
		copiedImage.PersistentVars[keyVar] = keyValue
	}
	copiedImage.Vars = map[string]interface{}{}
	for keyVar, keyValue := range i.Vars {
		copiedImage.Vars[keyVar] = keyValue
	}
	copiedImage.Vars = map[string]interface{}{}
	for keyVar, keyValue := range i.Vars {
		copiedImage.Vars[keyVar] = keyValue
	}

	if i.Children != nil {
		copiedImage.Children = map[string][]string{}
		for keyVar, keyValue := range i.Children {
			keyValueCopy := append([]string{}, keyValue...)
			copiedImage.Children[keyVar] = keyValueCopy
		}
	}

	return &copiedImage, nil
}

func (i *Image) ToArray() ([]string, error) {

	if i == nil {
		return nil, errors.New("(image::Image::ToArray)", "Image is nil")
	}

	arrayImage := []string{}
	arrayImage = append(arrayImage, i.Name)
	arrayImage = append(arrayImage, i.Version)
	arrayImage = append(arrayImage, i.getBuilderType())
	arrayImage = append(arrayImage, i.Namespace)
	arrayImage = append(arrayImage, i.Registry)

	return arrayImage, nil
}

// getBuilderType return the name of the builder or <in-line> when the builder is defined on the image
func (i *Image) getBuilderType() string {
	switch i.Builder.(type) {
	case string:
		return i.Builder.(string)
	default:
		return InlineBuilder
	}

}

// CheckCompatibility checks that image compatibility
func (i *Image) CheckCompatibility() {

	if i.Type != "" {
		console.Warn(fmt.Sprintf("DEPRECATION: On '%s', 'type' attribute must be replaced by 'builder' before 0.11.0", i.Name))

		if i.Builder == "" {
			i.Builder = i.Type
		} else {
			console.Warn(fmt.Sprintf(" On '%s', 'builder' value will be used instead of 'type'", i.Name))
		}

	}

	if i.Childs != nil && len(i.Childs) > 0 {
		console.Warn(fmt.Sprintf("DEPRECATION: On '%s', 'childs' attribute must be replaced by 'children' before 0.11.0", i.Name))

		if i.Children != nil && len(i.Children) > 0 {
			console.Warn(fmt.Sprintf(" On '%s', 'children' value will be used instead of 'childs'", i.Name))
		} else {
			i.Children = i.Childs
		}
	}
}
