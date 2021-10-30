package promote

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (p *ServiceMock) Promote(ctx context.Context, options *ServiceOptions, promoteType string) error {
	args := p.Mock.Called(ctx, options, promoteType)
	return args.Error(0)
}
