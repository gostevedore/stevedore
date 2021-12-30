package mockdriver

import (
	"context"

	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/stretchr/testify/mock"
)

// const (
// 	DriverName = "mock"
// )

// MockDriver is a mock implementation of driver.BuildDriverer
type MockDriver struct {
	mock.Mock
}

// NewMockDriver creates a new MockDriver
func NewMockDriver() *MockDriver {
	return &MockDriver{}
}

// Build simulate a new image build
func (d *MockDriver) Build(ctx context.Context, options *driver.BuildDriverOptions) error {
	args := d.Called(ctx, options)
	return args.Error(0)
}
