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

func (p *MockService) Promote(ctx context.Context, options *ServiceOptions, promoteType string) error {
	args := p.Mock.Called(ctx, options, promoteType)
	return args.Error(0)
}
