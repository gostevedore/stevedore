package functional

import (
	"testing"

	// "github.com/gostevedore/stevedore/pkg/terratest/modules/docker"
	"github.com/stretchr/testify/suite"
)

type BuildFunctionalTestsSuite struct {
	FunctionalTestsSuite
}

func (s *BuildFunctionalTestsSuite) SetupTest() {
	s.TearDownTest()

	err := s.stack.Execute("up -d registry")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func (s *BuildFunctionalTestsSuite) TearDownTest() {
	err := s.stack.Execute("rm --stop --force --volumes registry stevedore")
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

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore stevedore build app2 --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16")
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

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore stevedore build app1 --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-ubuntu-20.04")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-scratch-latest")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
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

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore stevedore build alpine --build-on-cascade --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app1:v1-alpine-3.16")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app2:v1-alpine-3.16")
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

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore stevedore build app3 --enable-semver-tags --push-after-build")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.2.3")
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

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore stevedore build app3 --image-version 1.3.0-rc0.1+1234 --push-after-build --enable-semver-tags")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1-rc0.1")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3-rc0.1-1234")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}

	err = s.stack.Execute("run -w /app/test/stack/client/stevedore stevedore docker pull registry.stevedore.test/stable/app3:1.3.0-rc0.1_1234")
	if err != nil {
		s.T().Log(err)
		s.T().Fail()
	}
}

func TestBuildFunctionalTests(t *testing.T) {
	suite.Run(t, new(BuildFunctionalTestsSuite))
}
