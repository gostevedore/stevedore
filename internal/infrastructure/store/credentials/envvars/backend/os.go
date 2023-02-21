package backend

import (
	"os"
)

// OSEnvvarsBackend is a backend for envvars that uses the OS envvars
type OSEnvvarsBackend struct{}

func NewOSEnvvarsBackend() *OSEnvvarsBackend {
	return &OSEnvvarsBackend{}
}

// Setenv sets the value of the environment variable named by the key. It returns an error, if any
func (b *OSEnvvarsBackend) Setenv(key, value string) {
	os.Setenv(key, value)
}

// Unsetenv unsets a single environment variable
func (b *OSEnvvarsBackend) Unsetenv(key string) {
	os.Unsetenv(key)
}

// Getenv retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present. To distinguish between an empty value and an unset value, use LookupEnv
func (b *OSEnvvarsBackend) Getenv(key string) string {
	return os.Getenv(key)
}

// LookupEnv retrieves the value of the environment variable named by the key. If the variable is set the value (which may be empty) is returned and the boolean is true. Otherwise the returned value will be empty and the boolean will be false
func (b *OSEnvvarsBackend) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Environ returns a copy of strings representing the environment, in the form "key=value"
func (b *OSEnvvarsBackend) Environ() []string {
	return os.Environ()
}
