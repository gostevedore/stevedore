package images

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/mock"
)

type MockOutput struct {
	mock.Mock
}

func NewMockOutput() *MockOutput {
	return &MockOutput{}
}

func (o *MockOutput) Output(list []*image.Image) error {
	args := o.Mock.Called(list)
	return args.Error(0)
}
