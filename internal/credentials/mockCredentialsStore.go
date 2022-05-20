package credentials

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

type CredentialsStoreMock struct {
	mock.Mock
}

func NewCredentialsStoreMock() *CredentialsStoreMock {
	return &CredentialsStoreMock{}
}

func (m *CredentialsStoreMock) Get(id string) (*credentials.UserPasswordAuth, error) {
	args := m.Mock.Called(id)

	// It is used when you need to return a nil UserPasswordAuth
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*credentials.UserPasswordAuth), args.Error(1)
	}
}
