package image

import (
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/docker/distribution/reference"
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
	RegistryHost string `yaml:"registry"`
	// RegistryNamespace is the namespace of the registry
	RegistryNamespace string `yaml:"namespace"`
	// Tags is a list of extra tags
	Tags []string `yaml:"tags"`
	// Vars is a map of variables
	Vars map[string]interface{} `yaml:"vars"`
	// Version is the version of the image
	Version string `yaml:"version"`
	// Parent is the parent image
	Parent *Image `yaml:"-"`
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

	referenceName, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return nil, errors.New(errContext, "Image could not be created", err)
	}
	// in case that no version is specified, it will be the tag it as 'latest'
	referenceName = reference.TagNameOnly(referenceName)

	fmt.Println(">>>>", referenceName.String())

	// Image
	image := &Image{
		Name:              reference.Path(referenceName)[strings.LastIndex(reference.Path(referenceName), "/")+1:],
		Version:           referenceName.String()[strings.IndexRune(referenceName.String(), ':')+1:],
		RegistryHost:      reference.Domain(referenceName),
		RegistryNamespace: reference.Path(referenceName)[:strings.LastIndex(reference.Path(referenceName), "/")],
	}
	for _, o := range opt {
		o(image)
	}
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

// AddChild adds a child image
func (i *Image) AddChild(child *Image) {
	i.Children = append(i.Children, child)
}
