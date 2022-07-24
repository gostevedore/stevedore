package loader

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type ConfigurationLoader struct {
	viper *viper.Viper
}

func NewConfigurationLoader(viper *viper.Viper) *ConfigurationLoader {
	return &ConfigurationLoader{viper: viper}
}

// AddConfigPath AddConfigPath adds a path for Viper to search for the config file in. Can be called multiple times to define multiple search paths
func (c *ConfigurationLoader) AddConfigPath(in string) {
	c.viper.AddConfigPath(in)
}

// AutomaticEnv has Viper check ENV variables for all. keys set in config, default & flags
func (c *ConfigurationLoader) AutomaticEnv() {
	c.viper.AutomaticEnv()
}

// GetBool returns the value associated with the key as a boolean
func (c *ConfigurationLoader) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetInt returns the value associated with the key as an integer
func (c *ConfigurationLoader) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetString returns the value associated with the key as a string
func (c *ConfigurationLoader) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings
func (c *ConfigurationLoader) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// ReadInConfig will discover and load the configuration file from disk and key/value stores, searching in one of the defined paths
func (c *ConfigurationLoader) ReadInConfig() error {
	return c.viper.ReadInConfig()
}

// SetConfigFile explicitly defines the path, name and extension of the config file. Viper will use this and not check any of the config paths
func (c *ConfigurationLoader) SetConfigFile(in string) {
	c.viper.SetConfigFile(in)
}

// SetConfigName sets name for the config file. Does not include extension
func (c *ConfigurationLoader) SetConfigName(in string) {
	c.viper.SetConfigName(in)
}

// SetConfigType sets the type of the configuration returned by the remote source, e.g. "json"
func (c *ConfigurationLoader) SetConfigType(in string) {
	c.viper.SetConfigType(in)
}

// SetDefault sets the default value for this key. SetDefault is case-insensitive for a key. Default only used when no value is provided by the user via flag, config or ENV
func (c *ConfigurationLoader) SetDefault(key string, value interface{}) {
	c.viper.SetDefault(key, value)
}

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use. E.g. if your prefix is "spf", the env registry will look for env variables that start with "SPF_"
func (c *ConfigurationLoader) SetEnvPrefix(in string) {
	c.viper.SetEnvPrefix(in)
}

// SetFs sets the filesystem to use to read configuration
func (c *ConfigurationLoader) SetFs(fs afero.Fs) {
	c.viper.SetFs(fs)
}

// ConfigFileUsed sets the filesystem to use to read configuration
func (c *ConfigurationLoader) ConfigFileUsed() string {
	return c.viper.ConfigFileUsed()
}
