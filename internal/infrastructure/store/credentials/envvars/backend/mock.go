package backend

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockEnvvarsBackend is a backend for envvars that uses the OS envvars
type MockEnvvarsBackend struct {
	mock.Mock
}

func NewMockEnvvarsBackend() *MockEnvvarsBackend {
	return &MockEnvvarsBackend{}
}

// Setenv sets the value of the environment variable named by the key. It returns an error, if any
func (b *MockEnvvarsBackend) Setenv(key, value string) {
	b.Mock.Called(key, value)
}

// Getenv retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present. To distinguish between an empty value and an unset value, use LookupEnv
func (b *MockEnvvarsBackend) Getenv(key string) string {
	args := b.Mock.Called(key)
	return args.Get(0).(string)
}

// LookupEnv retrieves the value of the environment variable named by the key. If the variable is set the value (which may be empty) is returned and the boolean is true. Otherwise the returned value will be empty and the boolean will be false
func (b *MockEnvvarsBackend) LookupEnv(key string) (string, bool) {
	args := b.Mock.Called(key)
	return args.String(0), args.Bool(1)
}

// Environ returns a copy of strings representing the environment, in the form "key=value"
func (b *MockEnvvarsBackend) Environ() []string {
	args := b.Mock.Called()
	return args.Get(0).([]string)
}

func (b *MockEnvvarsBackend) AchieveBadge(id string) (*credentials.Badge, error) {
	args := b.Mock.Called(id)
	return args.Get(0).(*credentials.Badge), args.Error(1)
}
