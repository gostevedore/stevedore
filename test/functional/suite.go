package functional

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FunctionalTestsSuite struct {
	suite.Suite
	stack            *DockerComposeStack
	setupSuiteFunc   func(*testing.T, *DockerComposeStack) error
	tearDownSuitFunc func(*testing.T, *DockerComposeStack) error
}

type OptionsFunc func(*FunctionalTestsSuite)

func NewTestSuite(opts ...OptionsFunc) *FunctionalTestsSuite {
	suite := new(FunctionalTestsSuite)

	suite.Options(opts...)

	return suite
}

func defaultSetupSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	return errors.New("You are using the default setup suite function")
}

func defaultTearDownSuiteFunc(t *testing.T, stack *DockerComposeStack) error {
	return errors.New("You are using the default tear down suite function")
}

func WithStack(stack *DockerComposeStack) OptionsFunc {
	return func(s *FunctionalTestsSuite) {
		s.stack = stack
	}
}

func WithSetupSuiteFunc(f func(*testing.T, *DockerComposeStack) error) OptionsFunc {
	return func(s *FunctionalTestsSuite) {
		s.setupSuiteFunc = f
	}
}

func WithTearDownSuiteFunc(f func(*testing.T, *DockerComposeStack) error) OptionsFunc {
	return func(s *FunctionalTestsSuite) {
		s.tearDownSuitFunc = f
	}
}

func (s *FunctionalTestsSuite) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *FunctionalTestsSuite) SetupSuite() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	if s.setupSuiteFunc == nil {
		err = defaultSetupSuiteFunc(s.T(), s.stack)
	} else {
		err = s.setupSuiteFunc(s.T(), s.stack)
	}

	if err != nil {
		s.TearDownSuite()
		s.T().Log(err)
		s.T().FailNow()
	}
}

func (s *FunctionalTestsSuite) TearDownSuite() {
	var err error

	if testing.Short() {
		s.T().Skip("functional test are skipped in short mode")
	}

	if s.stack != nil {
		err = s.stack.Down()
		if err != nil {
			s.T().Log(err)
			s.T().FailNow()
		}
	}
}
