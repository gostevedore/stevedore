package backend

import (
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/fatih/structs"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/spf13/viper"
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

func (b *OSEnvvarsBackend) AchieveBadge(id string) (*credentials.Badge, error) {
	errContext := "(store::credentials::envvars::backend::OSEnvvarsBackend::AchieveBadge)"

	badge := &credentials.Badge{}
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(id)

	badgeFields := structs.Fields(badge)
	for _, field := range badgeFields {
		attribute := field.Tag("mapstructure")
		if attribute == "" {
			continue
		}

		v.BindEnv(attribute)
	}

	v.ReadInConfig()
	err := v.Unmarshal(badge)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return badge, nil
}
