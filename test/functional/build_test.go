package functional

import (
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/pkg/terratest/modules/docker"
	"github.com/stretchr/testify/suite"
)

type BuildFunctionalTestsSuite struct {
	suite.Suite
	options *docker.Options
}

func (s *BuildFunctionalTestsSuite) SetupSuite() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	s.options = &docker.Options{
		WorkingDir:  ".",
		ProjectName: strings.ToLower(s.T().Name()),
	}

	err = start(s.T(), s.options)
	if err != nil {
		defer s.TearDownSuite()
		s.T().Log(err)
		s.T().FailNow()
	}
}

func (s *BuildFunctionalTestsSuite) TearDownSuite() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = dockerComposeCommand(s.T(), s.options, "down -v --remove-orphans")
	if err != nil {
		s.T().Log(err)
		s.T().FailNow()
	}
}

func (s *BuildFunctionalTestsSuite) TestBuildImageWithGitContext() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore stevedore build app2 --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
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

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore stevedore build app1 --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-20.04")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-scratch-latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
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

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore stevedore build alpine --build-on-cascade --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16")
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

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore stevedore build app3 --enable-semver-tags --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2.3")
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

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore stevedore build app3 --image-version 1.3.0-rc0.1+1234 --push-after-build --enable-semver-tags")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1-rc0.1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3-rc0.1-1234")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = dockerComposeCommand(s.T(), s.options, "run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3.0-rc0.1_1234")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func dockerComposeCommand(t *testing.T, options *docker.Options, cmd string) error {
	var err error
	cmds := strings.Split(cmd, " ")

	_, err = docker.RunDockerComposeE(t, options, cmds...)
	if err != nil {
		return err
	}
	return nil
}

func generateKeys(t *testing.T, options *docker.Options) error {
	var err error

	err = dockerComposeCommand(t, options, "run --rm openssh -t rsa -q -N password -f id_rsa -C \"apenella@stevedore.test\"")

	return err
}

func generateCerts(t *testing.T, options *docker.Options) error {

	var err error

	err = dockerComposeCommand(t, options, "run --rm openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config /root/ssl/stevedore.test.cnf")
	if err != nil {
		return err
	}

	err = dockerComposeCommand(t, options, "run --rm openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile /root/ssl/stevedore.test.cnf")
	if err != nil {
		return err
	}

	return nil
}

func start(t *testing.T, options *docker.Options) error {
	var err error

	err = dockerComposeCommand(t, options, "down -v --remove-orphans")
	if err != nil {
		return err
	}

	err = dockerComposeCommand(t, options, "build")
	if err != nil {
		return err
	}

	err = generateKeys(t, options)
	if err != nil {
		return err
	}
	err = generateCerts(t, options)
	if err != nil {
		return err
	}

	err = dockerComposeCommand(t, options, "up -d registry docker-hub gitserver stevedore")
	if err != nil {
		return err
	}

	err = dockerComposeCommand(t, options, "run stevedore /prepare-images.sh")
	if err != nil {
		return err
	}

	return nil
}

func TestBuildFunctionalTests(t *testing.T) {
	suite.Run(t, new(BuildFunctionalTestsSuite))
}
