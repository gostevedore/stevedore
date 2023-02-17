package builder

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"gopkg.in/yaml.v3"
)

// BuilderOptions are the options that can be set on a builder
type BuilderOptions struct {
	// AnsiblePlaybookDriverOptions are the options that can be set on a builder for ansible-playbook driver
	Playbook  string `yaml:"playbook"`
	Inventory string `yaml:"inventory"`
	// DockerDriverOptions are the options that can be set on a builder for docker driver
	Dockerfile string `yaml:"dockerfile"`
	// Context    []*DockerDriverContextOptions `yaml:"context"`
	Context interface{} `yaml:"context"`
}

func (o *BuilderOptions) GetContext() ([]*DockerDriverContextOptions, error) {

	var contextList []*DockerDriverContextOptions
	errContext := "(core::domain::builder::BuilderOptions::GetContext)"

	switch o.Context.(type) {
	case []*DockerDriverContextOptions:
		contextList = append(contextList, o.Context.([]*DockerDriverContextOptions)...)
	case *DockerDriverContextOptions:
		contextList = append([]*DockerDriverContextOptions{}, o.Context.(*DockerDriverContextOptions))
	case map[string]interface{}:
		context := &DockerDriverContextOptions{}

		contextDefinitionBytes, err := yaml.Marshal(o.Context)
		if err != nil {
			return nil, errors.New(errContext, "There is an error marshaling the Docker driver context options", err)
		}

		err = yaml.Unmarshal(contextDefinitionBytes, &context)
		if err != nil {
			return nil, errors.New(errContext, fmt.Sprintf("Docker driver context options could not be created.\nfound:\n'%s'\n", string(contextDefinitionBytes)), err)
		}

		contextList = append([]*DockerDriverContextOptions{}, context)

	case []interface{}:
		contextDefinitionBytes, err := yaml.Marshal(o.Context)
		if err != nil {
			return nil, errors.New(errContext, "There is an error marshaling the Docker driver context options", err)
		}

		err = yaml.Unmarshal(contextDefinitionBytes, &contextList)
		if err != nil {
			return nil, errors.New(errContext, fmt.Sprintf("Docker driver context options could not be created.\nfound:\n'%s'\n", string(contextDefinitionBytes)), err)
		}
	default:
		return nil, errors.New(errContext, "Docker driver context options format is not valid")
	}

	return contextList, nil
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
	// CredentialsID is the id of the credentials on credentials store to use to authenticate to the git repository
	CredentialsID string `yaml:"credentials_id"`
}
