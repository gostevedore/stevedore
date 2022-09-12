package credentials

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockOutput is an output for the builders
type MockOutput struct {
	mock.Mock
}

// NewMockOutput creates a new MockOutput
func NewMockOutput() *MockOutput {
	return &MockOutput{}
}

// MockOutput prints the credentials
func (o *MockOutput) Print(badges []*credentials.Badge) error {
	args := o.Mock.Called(badges)
	return args.Error(0)
}
