package functional

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/suite"
)

type InitializeFunctionalTestsSuite struct {
	*FunctionalTestsSuite
}

func NewInitializeFunctionalTestsSuite(opts ...OptionsFunc) *InitializeFunctionalTestsSuite {

	functional := NewTestSuite(opts...)
	s := &InitializeFunctionalTestsSuite{
		functional,
	}

	return s
}

func (s *InitializeFunctionalTestsSuite) SetupTest() {
	s.TearDownTest()

	err := s.stack.Execute("up --detach registry")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *InitializeFunctionalTestsSuite) TearDownTest() {
	err := s.stack.Execute("rm --stop --force --volumes registry")
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

	err = s.stack.Execute("exec stevedore stevedore initialize --builders-path /builders --concurrency 3 --config /stevedore.yaml --credentials-format yaml --credentials-storage-type envvars --enable-semver-tags --force --generate-credentials-encryption-key --images-path /images --log-path-file /logs --push-images --semver-tags-template '{{ .Major }}'")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func initializeSetupSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	var err error

	err = stack.Execute("run --detach --rm --entrypoint sleep stevedore infinity")
	return err
}

func initializeTearDownSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	err := stack.Down()
	return err
}

func TestInitializeFunctionalTests(t *testing.T) {

	options := &docker.Options{
		WorkingDir:  ".",
		ProjectName: strings.ToLower(t.Name()),
	}
	options.Logger = logger.New(&quiteLogger{})

	project := NewDockerComposeProject(options)
	command := NewDockerComposeCommand(t, project)

	stack := NewDockerComposeStack(
		WithCommand(command),
		WithStackPreUpAction("build"),
	)

	s := NewInitializeFunctionalTestsSuite(
		WithStack(stack),
		WithSetupSuiteFunc(initializeSetupSuiteFunc),
		WithTearDownSuiteFunc(initializeTearDownSuiteFunc),
	)

	suite.Run(t, s)
}
