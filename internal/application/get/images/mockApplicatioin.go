package images

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGetImagesApplication is a mock of build application
type MockGetImagesApplication struct {
	mock.Mock
}

// NewMockGetImagesApplication return a mock of build application
func NewMockGetImagesApplication() *MockGetImagesApplication {
	return &MockGetImagesApplication{}
}

// Build provides a mock function with given fields: ctx, buildPlan, name, version, options, optionsFunc
func (m *MockGetImagesApplication) Run(ctx context.Context, options *Options, optionsFunc ...OptionsFunc) error {
	args := m.Called(ctx, options, optionsFunc)
	return args.Error(0)
}
