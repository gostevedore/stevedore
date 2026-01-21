package helpers

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

// DockerComposeStackOptionsFunc is a function used to configure the service
type DockerComposeStackOptionsFunc func(*DockerComposeStack)

// DockerComposeStack is a struct that defines the docker compose stack where the test will be executed
type DockerComposeStack struct {
	DownCommand     Commander
	PostDownCommand []Commander
	PostUpCommand   []Commander
	PreDownCommand  []Commander
	PreUpCommand    []Commander
	UpCommand       Commander
}

// NewDockerComposeStack creates a new DockerComposeStack
func NewDockerComposeStack(opts ...DockerComposeStackOptionsFunc) *DockerComposeStack {
	stack := &DockerComposeStack{
		PreUpCommand:    make([]Commander, 0, 10),
		PostUpCommand:   make([]Commander, 0, 10),
		PreDownCommand:  make([]Commander, 0, 10),
		PostDownCommand: make([]Commander, 0, 10),
	}
	stack.Options(opts...)

	return stack
}

// Options sets options for the DockerComposeStack
func (s *DockerComposeStack) Options(opts ...DockerComposeStackOptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

// WithUpCommand sets the up command for the DockerComposeStack
func WithUpCommand(cmd Commander) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.UpCommand = cmd
	}
}

// WithDownCommand sets the down command for the DockerComposeStack
func WithDownCommand(cmd Commander) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.DownCommand = cmd
	}
}

// WithStackPreUpCommand sets the pre-up commands for the DockerComposeStack
func WithStackPreUpCommand(cmd ...Commander) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PreUpCommand = append(dcs.PreUpCommand, cmd...)
	}
}

// WithStackPostUpCommand sets the post-up commands for the DockerComposeStack
func WithStackPostUpCommand(cmd ...Commander) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PostUpCommand = append(dcs.PostUpCommand, cmd...)
	}
}

// WithStackPreDownCommand sets the pre-down commands for the DockerComposeStack
func WithStackPreDownCommand(cmd ...Commander) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PreDownCommand = append(dcs.PreDownCommand, cmd...)
	}
}

// WithStackPostDownCommand sets the post-down commands for the DockerComposeStack
func WithStackPostDownCommand(cmd ...Commander) DockerComposeStackOptionsFunc {
	return func(dcs *DockerComposeStack) {
		dcs.PostDownCommand = append(dcs.PostDownCommand, cmd...)
	}
}

// Up creates and starts the stack where the tests will be executed
func (s *DockerComposeStack) Up() error {

	if s.UpCommand == nil {
		return errors.New("Docker compose requires a command to create and start the stack")
	}

	color.Green(" Pre-up actions\n")
	for _, command := range s.PreUpCommand {
		_, err := command.Execute()
		if err != nil {
			return fmt.Errorf("Error executing pre-up command: %v. %s", command, err)
		}
	}
	fmt.Println()

	color.Green(" Execute up\n")
	_, err := s.UpCommand.Execute()
	if err != nil {
		return err
	}
	fmt.Println()

	color.Green(" Post-up actions\n")
	for _, command := range s.PostUpCommand {
		_, err := command.Execute()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

// Down stops and removes the stack where the tests were executed
func (s *DockerComposeStack) Down() error {

	if s.DownCommand == nil {
		return errors.New("Docker compose requires a command to clean up the stack")
	}

	color.Green(" Pre-down actions\n")
	for _, command := range s.PreDownCommand {
		_, err := command.Execute()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	color.Green(" Execute down\n")
	_, err := s.DownCommand.Execute()
	if err != nil {
		return err
	}
	fmt.Println()

	color.Green(" Post-down actions\n")
	for _, command := range s.PostDownCommand {
		_, err := command.Execute()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

// DownAndUp cleans up the stack and creates and starts it again
func (s *DockerComposeStack) DownAndUp() error {

	err := s.Down()
	if err != nil {
		return err
	}

	err = s.Up()
	if err != nil {
		return err
	}

	return nil
}

// Execute executes a command in the stack
func (c *DockerComposeStack) Execute(cmd Commander) error {
	var err error
	_, err = cmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

// ExecuteAndCompare executes a command in the stack and compares the result with the expected result
func (c *DockerComposeStack) ExecuteAndCompare(cmd AssertAndCommander, expectedResult string) error {
	var err error
	var result string

	result, err = cmd.Execute()
	if err != nil {
		return err
	}

	cmd.AssertExectedResult(expectedResult, result)

	return nil
}
