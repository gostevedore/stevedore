package image

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
func NewImage(name, version string, opt ...OptionFunc) *Image {
	image := &Image{
		Name:    name,
		Version: version,
	}
	for _, o := range opt {
		o(image)
	}
	return image
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
func WithRegistryHost(registryHost string) OptionFunc {
	return func(i *Image) {
		i.RegistryHost = registryHost
	}
}

// WithRegistryNamespace sets the registry namespace
func WithRegistryNamespace(registryNamespace string) OptionFunc {
	return func(i *Image) {
		i.RegistryNamespace = registryNamespace
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

// AddChild adds a child image
func (i *Image) AddChild(child *Image) {
	i.Children = append(i.Children, child)
}
