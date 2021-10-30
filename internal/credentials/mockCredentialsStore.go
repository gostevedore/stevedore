package credentials

import "github.com/stretchr/testify/mock"

type CredentialsStoreMock struct {
	mock.Mock
}

func NewCredentialsStoreMock() *CredentialsStoreMock {
	return &CredentialsStoreMock{}
}

func (m *CredentialsStoreMock) GetCredentials(registry string) (*RegistryUserPassAuth, error) {
	args := m.Mock.Called(registry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*RegistryUserPassAuth), args.Error(1)
	}
}
