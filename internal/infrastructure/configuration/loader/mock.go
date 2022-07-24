package loader

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
)

type MockConfigurationLoader struct {
	mock.Mock
}

func NewMockConfigurationLoader() *MockConfigurationLoader {
	return &MockConfigurationLoader{}
}

// AddConfigPath adds a path for Viper to search for the config file in. Can be called multiple times to define multiple search paths
func (c *MockConfigurationLoader) AddConfigPath(in string) {
	c.Called(in)
}

// AutomaticEnv has Viper check ENV variables for all. keys set in config, default & flags
func (c *MockConfigurationLoader) AutomaticEnv() {
	c.Called()
}

// GetBool returns the value associated with the key as a boolean
func (c *MockConfigurationLoader) GetBool(key string) bool {
	args := c.Called(key)
	return args.Bool(0)
}

// GetInt returns the value associated with the key as an integer
func (c *MockConfigurationLoader) GetInt(key string) int {
	args := c.Called(key)
	return args.Int(0)
}

// GetString returns the value associated with the key as a string
func (c *MockConfigurationLoader) GetString(key string) string {
	args := c.Called(key)
	return args.String(0)
}

// GetStringSlice returns the value associated with the key as a slice of strings
func (c *MockConfigurationLoader) GetStringSlice(key string) []string {
	args := c.Called(key)
	return args.Get(0).([]string)
}

// ReadInConfig will discover and load the configuration file from disk and key/value stores, searching in one of the defined paths
func (c *MockConfigurationLoader) ReadInConfig() error {
	args := c.Called()
	return args.Error(0)
}

// SetConfigFile explicitly defines the path, name and extension of the config file. Viper will use this and not check any of the config paths
func (c *MockConfigurationLoader) SetConfigFile(in string) {
	c.Called(in)
}

// SetConfigName sets name for the config file. Does not include extension
func (c *MockConfigurationLoader) SetConfigName(in string) {
	c.Called(in)
}

// SetConfigType sets the type of the configuration returned by the remote source, e.g. "json"
func (c *MockConfigurationLoader) SetConfigType(in string) {
	c.Called(in)
}

// SetDefault sets the default value for this key. SetDefault is case-insensitive for a key. Default only used when no value is provided by the user via flag, config or ENV
func (c *MockConfigurationLoader) SetDefault(key string, value interface{}) {
	c.Called(key, value)
}

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use. E.g. if your prefix is "spf", the env registry will look for env variables that start with "SPF_"
func (c *MockConfigurationLoader) SetEnvPrefix(in string) {
	c.Called(in)
}

// SetFs sets the filesystem to use to read configuration
func (c *MockConfigurationLoader) SetFs(fs afero.Fs) {
	c.Called(fs)
}

// ConfigFileUsed sets the filesystem to use to read configuration
func (c *MockConfigurationLoader) ConfigFileUsed() string {
	args := c.Called()
	return args.String(0)
}
