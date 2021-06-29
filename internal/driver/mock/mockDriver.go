package mockdriver

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/types"
)

const (
	DriverName = "mock"
)

// MockBuilder
type MockDriver struct{}

func NewMockDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
	return &MockDriver{}, nil
}

func (b *MockDriver) Run(ctx context.Context) error {
	return nil
}

// MockBuilderRunErr
type MockDriverErr struct{}

func NewMockDriverErr(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
	return &MockDriverErr{}, nil
}

func (b *MockDriverErr) Run(ctx context.Context) error {
	return errors.New("(MockDriverRunErr)", "Error")
}

// NewMockDrivererOnNew
func NewMockDriverErrOnNew(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
	return nil, errors.New("(NewMockDriverErrOnNew)", "Error")
}
