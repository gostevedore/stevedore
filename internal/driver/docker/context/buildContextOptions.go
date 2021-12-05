package context

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"gopkg.in/yaml.v2"
)

const (
	// GitContextOptionsKey is the key which identify git context options
	GitContextOptionsKey = "git"
	// PathContextOptionsKey is the key which identify path context options
	PathContextOptionsKey = "path"
)

// DockerBuildContextOptions are the context definitions to build a docker image
type DockerBuildContextOptions struct {
	Path string             `yaml:"path"`
	Git  *GitContextOptions `yaml:"git"`
}

// GenerateBuildContextOptions returns the build context structure
func GenerateBuildContextOptions(context interface{}) (*DockerBuildContextOptions, error) {
	errContext := "(DockerBuildContextFactory::GenerateDockerBuildContext)"
	dockerBuildContext := &DockerBuildContextOptions{}

	if context == nil {
		return nil, errors.New(errContext, "Docker build context options are expected to build an image")
	}

	// data, err := yaml.Marshal(context)
	// if err != nil {
	// 	return nil, errors.New(errContext, "Error marshalling the context options", err)
	// }

	err := yaml.Unmarshal([]byte(context.(string)), dockerBuildContext)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Docker build context options are not properly configured\n found:\n%s\n", context.(string)), err)
	}

	return dockerBuildContext, nil
}

// GitContextOptions defines a build context from a git repository
type GitContextOptions struct {
	// Path must be set when docker build context is located in a subpath inside the repository
	Path string `yaml:"path"`
	// Repository which will be used as docker build context
	Repository string `yaml:"repository"`
	// Reference is the name of the branch to clone. By default is used 'master'
	Reference string `yaml:"reference"`
	// Auth is the authentication configuration for the repository
	Auth *GitContextAuthOptions `yaml:"auth"`
}

// GitContextAuthDefinition defines the authentication for a git context
type GitContextAuthOptions struct {
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
	// CredentialsId is the id of the credentials to use for docker registry authentication
	CredentialsId string `yaml:"credentials_id"`
}
