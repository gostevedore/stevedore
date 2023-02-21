package job

import (
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/stretchr/testify/mock"
)

// MockJobFactory is a factory to create a mock build command
type MockJobFactory struct {
	mock.Mock
}

// NewMockJobFactory returns a new mock build command factory
func NewMockJobFactory() *MockJobFactory {
	return &MockJobFactory{}
}

// New returns a new build command constructor
func (f *MockJobFactory) New(command Commander) scheduler.Jobber {
	args := f.Called(command)
	return args.Get(0).(scheduler.Jobber)
}
