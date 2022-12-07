package image

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	data "github.com/apenella/go-common-utils/data"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/docker/distribution/reference"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"gopkg.in/yaml.v3"
)

const (
	// ImageWildcardVersionSymbol is the wildcard version
	ImageWildcardVersionSymbol = "*"
	// NameFilterAttribute is the attribute's filter value to filter by name
	NameFilterAttribute = "name"
	// VersionFilterAttribute is the attribute's filter value to filter by version
	VersionFilterAttribute = "version"
	// RegistryHostFilterAttribute is the attribute's filter value to filter by registry host
	RegistryHostFilterAttribute = "registry"
	// RegistryNamespaceFilterAttribute is the attribute's filter value to filter by namespace
	RegistryNamespaceFilterAttribute = "namespace"
	// UndefinedStringValue defines an empty value rather that and empty string
	UndefinedStringValue = "-"
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
	// Parent is the parent image
	Parent *Image `yaml:"-"`
	// PersistentLabels persistent labels
	PersistentLabels map[string]string `yaml:"persistent_labels"`
	// PresistentVars are persistent variables
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

	addChildMutex sync.RWMutex
}

// OptionFunc is an option to pass to NewImage
type OptionFunc func(*Image)

// NewImage creates a new image
func NewImage(name, version, registryHost, registryNamesapace string, opt ...OptionFunc) (*Image, error) {

	errContext := "(core::domain::image::NewImage)"

	// Image name normalization
	imageName := name

	if imageName == "" {
		return nil, errors.New(errContext, "Image name is not provided")
	}

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
		return nil, errors.New(errContext, "", err)
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

// WithPersistentLabels sets the persistent variables
func WithPersistentLabels(persistentLabels map[string]string) OptionFunc {
	return func(i *Image) {
		i.PersistentLabels = persistentLabels
	}
}

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
	errContext := "(core::domain::image::Parse)"

	referenceName, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Image '%s' could not be parsed", name), err)
	}
	// in case that no version is specified, it will be the tag it as 'latest'
	referenceName = reference.TagNameOnly(referenceName)

	return &Image{
		Name:              reference.Path(referenceName)[strings.LastIndex(reference.Path(referenceName), "/")+1:],
		Version:           referenceName.String()[strings.LastIndex(referenceName.String(), ":")+1:],
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

// // NormalizedNamed normalizes the image name
// func (i *Image) DockerNormalizedNamed() (string, error) {
// 	var err error
// 	errContext := "(core::domain::image::DockerNormalizedNamed)"

// 	if i.Name == "" {
// 		return "", errors.New(errContext, "Image name is empty")
// 	}

// 	if i.Version == "" {
// 		return "", errors.New(errContext, "Image version is empty")
// 	}

// 	if i.RegistryHost == "" {
// 		return "", errors.New(errContext, "Registry host is empty")
// 	}

// 	if i.RegistryNamespace == "" {
// 		return "", errors.New(errContext, "Registry namespace is empty")
// 	}

// 	name := fmt.Sprintf("%s/%s/%s:%s", i.RegistryHost, i.RegistryNamespace, i.Name, i.Version)

// 	_, err = reference.ParseNormalizedNamed(name)
// 	if err != nil {
// 		return "", errors.New(errContext, fmt.Sprintf("Image name '%s' could not be normalized", name), err)
// 	}

// 	return name, nil
// }

// Sanetize normalizes the image name
func (i *Image) Sanetize() error {
	var err error
	errContext := "(core::domain::image::Sanetize)"
	sanitizeTable := map[string]string{
		"+": "_",
	}

	for dirty, sane := range sanitizeTable {
		i.Version = strings.ReplaceAll(i.Version, dirty, sane)
	}

	err = i.sanetizeBuilder()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func (i *Image) sanetizeBuilder() error {

	var err error
	var str string
	var builderSane *builder.Builder
	errContext := "(core::domain::image::sanetizeBuilder)"

	value := reflect.ValueOf(i.Builder)
	switch value.Kind() {
	case reflect.String:
		return nil

	case reflect.Map: // interface conversion: interface {} is map[interface {}]interface {}, not *builder.Builder
		str, err = data.ObjectToYamlString(i.Builder)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal([]byte(str), &builderSane)
		if err != nil {
			return errors.New(errContext, "", err)
		}

		if builderSane.Name == "" && i.Name != "" && i.Version != "" {
			builderSane.Name = fmt.Sprintf("%s:%s", i.Name, i.Version)
		}
		i.Builder = builderSane

		return nil

	case reflect.TypeOf(builderSane).Kind():
		if i.Builder.(*builder.Builder).Name == "" && i.Name != "" {
			i.Builder.(*builder.Builder).Name = i.Name
		}

		return nil
	default:
		return nil
	}
}

// Copy method return a copy of the instanced Image
func (i *Image) Copy() (*Image, error) {

	copiedImage := &Image{
		Name:              i.Name,
		Version:           i.Version,
		RegistryHost:      i.RegistryHost,
		RegistryNamespace: i.RegistryNamespace,
	}

	opts := []OptionFunc{}

	if i.Builder != nil {
		opts = append(opts, WithBuilder(i.Builder))
	}

	if i.Parent != nil {
		opts = append(opts, WithParent(i.Parent))
	}
	copiedImage.Options(opts...)

	copiedImage.Children = append([]*Image{}, i.Children...)

	copiedImage.Tags = append([]string{}, i.Tags...)

	copiedImage.PersistentVars = map[string]interface{}{}
	for keyVar, keyValue := range i.PersistentVars {
		copiedImage.PersistentVars[keyVar] = keyValue
	}
	copiedImage.Vars = map[string]interface{}{}
	for keyVar, keyValue := range i.Vars {
		copiedImage.Vars[keyVar] = keyValue
	}
	copiedImage.PersistentLabels = map[string]string{}
	for keyVar, keyValue := range i.PersistentLabels {
		copiedImage.PersistentLabels[keyVar] = keyValue
	}
	copiedImage.Labels = map[string]string{}
	for keyVar, keyValue := range i.Labels {
		copiedImage.Labels[keyVar] = keyValue
	}

	return copiedImage, nil
}

// IsWildcardImage returns true if the image is a wildcard image
func (i *Image) IsWildcardImage() bool {
	return i.Version == ImageWildcardVersionSymbol
}

// YAMLMarshal marshals the image to YAML
func (i *Image) YAMLMarshal() ([]byte, error) {
	errContext := "(core::domain::image::YAMLMarshal)"

	marshaled, err := yaml.Marshal(i)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return marshaled, nil
}

// YAMLUnmarshal unmarshals the image from a YAML string
func (i *Image) YAMLUnmarshal(in []byte) error {
	errContext := "(core::domain::image::YAMLUnmarshal)"

	err := yaml.Unmarshal(in, i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
