package helpers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/mattn/go-shellwords"
	"github.com/stretchr/testify/require"
)

const (
	// BuildCommand is the command to build or rebuild services
	BuildCommand = "build"
	// DownCommand is the command to stop and remove containers and networks
	DownCommand = "down"
	// ExecCommand is the command to execute a command in a running container.
	ExecCommand = "exec"
	// KillCommand is the command to force stop service containers.
	KillCommand = "kill"
	// LogsCommand is the command to view output from containers
	LogsCommand = "logs"
	// RestartCommand is the command to restart service containers
	RestartCommand = "restart"
	// RmCommand is the command to removes stopped service containers
	RmCommand = "rm"
	// RunCommand is the command to run a one-off command on a service.
	RunCommand = "run"
	// StartCommand is the command to start services
	StartCommand = "start"
	// StopCommand is the command to stop services
	StopCommand = "stop"
	// UpCommand is the command to create and start containers
	UpCommand = "up"
)

// DockerComposeTerratestCommand that defines a Terratest docker-compose command
type DockerComposeTerratestCommand struct {
	command     string
	commandArgs string
	options     *docker.Options
	testing     *testing.T
	verbose     bool
}

// NewDockerComposeTerratestCommand creates a new DockerComposeTerratestCommand
func NewDockerComposeTerratestCommand(t *testing.T, options *docker.Options) *DockerComposeTerratestCommand {
	return &DockerComposeTerratestCommand{
		options: options,
		testing: t,
	}
}

// WithCommand sets the command for the DockerComposeTerratestCommand
func (c *DockerComposeTerratestCommand) WithCommand(cmd string) *DockerComposeTerratestCommand {
	c.command = cmd
	return c
}

// WithCommandArgs sets the command arguments for the DockerComposeTerratestCommand
func (c *DockerComposeTerratestCommand) WithCommandArgs(args string) *DockerComposeTerratestCommand {
	c.commandArgs = args
	return c
}

// WithVerbose sets the verbose flag for the DockerComposeTerratestCommand
func (c *DockerComposeTerratestCommand) WithVerbose() *DockerComposeTerratestCommand {
	c.verbose = true
	return c
}

// Execute runs the DockerComposeTerratestCommand
func (c *DockerComposeTerratestCommand) Execute() (string, error) {
	var err error
	var result string
	_ = err
	if c.command == "" {
		return "", errors.New("Docker-compose command requires a command")
	}

	if c.options == nil {
		return "", errors.New("Docker-compose command requires a project")
	}

	cmd := []string{c.command}
	cmd = append(cmd, parseArgs(c.commandArgs)...)

	// if c.verbose {
	fmt.Printf(" - [%s] %s\n", c.options.ProjectName, strings.Join(cmd, " "))
	// }

	result, err = docker.RunDockerComposeE(c.testing, c.options, cmd...)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Docker-compose command failed running '%s'. %s", strings.Join(cmd, " "), err.Error()))
	}

	return result, nil
}

func (c *DockerComposeTerratestCommand) AssertExectedResult(expected, result string) {
	require.Contains(c.testing, result, expected)
}

// Custom function to parse command string into arguments
func parseArgs(command string) []string {
	// Use the shellwords package to split the command string
	args, err := shellwords.Parse(command)
	if err != nil {
		fmt.Printf("Error parsing command string: %s", err)
	}

	// If you want to remove empty arguments, you can do that here
	var cleanedArgs []string
	for _, arg := range args {
		if arg != "" {
			cleanedArgs = append(cleanedArgs, arg)
		}
	}

	return cleanedArgs
}
