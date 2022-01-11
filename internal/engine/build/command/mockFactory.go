package command

import (
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/stretchr/testify/mock"
)

// MockBuildCommandFactory is a factory to create a mock build command
type MockBuildCommandFactory struct {
	mock.Mock
}

// NewMockBuildCommandFactory returns a new mock build command factory
func NewMockBuildCommandFactory() *MockBuildCommandFactory {
	return &MockBuildCommandFactory{}
}

// New returns a new build command constructor
func (f *MockBuildCommandFactory) New(driver driver.BuildDriverer, options *driver.BuildDriverOptions) BuildCommander {
	args := f.Called(driver, options)
	return args.Get(0).(BuildCommander)
}
