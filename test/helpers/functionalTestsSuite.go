package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

// FunctionalTestsSuite is a struct that defines the functional test suite
type FunctionalTestsSuite struct {
	suite.Suite
	CommandFactory   CommandFactorier
	SetupSuiteFunc   func() error
	Stack            *DockerComposeStack
	TearDownSuitFunc func() error
}

// OptionsFunc is a function used to configure the FunctionalTestsSuite
type OptionsFunc func(*FunctionalTestsSuite)

// NewFunctionalTestsSuite creates a new FunctionalTestsSuite
func NewFunctionalTestsSuite(opts ...OptionsFunc) *FunctionalTestsSuite {
	suite := new(FunctionalTestsSuite)

	suite.Options(opts...)

	return suite
}

// defaultSetupSuiteFunc is the default function to execute in the SetupSuite
func defaultSetupSuiteFunc() error {
	return errors.New("You need to define a function to setup the test suite")
}

// defaultTearDownSuiteFunc is the default function to execute in the TearDownSuite
func defaultTearDownSuiteFunc() error {
	return errors.New("You need to define a function to tear down the test suite")
}

// WithStack sets the stack for the FunctionalTestsSuite
func WithStack(stack *DockerComposeStack) OptionsFunc {
	return func(suite *FunctionalTestsSuite) {
		suite.Stack = stack
	}
}

// WithSetupSuiteFunc sets the setup suite function for the FunctionalTestsSuite
func WithSetupSuiteFunc(f func() error) OptionsFunc {
	return func(suite *FunctionalTestsSuite) {
		suite.SetupSuiteFunc = f
	}
}

// WithTearDownSuiteFunc sets the tear down suite function for the FunctionalTestsSuite
func WithTearDownSuiteFunc(f func() error) OptionsFunc {
	return func(suite *FunctionalTestsSuite) {
		suite.TearDownSuitFunc = f
	}
}

// WithCommandFactory sets the command factory for the FunctionalTestsSuite
func WithCommandFactory(factory CommandFactorier) OptionsFunc {
	return func(suite *FunctionalTestsSuite) {
		suite.CommandFactory = factory
	}
}

// Options sets options for the FunctionalTestsSuite
func (suite *FunctionalTestsSuite) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(suite)
	}
}

// SetupSuite runs before the tests and setups the test suite
func (suite *FunctionalTestsSuite) SetupSuite() {
	var err error

	if testing.Short() {
		suite.T().Skip("functional test are skipped in short mode")
	}

	if suite.SetupSuiteFunc == nil {
		err = defaultSetupSuiteFunc()
	} else {
		err = suite.SetupSuiteFunc()
	}

	if err != nil {
		suite.TearDownSuite()
		suite.T().Log(err)
		suite.T().FailNow()
	}
}

// TearDownSuite runs after the tests and tears down the test suite
func (suite *FunctionalTestsSuite) TearDownSuite() {
	var err error

	if testing.Short() {
		suite.T().Skip("functional test are skipped in short mode")
	}

	if suite.SetupSuiteFunc == nil {
		err = defaultTearDownSuiteFunc()
	} else {
		err = suite.TearDownSuitFunc()
	}

	if suite.Stack != nil {
		err = suite.Stack.Down()
		if err != nil {
			suite.T().Log(err)
			suite.T().FailNow()
		}
	}
}
