package builders

import (
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/stretchr/testify/mock"
)

type MockOutput struct {
	mock.Mock
}

func NewMockOutput() *MockOutput {
	return &MockOutput{}
}

func (o *MockOutput) Output(list []*builder.Builder) error {
	args := o.Mock.Called(list)
	return args.Error(0)
}
