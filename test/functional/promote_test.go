package functional

import (
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/test/helpers"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/suite"
)

// PromoteFunctionalTestsSuite is a struct that defines the functional test suite for promote command
type PromoteFunctionalTestsSuite struct {
	*helpers.FunctionalTestsSuite
}

func NewPromoteFunctionalTestsSuite(opts ...helpers.OptionsFunc) *PromoteFunctionalTestsSuite {

	functional := helpers.NewFunctionalTestsSuite(opts...)
	s := &PromoteFunctionalTestsSuite{
		functional,
	}

	return s
}

func (s *PromoteFunctionalTestsSuite) SetupTest() {
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

func (s *PromoteFunctionalTestsSuite) TearDownTest() {
	if s.CommandFactory == nil {
		s.T().Fatal("You need to define a command factory")
	}

	err := s.Stack.Execute(s.CommandFactory.Command(helpers.RmCommand).WithCommandArgs("--stop --force --volumes registry"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *PromoteFunctionalTestsSuite) TestPromoteImage() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore promote docker-hub.stevedore.test:5000/library/busybox:latest --promote-image-registry-host registry.stevedore.test --promote-image-tag 1.2.3 --force-promote-source-image --use-source-image-from-remote --enable-semver-tags --remove-local-images-after-push"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/busybox:latest"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/busybox:1.2.3"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/busybox:1.2"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/busybox:1"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *PromoteFunctionalTestsSuite) TestPromoteImageOverwriteSemversTagTemplates() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore stevedore promote docker-hub.stevedore.test:5000/library/busybox:latest --use-source-image-from-remote --promote-image-registry-host registry.stevedore.test --promote-image-tag 1.2.3 --force-promote-source-image --use-source-image-from-remote --enable-semver-tags --remove-local-images-after-push --semver-tags-template \"{{ .Major }}_test\" --semver-tags-template \"{{ .Major }}.{{ .Minor }}_test\""))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/busybox:1_test"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.Stack.Execute(s.CommandFactory.Command(helpers.ExecCommand).WithCommandArgs("--workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/busybox:1.2_test"))
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func TestPromoteFunctionalTests(t *testing.T) {

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
		helpers.WithStackPostUpCommand(factory.Command(helpers.ExecCommand).WithCommandArgs("stevedore stevedore copy busybox:latest --use-source-image-from-remote --promote-image-registry-host docker-hub.stevedore.test:5000")),
	)

	// Define the suite
	s := NewPromoteFunctionalTestsSuite(
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
