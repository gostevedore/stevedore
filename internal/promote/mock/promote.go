package mockpromote

import (
	"context"

	"github.com/gostevedore/stevedore/internal/promote"
	"github.com/stretchr/testify/mock"
)

type MockPromote struct {
	mock.Mock
}

func NewMockPromote() *MockPromote {
	return &MockPromote{}
}

func (p *MockPromote) Promote(ctx context.Context, options *promote.PromoteOptions) error {
	args := p.Mock.Called(ctx, options)
	return args.Error(0)
}
