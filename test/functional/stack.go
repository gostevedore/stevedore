package functional

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/require"
)

type DockerComposeProject struct {
	options *docker.Options
}

func NewDockerComposeProject(options *docker.Options) *DockerComposeProject {

	return &DockerComposeProject{
		options: options,
	}
}

type DockerComposeCommand struct {
	project *DockerComposeProject
	testing *testing.T
}

func NewDockerComposeCommand(t *testing.T, project *DockerComposeProject) *DockerComposeCommand {
	return &DockerComposeCommand{
		project: project,
		testing: t,
	}
}

func (c *DockerComposeCommand) Execute(cmd string) (string, error) {
	var err error
	var result string

	if c.project == nil {
		return "", errors.New("Docker-compose command requires a project")
	}

	cmds := strings.Split(cmd, " ")
	fmt.Println("  -", cmd)
	result, err = docker.RunDockerComposeE(c.testing, c.project.options, cmds...)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *DockerComposeCommand) ExpectedInResult(expected, result string) {
	require.Contains(c.testing, expected, result)
}

// DockerComposeStackOptionsFunc is a function used to configure the service
type DockerComposeStackOptionsFunc func(*DockerComposeStack)

type DockerComposeStack struct {
	command        *DockerComposeCommand
	PreUpAction    []string
	PostUpAction   []string
	PreDownAction  []string
	PostDownAction []string
}

func NewDockerComposeStack(opts ...DockerComposeStackOptionsFunc) *DockerComposeStack {
	stack := &DockerComposeStack{
		PreUpAction:    make([]string, 0, 10),
		PostUpAction:   make([]string, 0, 10),
		PreDownAction:  make([]string, 0, 10),
		PostDownAction: make([]string, 0, 10),
	}
	stack.Options(opts...)

	return stack
}

func (s *DockerComposeStack) Options(opts ...DockerComposeStackOptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

func WithCommand(cmd *DockerComposeCommand) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.command = cmd
	}
}

func WithStackPreUpAction(cmd ...string) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PreUpAction = append(dcs.PreUpAction, cmd...)
	}
}

func WithStackPostUpAction(cmd ...string) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PostUpAction = append(dcs.PostUpAction, cmd...)
	}
}

func WithStackPreDownAction(cmd ...string) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PreDownAction = append(dcs.PreDownAction, cmd...)
	}
}

func WithStackPostDownAction(cmd ...string) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PostDownAction = append(dcs.PostDownAction, cmd...)
	}
}

func (s *DockerComposeStack) Up(options ...string) error {
	if s.command == nil {
		return errors.New("Docker-compose stack requires command")
	}

	color.Green(" Pre-up actions")
	for _, action := range s.PreUpAction {
		_, err := s.command.Execute(action)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	color.Green(" Execute up")
	_, err := s.command.Execute(strings.Join(append([]string{"up"}, options...), " "))
	if err != nil {
		return err
	}
	fmt.Println()

	color.Green(" Post-up actions")
	for _, action := range s.PostUpAction {
		_, err := s.command.Execute(action)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

func (s *DockerComposeStack) DownAndUp(options ...string) error {
	if s.command == nil {
		return errors.New("Docker-compose stack requires command")
	}

	err := s.Down()
	if err != nil {
		return err
	}

	err = s.Up(options...)
	if err != nil {
		return err
	}

	return nil
}

func (s *DockerComposeStack) Down(options ...string) error {
	if s.command == nil {
		return errors.New("Docker-compose stack requires command")
	}

	color.Green(" Pre-down actions")
	for _, action := range s.PreDownAction {
		_, err := s.command.Execute(action)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	color.Green(" Execute down")
	_, err := s.command.Execute(strings.Join(append([]string{"down --volumes --remove-orphans --timeout 3"}, options...), " "))
	if err != nil {
		return err
	}
	fmt.Println()

	color.Green(" Post-down actions")
	for _, action := range s.PostDownAction {
		_, err := s.command.Execute(action)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

func (c *DockerComposeStack) Execute(cmd string) error {
	var err error
	_, err = c.command.Execute(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *DockerComposeStack) ExecuteAndCompare(cmd string, expectedResult string) error {
	var err error
	var result string

	result, err = c.command.Execute(cmd)
	if err != nil {
		return err
	}

	c.command.ExpectedInResult(expectedResult, result)

	return nil
}
