package configuration

import (
	"strings"

	"github.com/spf13/afero"
)

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}

// ConfigurationLoader is the interface for the configuration loader
type ConfigurationLoader interface {
	AddConfigPath(in string)
	AutomaticEnv()
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetStringSlice(key string) []string
	ReadInConfig() error
	SetConfigFile(in string)
	SetConfigName(in string)
	SetConfigType(in string)
	SetDefault(key string, value interface{})
	SetEnvPrefix(in string)
	SetFs(fs afero.Fs)
	ConfigFileUsed() string
	SetEnvKeyReplacer(*strings.Replacer)
}

type ConfigurationWriter interface {
	Write(config *Configuration) error
}
