package image

import (
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/docker/distribution/reference"
	"gopkg.in/yaml.v2"
)

// Image defines the image on the system
type Image struct {
	// Builder is the builder to use to build the image
	Builder interface{} `yaml:"builder"`
	// Children list of children images
	Children []*Image `yaml:"-"`
	// Labels is a map of image labels
	Labels map[string]string `yaml:"labels"`
	// Name is the name of the image
	Name string `yaml:"name"`
	// PresistentVars is a map of persistent variables
	PersistentVars map[string]interface{} `yaml:"persistent_vars"`
	// RegistryHost is the host of the registry
	RegistryHost string `yaml:"registry_host"`
	// RegistryNamespace is the namespace of the registry
	RegistryNamespace string `yaml:"registry_namespace"`
	// Tags is a list of extra tags
	Tags []string `yaml:"tags"`
	// Vars is a map of variables
	Vars map[string]interface{} `yaml:"vars"`
	// Version is the version of the image
	Version string `yaml:"version"`
	// Parent is the parent image
	Parent *Image `yaml:"-"`

	addChildMutex sync.RWMutex
}

// OptionFunc is an option to pass to NewImage
type OptionFunc func(*Image)

// NewImage creates a new image
func NewImage(name, version, registryHost, registryNamesapace string, opt ...OptionFunc) (*Image, error) {

	errContext := "(image::NewImage)"

	// Image name normalization
	imageName := name
	if version != "" {
		imageName = fmt.Sprintf("%s:%s", imageName, version)
	}

	if registryNamesapace != "" {
		imageName = fmt.Sprintf("%s/%s", registryNamesapace, imageName)
	}

	if registryHost != "" {

		if strings.IndexRune(registryHost, '.') < 0 {
			return nil, errors.New(errContext, "Registry host name must by a FQDN")
		}

		imageName = fmt.Sprintf("%s/%s", registryHost, imageName)
	}

	image, err := Parse(imageName)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	image.Options(opt...)
	return image, nil
}

// WithBuilder sets the builder
func WithBuilder(builder interface{}) OptionFunc {
	return func(i *Image) {
		i.Builder = builder
	}
}

// WithChildren sets the children
func WithChildren(children ...*Image) OptionFunc {
	return func(i *Image) {
		i.Children = children
	}
}

// WithLabels sets the labels
func WithLabels(labels map[string]string) OptionFunc {
	return func(i *Image) {
		i.Labels = labels
	}
}

// WithParent sets the parent image
func WithParent(parent *Image) OptionFunc {
	return func(i *Image) {
		i.Parent = parent
	}
}

// WithPersistentVars sets the persistent variables
func WithPersistentVars(persistentVars map[string]interface{}) OptionFunc {
	return func(i *Image) {
		i.PersistentVars = persistentVars
	}
}

// WithRegistryHost sets the registry host
// func WithRegistryHost(registryHost string) OptionFunc {
// 	return func(i *Image) {
// 		i.RegistryHost = registryHost
// 	}
// }

// WithRegistryNamespace sets the registry namespace
// func WithRegistryNamespace(registryNamespace string) OptionFunc {
// 	return func(i *Image) {
// 		i.RegistryNamespace = registryNamespace
// 	}
// }

// WithTags sets the tags
func WithTags(tags ...string) OptionFunc {
	return func(i *Image) {
		i.Tags = tags
	}
}

// WithVars sets the variables
func WithVars(vars map[string]interface{}) OptionFunc {
	return func(i *Image) {
		i.Vars = vars
	}
}

// Parse generates an image skeleton using a docker images name passed as parameter
func Parse(name string) (*Image, error) {
	errContext := "(image::Parse)"

	referenceName, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return nil, errors.New(errContext, "Image could not be parsed", err)
	}
	// in case that no version is specified, it will be the tag it as 'latest'
	referenceName = reference.TagNameOnly(referenceName)

	return &Image{
		Name:              reference.Path(referenceName)[strings.LastIndex(reference.Path(referenceName), "/")+1:],
		Version:           referenceName.String()[strings.IndexRune(referenceName.String(), ':')+1:],
		RegistryHost:      reference.Domain(referenceName),
		RegistryNamespace: reference.Path(referenceName)[:strings.LastIndex(reference.Path(referenceName), "/")],
	}, nil
}

// Options returns the options of the image
func (i *Image) Options(o ...OptionFunc) {
	for _, opt := range o {
		opt(i)
	}
}

// AddChild adds a child image
func (i *Image) AddChild(child *Image) {
	i.addChildMutex.Lock()
	defer i.addChildMutex.Unlock()

	i.Children = append(i.Children, child)
}

// NormalizedNamed normalizes the image name
func (i *Image) DockerNormalizedNamed() (string, error) {
	errContext := "(image::DockerNormalizedNamed)"

	if i.Name == "" {
		return "", errors.New(errContext, "Image name is empty")
	}

	if i.Version == "" {
		return "", errors.New(errContext, "Image version is empty")
	}

	if i.RegistryHost == "" {
		return "", errors.New(errContext, "Registry host is empty")
	}

	if i.RegistryNamespace == "" {
		return "", errors.New(errContext, "Registry namespace is empty")
	}

	return fmt.Sprintf("%s/%s/%s:%s", i.RegistryHost, i.RegistryNamespace, i.Name, i.Version), nil
}

// Copy method return a copy of the instanced Image
func (i *Image) Copy() (*Image, error) {

	errContext := "(image::Copy)"

	copiedImage, err := NewImage(i.Name, i.Version, i.RegistryHost, i.RegistryNamespace)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}
	copiedImage.Tags = append([]string{}, i.Tags...)

	copiedImage.PersistentVars = map[string]interface{}{}
	for keyVar, keyValue := range i.PersistentVars {
		copiedImage.PersistentVars[keyVar] = keyValue
	}
	copiedImage.Vars = map[string]interface{}{}
	for keyVar, keyValue := range i.Vars {
		copiedImage.Vars[keyVar] = keyValue
	}
	copiedImage.Labels = map[string]string{}
	for keyVar, keyValue := range i.Labels {
		copiedImage.Labels[keyVar] = keyValue
	}

	if i.Children != nil {
		for _, child := range i.Children {
			copiedImage.Children = append(copiedImage.Children, child)
		}
	}

	return copiedImage, nil
}

// YAMLMarshal marshals the image to YAML
func (i *Image) YAMLMarshal() ([]byte, error) {
	errContext := "(image::YAMLMarshal)"

	marshaled, err := yaml.Marshal(i)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return marshaled, nil
}

// YAMLUnmarshal unmarshals the image from a YAML string
func (i *Image) YAMLUnmarshal(in []byte) error {
	errContext := "(image::YAMLUnmarshal)"

	err := yaml.Unmarshal(in, i)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
