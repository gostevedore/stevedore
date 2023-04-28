package mock

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

type MockFormater struct {
	mock.Mock
}

func NewMockFormater() *MockFormater {
	return &MockFormater{}
}

func (f *MockFormater) Marshal(badge *credentials.Credential) (string, error) {
	args := f.Called(badge)
	return args.String(0), args.Error(1)
}

func (f *MockFormater) Unmarshal(data []byte) (*credentials.Credential, error) {
	args := f.Called(data)
	return args.Get(0).(*credentials.Credential), args.Error(1)
}
