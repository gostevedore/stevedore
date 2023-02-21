package promote

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockApplication struct {
	mock.Mock
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (p *MockApplication) Promote(ctx context.Context, options *Options) error {
	args := p.Mock.Called(ctx, options)
	return args.Error(0)
}
