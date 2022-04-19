package promote

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func NewMockService() *MockService {
	return &MockService{}
}

func (p *MockService) Promote(ctx context.Context, options *ServiceOptions) error {
	args := p.Mock.Called(ctx, options)
	return args.Error(0)
}
