package functional

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/suite"
)

type BuildFunctionalTestsSuite struct {
	*FunctionalTestsSuite
}

func NewBuildFunctionalTestsSuite(opts ...OptionsFunc) *BuildFunctionalTestsSuite {

	functional := NewTestSuite(opts...)
	s := &BuildFunctionalTestsSuite{
		functional,
	}

	return s
}

func (s *BuildFunctionalTestsSuite) SetupTest() {
	s.TearDownTest()

	err := s.stack.Execute("up --detach registry")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TearDownTest() {
	err := s.stack.Execute("rm --stop --force --volumes registry")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageWithGitContext() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore stevedore build app2 --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16")
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

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore stevedore build app1 --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-20.04")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-busybox-latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
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

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore stevedore build alpine --build-on-cascade --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16")
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

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore stevedore build app3 --enable-semver-tags --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2.3")
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

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore stevedore build app3 --image-version 1.3.0-rc0.1+1234 --push-after-build --enable-semver-tags")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1-rc0.1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3-rc0.1-1234")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec --workdir /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3.0-rc0.1_1234")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func buildSetupSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	var err error

	err = stack.DownAndUp("--detach docker-hub gitserver stevedore")
	return err
}

func buildTearDownSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	err := stack.Down()
	return err
}

func TestBuildFunctionalTests(t *testing.T) {

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
		WithStackPreUpAction("run --rm openssh -t rsa -q -N password -f id_rsa -C \"apenella@stevedore.test\""),
		WithStackPreUpAction("run --rm openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config /root/ssl/stevedore.test.cnf"),
		WithStackPreUpAction("run --rm openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile /root/ssl/stevedore.test.cnf"),
		WithStackPostUpAction("exec stevedore /usr/local/bin/wait-for-dockerd.sh"),
		WithStackPostUpAction("exec stevedore /prepare-images"),
	)

	s := NewBuildFunctionalTestsSuite(
		WithStack(stack),
		WithSetupSuiteFunc(buildSetupSuiteFunc),
		WithTearDownSuiteFunc(buildTearDownSuiteFunc),
	)

	suite.Run(t, s)
}
