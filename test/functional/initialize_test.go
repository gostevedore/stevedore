package functional

import (
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/test/helpers"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/suite"
)

// InitializeFunctionalTestsSuite is a struct that defines the functional test suite for initialize command
type InitializeFunctionalTestsSuite struct {
	*helpers.FunctionalTestsSuite
}

// NewInitializeFunctionalTestsSuite creates a new InitializeFunctionalTestsSuite
func NewInitializeFunctionalTestsSuite(opts ...helpers.OptionsFunc) *InitializeFunctionalTestsSuite {

	functional := helpers.NewFunctionalTestsSuite(opts...)
	s := &InitializeFunctionalTestsSuite{
		functional,
	}

	return s
}

// SetupTest is executed before executing each test
func (s *InitializeFunctionalTestsSuite) SetupTest() {
	s.TearDownTest()

	if s.CommandFactory == nil {
		s.T().Fatal("You need to define a command factory")
	}

	err := s.Stack.Execute(s.CommandFactory.Command(helpers.UpCommand).WithCommandArgs("--detach registry"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

// TearDownTest is executed after executing each test
func (s *InitializeFunctionalTestsSuite) TearDownTest() {

	if s.CommandFactory == nil {
		s.T().Fatal("You need to define a command factory")
	}

	err := s.Stack.Execute(s.CommandFactory.Command(helpers.RmCommand).WithCommandArgs("--stop --force --volumes registry"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *InitializeFunctionalTestsSuite) TestInitialize() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	expectedResult := "6fd4d739c5c98992db8f85c4889308fa"
	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("stevedore stevedore initialize --builders-path /builders --concurrency 3 --config /stevedore.yaml --credentials-format yaml --credentials-storage-type envvars --enable-semver-tags --force --credentials-encryption-key 12345asdfg --images-path /images --log-path-file /logs --push-images --semver-tags-template '{{ .Major }}'"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.ExecuteAndCompare(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("stevedore md5sum /stevedore.yaml"), expectedResult)
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

}

// TestInitializeFunctionalTests is the entrypoint to execute the test suite that runs the functional tests for initialize command
func TestInitializeFunctionalTests(t *testing.T) {

	t.Parallel()

	options := &docker.Options{
		WorkingDir:     "../functional",
		ProjectName:    strings.ToLower(t.Name()),
		EnableBuildKit: true,
	}
	options.Logger = logger.New(&helpers.QuiteLogger{})

	// Command factory creates commands to be executed in the stack
	factory := helpers.NewDockerComposeTerratestCommandFactory(t, options)

	// Stack or the environment where the tests are executed
	stack := helpers.NewDockerComposeStack(
		// Define the command that is going to be executed when the stack is up
		helpers.WithUpCommand(factory.Command(helpers.RunCommand).WithCommandArgs("--detach --rm --entrypoint sleep stevedore infinity")),
		// Define the command that is going to be executed when the stack is down
		helpers.WithDownCommand(factory.Command(helpers.DownCommand).WithCommandArgs("--remove-orphans --volumes --timeout 3")),

		// Define the commands that are going to be executed before the stack is up
		helpers.WithStackPreUpCommand(factory.Command(helpers.BuildCommand).WithCommandArgs("stevedore")),
	)

	// Define the suite
	s := NewInitializeFunctionalTestsSuite(
		helpers.WithCommandFactory(factory),
		helpers.WithStack(stack),
		// Define the function to execute before executing the tests suite
		helpers.WithSetupSuiteFunc(
			func() error {
				var err error

				err = stack.DownAndUp()
				if err != nil {
					t.Log(err)
					t.Fail()
				}

				return err
			},
		),
		// Define the function to execute after executing the tests suite
		helpers.WithTearDownSuiteFunc(func() error {
			var err error

			err = stack.Down()
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			return err
		}),
	)

	// Run the suite InitializeFunctionalTestsSuite
	suite.Run(t, s)
}
