package builder

import (
	"bytes"
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/varsmap"
	"gopkg.in/yaml.v2"
)

const (
	arrayOptionAssignment = "="
)

// Builder serializes each builder defined on user configuration
type Builder struct {
	Name       string          `yaml:"name"`
	Driver     string          `yaml:"driver"`
	Options    *BuilderOptions `yaml:"options"`
	VarMapping varsmap.Varsmap `yaml:"variables_mapping"`
}

// BuilderOptions are the options that can be set on a builder
type BuilderOptions struct {
	// AnsiblePlaybookDriverOptions are the options that can be set on a builder for ansible-playbook driver
	Playbook  string `yaml:"playbook"`
	Inventory string `yaml:"inventory"`
	// DockerDriverOptions are the options that can be set on a builder for docker driver
	Dockerfile string                        `yaml:"dockerfile"`
	Context    []*DockerDriverContextOptions `yaml:"context"`
}

// DockerDriverContextOptions are the context definitions to build a docker image
type DockerDriverContextOptions struct {
	Path string                         `yaml:"path"`
	Git  *DockerDriverGitContextOptions `yaml:"git"`
}

// DockerDriverGitContextOptions defines a build context from a git repository
type DockerDriverGitContextOptions struct {
	// Path must be set when docker build context is located in a subpath inside the repository
	Path string `yaml:"path"`
	// Repository which will be used as docker build context
	Repository string `yaml:"repository"`
	// Reference is the name of the branch to clone. By default is used 'master'
	Reference string `yaml:"reference"`
	// Auth is the authentication configuration for the repository
	Auth *DockerDriverGitContextAuthOptions `yaml:"auth"`
}

// DockerDriverGitContextAuthOptions defines the authentication for a git context
type DockerDriverGitContextAuthOptions struct {
	// Username is the username to use for basic authentication and used for git over http
	Username string `yaml:"username"`
	// Password is the password to use for basic authentication for git over http
	Password string `yaml:"password"`
	// GitSSHUser is the git ssh user to use for ssh authentication
	GitSSHUser string `yaml:"git_ssh_user"`
	// PrivateKeyFile is the path to the private key to use for ssh authentication
	PrivateKeyFile string `yaml:"private_key_file"`
	// PrivateKeyPassword is the password to use for the private key
	PrivateKeyPassword string `yaml:"private_key_password"`
	// CredentialsID is the id of the credentials to use for docker registry authentication
	CredentialsID string `yaml:"credentials_id"`
}

// NewBuilder creates a new builder
func NewBuilder() *Builder {
	return &Builder{
		Options:    &BuilderOptions{},
		VarMapping: varsmap.New(),
	}
}

// NewBuilderFromByteArray creates a new builder from a byte array
func NewBuilderFromByteArray(data []byte) (*Builder, error) {
	var builder *Builder

	errContext := "(builder::NewBuilderFromByteArray)"

	err := yaml.Unmarshal(data, &builder)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Builder could not be created.\nfound:\n'%s'\n", string(data)), err)
	}

	return builder, nil
}

// NewBuilderFromIOReader creates a new builder from an io reader
func NewBuilderFromIOReader(reader io.Reader) (*Builder, error) {
	var builder *Builder
	var buff bytes.Buffer
	var err error

	errContext := "(builder::NewBuilderFromIOReader)"

	_, err = buff.ReadFrom(reader)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	err = yaml.Unmarshal(buff.Bytes(), &builder)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Builder could not be created.\nfound:\n'%s'\n", buff.String()), err)
	}

	return builder, nil
}

// WithName sets the name of the builder
func (b *Builder) WithName(name string) {
	b.Name = name
}

// WithDriver sets the driver of the builder
func (b *Builder) WithDriver(driver string) {
	b.Driver = driver
}

// WithOptions sets the options of the builder
func (b *Builder) WithOptions(options *BuilderOptions) {
	b.Options = options
}

// WithVarMapping sets the variable mapping of the builder
func (b *Builder) WithVarMapping(mapping varsmap.Varsmap) {
	b.VarMapping = mapping
}

// // New
// func NewBuilder(name, driver string, options map[string]interface{}, vmap varsmap.Varsmap) (*Builder, error) {

// 	if name == "" {
// 		return nil, errors.New("(builder::NewBuilder", "Name must be provided to create a builder")
// 	}

// 	if driver == "" {
// 		return nil, errors.New("(builder::NewBuilder", "Driver must be provided to create a builder")
// 	}

// 	if options == nil {
// 		options = map[string]interface{}{}
// 	}

// 	if vmap == nil {
// 		vmap = varsmap.New()
// 	} else {
// 		vmap.Combine(varsmap.New())
// 	}

// 	b := &Builder{
// 		Name:       name,
// 		Driver:     driver,
// 		Options:    options,
// 		VarMapping: vmap,
// 	}
// 	return b, nil
// }

// // SanetizeBuilder ensures that a builders has been created with all required parameters
// func (b *Builder) SanetizeBuilder(name string) error {

// 	if b == nil {
// 		return errors.New("(builder::SanetizeBuilder)", "Builder is nil")
// 	}
// 	if len(b.Name) <= 0 {
// 		b.Name = name
// 	}

// 	if len(b.Driver) <= 0 {
// 		b.Driver = defaultbuilder.DriverName
// 	}

// 	if b.VarMapping == nil {
// 		b.VarMapping = varsmap.New()
// 	} else {
// 		b.VarMapping.Combine(varsmap.New())
// 	}

// 	return nil
// }

// // ToArray
// func (b *Builder) ToArray() ([]string, error) {
// 	arrayBuilder := []string{}
// 	arrayBuilder = append(arrayBuilder, b.Name)
// 	arrayBuilder = append(arrayBuilder, b.Driver)
// 	arrayBuilder = append(arrayBuilder, b.listArrayOptions()...)

// 	return arrayBuilder, nil
// }

// // listArrayOptions
// func (b *Builder) listArrayOptions() []string {
// 	options := []string{}
// 	for option, value := range b.Options {
// 		switch value.(type) {

// 		case string:
// 			options = append(options, fmt.Sprintf("%s%s%s", option, arrayOptionAssignment, value.(string)))
// 		case int:
// 			options = append(options, fmt.Sprintf("%s%s%v", option, arrayOptionAssignment, fmt.Sprintf("%v", value.(int))))
// 		default:
// 			options = append(options, fmt.Sprintf("%s%s%v", option, arrayOptionAssignment, fmt.Sprintf("%v", value)))
// 		}
// 	}

// 	sort.Strings(options)
// 	return options
// }
