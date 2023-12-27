package functional

import (
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/test/helpers"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/suite"
)

// BuildFunctionalTestsSuite is a struct that defines the functional test suite for build command
type BuildFunctionalTestsSuite struct {
	*helpers.FunctionalTestsSuite
}

// NewBuildFunctionalTestsSuite creates a new BuildFunctionalTestsSuite
func NewBuildFunctionalTestsSuite(opts ...helpers.OptionsFunc) *BuildFunctionalTestsSuite {

	functional := helpers.NewFunctionalTestsSuite(opts...)
	s := &BuildFunctionalTestsSuite{
		functional,
	}

	return s
}

// SetupTest is executed before executing each test
func (s *BuildFunctionalTestsSuite) SetupTest() {
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
func (s *BuildFunctionalTestsSuite) TearDownTest() {

	if s.CommandFactory == nil {
		s.T().Fatal("You need to define a command factory")
	}

	err := s.Stack.Execute(s.CommandFactory.Command(helpers.RmCommand).WithCommandArgs("--stop --force --volumes registry"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageWithGitContext() {
	var err error

	s.T().Log("Testing build image with Git context")

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	// err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore stevedore build app2 --push-after-build")
	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore build app2 --push-after-build"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	// err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16")
	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageWithMultipleParents() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore build app1 --push-after-build"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-latest"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-20.04"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-busybox-latest"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageOnCascade() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore build alpine --build-on-cascade --push-after-build"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageWithSemVerEnabled() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore build app3 --enable-semver-tags --push-after-build"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2.3"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageWithWildcardVersion() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore build app3 --image-version 1.3.0-rc0.1+1234 --push-after-build --enable-semver-tags"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1-rc0.1"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3-rc0.1-1234"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3.0-rc0.1_1234"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

// TestBuildFunctionalTests is the entrypoint to execute the test suite that runs the functional tests for build command
func TestBuildFunctionalTests(t *testing.T) {

	t.Parallel()

	options := &docker.Options{
		WorkingDir:  "../functional",
		ProjectName: strings.ToLower(t.Name()),
	}
	options.Logger = logger.New(&helpers.QuiteLogger{})

	// Command factory creates commands to be executed in the stack
	factory := helpers.NewDockerComposeTerratestCommandFactory(t, options)

	// Stack or the environment where the tests are executed
	stack := helpers.NewDockerComposeStack(
		// Define the command that is going to be executed when the stack is up
		helpers.WithUpCommand(factory.Command(helpers.UpCommand).WithCommandArgs("--detach docker-hub gitserver stevedore")),
		// Define the command that is going to be executed when the stack is down
		helpers.WithDownCommand(factory.Command(helpers.DownCommand).WithCommandArgs("--remove-orphans --volumes --timeout 3")),

		// Define the commands that are going to be executed before the stack is up
		helpers.WithStackPreUpCommand(factory.Command(helpers.BuildCommand).WithCommandArgs("stevedore")),
		helpers.WithStackPreUpCommand(factory.Command(helpers.RunCommand).WithCommandArgs("--rm openssh -t rsa -q -N password -f id_rsa -C \"apenella@stevedore.test\"")),
		helpers.WithStackPreUpCommand(factory.Command(helpers.RunCommand).WithCommandArgs("--rm openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config /root/ssl/stevedore.test.cnf")),
		helpers.WithStackPreUpCommand(factory.Command(helpers.RunCommand).WithCommandArgs("--rm openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile /root/ssl/stevedore.test.cnf")),

		// Define the commands that are going to be executed after the stack is up
		helpers.WithStackPostUpCommand(
			helpers.NewRetryCommand(
				factory.Command(helpers.ExecCommand).WithCommandArgs("stevedore /usr/local/bin/wait-for-dockerd.sh")).
				WithPostRetryCommand(
					factory.Command(helpers.RestartCommand).WithCommandArgs("stevedore"),
				)),
		helpers.WithStackPostUpCommand(factory.Command(helpers.ExecCommand).WithCommandArgs("stevedore /prepare-images")),
	)

	// Define the suite
	s := NewBuildFunctionalTestsSuite(
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

	// Run the suite BuildFunctionalTestsSuite
	suite.Run(t, s)
}
