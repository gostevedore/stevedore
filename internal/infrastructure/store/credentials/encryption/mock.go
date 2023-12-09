package encryption

import "github.com/stretchr/testify/mock"

type MockEncription struct {
	mock.Mock
}

func NewMockEncryption() *MockEncription {
	return &MockEncription{}
}

func (e *MockEncription) Encrypt(text string) (string, error) {
	args := e.Mock.Called(text)
	return args.String(0), args.Error(1)
}

func (e *MockEncription) Decrypt(ciphertext string) (string, error) {
	args := e.Mock.Called(ciphertext)
	return args.String(0), args.Error(1)
}

func (e *MockEncription) GenerateEncryptionKey() (string, error) {
	args := e.Mock.Called()
	return args.String(0), args.Error(1)
}
