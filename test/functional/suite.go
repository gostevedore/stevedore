package functional

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/suite"
)

type FunctionalTestsSuite struct {
	suite.Suite
	stack DockerComposeStack
}

func (s *FunctionalTestsSuite) SetupSuite() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	options := &docker.Options{
		WorkingDir:  ".",
		ProjectName: strings.ToLower(s.T().Name()),
	}

	project := NewDockerComposeProject(options)
	command := NewDockerComposeCommand(s.T(), &project)
	s.stack = NewDockerComposeStack(
		WithCommand(&command),
		WithStackPreUpAction("build"),
		WithStackPreUpAction("run --rm openssh -t rsa -q -N password -f id_rsa -C \"apenella@stevedore.test\""),
		WithStackPreUpAction("run --rm openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config /root/ssl/stevedore.test.cnf"),
		WithStackPreUpAction("run --rm openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile /root/ssl/stevedore.test.cnf"),
		WithStackPostUpAction("run stevedore /prepare-images"),
	)

	err = s.stack.DownAndUp("-d registry docker-hub gitserver stevedore")
	if err != nil {
		defer s.TearDownSuite()
		s.T().Log(err)
		s.T().FailNow()
	}
}

func (s *FunctionalTestsSuite) TearDownSuite() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = s.stack.Down()
	if err != nil {
		s.T().Log(err)
		s.T().FailNow()
	}
}
