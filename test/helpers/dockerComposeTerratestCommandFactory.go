package helpers

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
)

// DockerComposeTerratestCommandFactory that defines a Terratest docker-compose command factory
type DockerComposeTerratestCommandFactory struct {
	testing *testing.T
	options *docker.Options
}

// NewDockerComposeTerratestCommandFactory creates a new DockerComposeTerratestCommandFactory
func NewDockerComposeTerratestCommandFactory(t *testing.T, options *docker.Options) *DockerComposeTerratestCommandFactory {
	return &DockerComposeTerratestCommandFactory{
		options: options,
		testing: t,
	}
}

// Command creates a new DockerComposeTerratestCommand for the command cmd
func (f *DockerComposeTerratestCommandFactory) Command(cmd string) *DockerComposeTerratestCommand {
	return &DockerComposeTerratestCommand{
		options: f.options,
		testing: f.testing,
		command: cmd,
	}
}
