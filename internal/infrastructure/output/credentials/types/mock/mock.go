package mock

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockOutput is a mock implementation of the Outputter interface
type MockOutput struct {
	mock.Mock
}

// NewMockOutput creates a new MockOutput
func NewMockOutput() *MockOutput {
	return &MockOutput{}
}

// Output is a mock implementation of the Outputter interface
func (o *MockOutput) Output(credential *credentials.Credential) (string, string, error) {
	args := o.Called(credential)
	return args.String(0), args.String(1), args.Error(2)
}
