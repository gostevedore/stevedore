package functional

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/suite"
)

type PromoteFunctionalTestsSuite struct {
	FunctionalTestsSuite
}

func NewPromoteFunctionalTestsSuite(opts ...OptionsFunc) *PromoteFunctionalTestsSuite {

	functional := NewTestSuite(opts...)
	s := &PromoteFunctionalTestsSuite{
		*functional,
	}

	return s
}

func (s *PromoteFunctionalTestsSuite) SetupTest() {
	s.TearDownTest()

	err := s.stack.Execute("up -d registry")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *PromoteFunctionalTestsSuite) TearDownTest() {
	err := s.stack.Execute("rm --stop --force --volumes registry")
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

	err = s.stack.Execute("exec -w /app/test/stack/client/stevedore stevedore stevedore promote docker-hub.stevedore.test:5000/library/ubuntu:latest --promote-image-registry-host registry.stevedore.test --promote-image-tag 1.2.3 --force-promote-source-image --use-source-image-from-remote --enable-semver-tags --remove-local-images-after-push")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/ubuntu:latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/ubuntu:1.2.3")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/ubuntu:1.2")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("exec -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/library/ubuntu:1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func promoteSetupSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	var err error

	err = stack.DownAndUp("-d docker-hub registry stevedore")
	return err
}

func promoteTearDownSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	err := stack.Down()
	return err
}

func TestPromoteFunctionalTests(t *testing.T) {

	options := &docker.Options{
		WorkingDir:  ".",
		ProjectName: strings.ToLower(t.Name()),
	}

	project := NewDockerComposeProject(options)
	command := NewDockerComposeCommand(t, project)

	stack := NewDockerComposeStack(
		WithCommand(command),
		WithStackPreUpAction("build"),
		WithStackPreUpAction("run --rm openssh -t rsa -q -N password -f id_rsa -C \"apenella@stevedore.test\""),
		WithStackPreUpAction("run --rm openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config /root/ssl/stevedore.test.cnf"),
		WithStackPreUpAction("run --rm openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile /root/ssl/stevedore.test.cnf"),
		// fixed timeout, try to improve by checking the status of dockerd with while !nc -vz localhost 2376; do sleep 1s;done
		WithStackPostUpAction("exec stevedore sleep 10s"),
		WithStackPostUpAction("exec stevedore stevedore copy ubuntu:latest --use-source-image-from-remote --promote-image-registry-host docker-hub.stevedore.test:5000"),
	)

	s := NewPromoteFunctionalTestsSuite(
		WithStack(stack),
		WithSetupSuiteFunc(promoteSetupSuiteFunc),
		WithTearDownSuiteFunc(promoteTearDownSuiteFunc),
	)

	suite.Run(t, s)
}
