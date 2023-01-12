package configuration

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/spf13/afero"
)

type CredentialsConfiguration struct {
	// StorageType is the backend used to store credentials
	StorageType string
	// LocalStoragePath is the local storage path where credentials are stored
	LocalStoragePath string
	// Format defines the format to store credentials, in case a format is required
	Format string
	// EncryptionKey is the key used to encrypt credentials
	EncryptionKey string
}

type Configuration struct {
	// BuildersPath is the path where the builders are stored
	BuildersPath string
	// Concurrency is the number of concurrent builds
	Concurrency int
	// Credentials is the credentials configuration block
	Credentials *CredentialsConfiguration
	// DEPRECATEDBuilderPath is the path where the builders are stored
	DEPRECATEDBuilderPath string
	// DEPRECATEDBuildOnCascade is the flag to build on cascade
	DEPRECATEDBuildOnCascade bool
	// DEPRECATEDDockerCredentialsDir is the path to the docker credentials directory
	DEPRECATEDDockerCredentialsDir string
	// DEPRECATEDNumWorkers is the number of concurrent workers
	DEPRECATEDNumWorkers int
	// DEPRECATEDTreePathFile is the path to the tree path file
	DEPRECATEDTreePathFile string
	// EnableSemanticVersionTags is the flag to enable semantic version tags
	EnableSemanticVersionTags bool
	// ImagesPath is the path where the images are stored
	ImagesPath string
	// LogPathFile is the path to the log file
	LogPathFile string
	// LogWriter is the writer to the log file
	LogWriter io.Writer
	// PushImages is the flag to push images automatically after build
	PushImages bool
	// SemanticVersionTagsTemplates is the list of semantic version tags templates
	SemanticVersionTagsTemplates []string

	compatibility Compatibilitier
	configFile    string
	fs            afero.Fs
	loader        ConfigurationLoader
}

const (

	// DefaultConfigFile is the name of the default configuration file
	DefaultConfigFile = "./stevedore"
	// DefaultConfigFileExtention is the default configuration file extention
	DefaultConfigFileExtention = "yaml"
	// DefaultConfigFolder is the default configuration folder
	DefaultConfigFolder = "."

	// DefaultBuildersPath is the default builders path
	DefaultBuildersPath = "stevedore.yaml"
	// DefaultCredentialsFormat is the default credentials format
	DefaultCredentialsFormat = credentials.JSONFormat
	// DefaultCredentialsLocalStoragePath is the default credentials local storage path
	DefaultCredentialsLocalStoragePath = "credentials"
	// DefaultCredentialsStorage is the default credentials storage
	DefaultCredentialsStorage = credentials.LocalStore
	// DefaultEnableSemanticVersionTags is the default enable semantic version tags
	DefaultEnableSemanticVersionTags = false
	// DefaultImagesPath is the default images path
	DefaultImagesPath = "stevedore.yaml"
	// DefaultLogPathFile is the default log path file
	DefaultLogPathFile = ""
	// DefaultPushImages by default images won't be pushed
	DefaultPushImages = false
	// DefaultSemanticVersionTagsTemplates is the default semantic version tags templates
	DefaultSemanticVersionTagsTemplates = "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
	// DEPRECATEDDefaultBuilderPath is the default builder path
	DEPRECATEDDefaultBuilderPath = "stevedore.yaml"
	// DEPRECATEDDefaultBuildOnCascade
	DEPRECATEDDefaultBuildOnCascade = false
	// DEPRECATEDDefaultDockerCredentialsDir
	DEPRECATEDDefaultDockerCredentialsDir = "credentials"
	// DEPRECATEDDefaultNumWorker
	DEPRECATEDDefaultNumWorker = 4
	// DEPRECATEDDefaultTreePathFile
	DEPRECATEDDefaultTreePathFile = "stevedore.yaml"

	// BuildersPathKey is the key for the builders path
	BuildersPathKey = "builders_path"
	// ConcurrencyKey is the key for the concurrency value
	ConcurrencyKey = "concurrency"
	// CredentialsFormatKey is the key for the credentials format
	CredentialsFormatKey = "format"
	// CredentialsKey is the key for the credentials block
	CredentialsKey = "credentials"
	// CredentialsLocalStoragePathKey is the key for the credentials local storage path
	CredentialsLocalStoragePathKey = "local_storage_path"
	// CredentialsEncryptionKeyKey is the key for the credentials encryption token
	CredentialsEncryptionKeyKey = "encryption_key"
	// CredentialsStorageTypeKey is the key for the credentials storage type
	CredentialsStorageTypeKey = "storage_type"
	// DEPRECATEDBuilderPathKey is the key for the deprecated builder path
	DEPRECATEDBuilderPathKey = "builder_path"
	// DEPRECATEDBuildOnCascadeKey is the key for the deprecated build on cascade value
	DEPRECATEDBuildOnCascadeKey = "build_on_cascade"
	// DEPRECATEDDockerCredentialsDirKey is the key for the deprecated docker credentials dir
	DEPRECATEDDockerCredentialsDirKey = "docker_registry_credentials_dir"
	// DEPRECATEDNumWorkerKey is the key for the deprecated number of workers
	DEPRECATEDNumWorkerKey = "num_workers"
	// DEPRECATEDTreePathFileKey is the key for the deprecated tree path file
	DEPRECATEDTreePathFileKey = "tree_path"
	// EnableSemanticVersionTagsKey is the key for the enable semantic version tags value
	EnableSemanticVersionTagsKey = "semantic_version_tags_enabled"
	// ImagesPathKey is the key for the images path
	ImagesPathKey = "images_path"
	// LogPathFileKey is the key for the log path file
	LogPathFileKey = "log_path"
	// PushImagesKey is the key for the push images value
	PushImagesKey = "push_images"
	// SemanticVersionTagsTemplatesKey is the key for the semantic version tags templates
	SemanticVersionTagsTemplatesKey = "semantic_version_tags_templates"
)

func DefaultConfig() *Configuration {

	config := &Configuration{}

	// dynamic default values
	defaultConcurrency := concurrencyValue()

	config.BuildersPath = filepath.Join(DefaultConfigFolder, DefaultBuildersPath)
	config.Concurrency = defaultConcurrency
	config.EnableSemanticVersionTags = DefaultEnableSemanticVersionTags
	config.ImagesPath = filepath.Join(DefaultConfigFolder, DefaultImagesPath)
	config.LogWriter = io.Discard
	config.PushImages = DefaultPushImages
	config.SemanticVersionTagsTemplates = []string{DefaultSemanticVersionTagsTemplates}

	config.Credentials = &CredentialsConfiguration{
		StorageType:      DefaultCredentialsStorage,
		LocalStoragePath: DefaultCredentialsLocalStoragePath,
		Format:           DefaultCredentialsFormat,
	}

	return config
}

// New method create a new configuration object
func New(fs afero.Fs, loader ConfigurationLoader, compatibility Compatibilitier) (*Configuration, error) {

	var logWriter io.Writer

	errContext := "(Configuration::New)"

	if compatibility == nil {
		return nil, errors.New(errContext, "Comptabilitier must be provided to create a new configuration")
	}

	if fs == nil {
		return nil, errors.New(errContext, "File system must be provided to create a new configuration")
	}

	if loader == nil {
		return nil, errors.New(errContext, "Configuration loader must be provided to create a new configuration")
	}

	user, err := user.Current()
	if err != nil {
		return nil, errors.New(errContext, "Current user information can not be cached", err)
	}

	alternativesConfigFolders := []string{
		filepath.Join(user.HomeDir, ".config", "stevedore"),
		user.HomeDir,
	}

	config := &Configuration{
		fs:            fs,
		loader:        loader,
		compatibility: compatibility,
	}

	loader.SetFs(fs)

	loader.AutomaticEnv()
	loader.SetEnvPrefix("stevedore")

	loader.SetConfigName(DefaultConfigFile)
	loader.SetConfigType(DefaultConfigFileExtention)

	// dynamic default values
	defaultConcurrency := concurrencyValue()

	loader.SetDefault(BuildersPathKey, filepath.Join(DefaultConfigFolder, DefaultBuildersPath))
	loader.SetDefault(ConcurrencyKey, defaultConcurrency)
	loader.SetDefault(EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags)
	loader.SetDefault(ImagesPathKey, filepath.Join(DefaultConfigFolder, DefaultImagesPath))
	loader.SetDefault(LogPathFileKey, DefaultLogPathFile)
	loader.SetDefault(PushImagesKey, DefaultPushImages)
	loader.SetDefault(SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates})
	loader.SetDefault(
		strings.Join([]string{CredentialsKey, CredentialsStorageTypeKey}, "."), DefaultCredentialsStorage)
	loader.SetDefault(
		strings.Join([]string{CredentialsKey, CredentialsLocalStoragePathKey}, "."), DefaultCredentialsLocalStoragePath)
	loader.SetDefault(
		strings.Join([]string{CredentialsKey, CredentialsFormatKey}, "."), DefaultCredentialsFormat)

	for _, alternativeConfigFolder := range alternativesConfigFolders {
		loader.AddConfigPath(alternativeConfigFolder)
	}
	// Set the default config folder as last default option
	loader.AddConfigPath(DefaultConfigFolder)

	// when configuration is created no error is shown if readinconfig files. It will use the defaults
	loader.ReadInConfig()

	logWriter, err = createLogWriter(fs, loader.GetString(LogPathFileKey))
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if loader.GetInt(DEPRECATEDNumWorkerKey) > 0 {
		config.DEPRECATEDNumWorkers = loader.GetInt(DEPRECATEDNumWorkerKey)
	}

	if loader.GetString(DEPRECATEDTreePathFileKey) != "" {
		config.DEPRECATEDTreePathFile = loader.GetString(DEPRECATEDTreePathFileKey)
	}

	if loader.GetString(DEPRECATEDBuilderPathKey) != "" {
		config.DEPRECATEDBuilderPath = loader.GetString(DEPRECATEDBuilderPathKey)
	}

	if loader.GetBool(DEPRECATEDBuildOnCascadeKey) {
		config.DEPRECATEDBuildOnCascade = loader.GetBool(DEPRECATEDBuildOnCascadeKey)
	}

	if loader.GetString(DEPRECATEDDockerCredentialsDirKey) != "" {
		config.DEPRECATEDDockerCredentialsDir = loader.GetString(DEPRECATEDDockerCredentialsDirKey)
	}

	config.BuildersPath = loader.GetString(BuildersPathKey)
	config.Concurrency = loader.GetInt(ConcurrencyKey)
	config.EnableSemanticVersionTags = loader.GetBool(EnableSemanticVersionTagsKey)
	config.ImagesPath = loader.GetString(ImagesPathKey)
	config.LogWriter = logWriter
	config.PushImages = loader.GetBool(PushImagesKey)
	config.SemanticVersionTagsTemplates = loader.GetStringSlice(SemanticVersionTagsTemplatesKey)

	config.Credentials = &CredentialsConfiguration{
		StorageType:      loader.GetString(strings.Join([]string{CredentialsKey, CredentialsStorageTypeKey}, ".")),
		LocalStoragePath: loader.GetString(strings.Join([]string{CredentialsKey, CredentialsLocalStoragePathKey}, ".")),
		Format:           loader.GetString(strings.Join([]string{CredentialsKey, CredentialsFormatKey}, ".")),
	}

	config.configFile = loader.ConfigFileUsed()

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	err = config.ValidateConfiguration()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return config, nil
}

// LoadFromFile method returns a configuration object loaded from a file
func LoadFromFile(fs afero.Fs, loader ConfigurationLoader, file string, compatibility Compatibilitier) (*Configuration, error) {

	var err error
	var logWriter io.Writer
	var config *Configuration

	errContext := "(configuration::LoadFromFile)"

	if compatibility == nil {
		return nil, errors.New(errContext, "Comptabilitier must be provided to create a new configuration")
	}

	if fs == nil {
		return nil, errors.New(errContext, "File system must be provided to create a new configuration")
	}

	if loader == nil {
		return nil, errors.New(errContext, "Configuration loader must be provided to create a new configuration")
	}

	loader.SetFs(fs)
	loader.SetConfigFile(file)
	err = loader.ReadInConfig()
	if err != nil {
		return nil, errors.New(errContext, "Configuration file could not be loaded", err)
	}

	config = &Configuration{
		BuildersPath: loader.GetString(BuildersPathKey),
		Concurrency:  loader.GetInt(ConcurrencyKey),
		Credentials: &CredentialsConfiguration{
			StorageType:      loader.GetString(strings.Join([]string{CredentialsKey, CredentialsStorageTypeKey}, ".")),
			LocalStoragePath: loader.GetString(strings.Join([]string{CredentialsKey, CredentialsLocalStoragePathKey}, ".")),
			Format:           loader.GetString(strings.Join([]string{CredentialsKey, CredentialsFormatKey}, ".")),
		},
		DEPRECATEDBuilderPath:          loader.GetString(DEPRECATEDBuilderPathKey),
		DEPRECATEDBuildOnCascade:       loader.GetBool(DEPRECATEDBuildOnCascadeKey),
		DEPRECATEDNumWorkers:           loader.GetInt(DEPRECATEDNumWorkerKey),
		DEPRECATEDTreePathFile:         loader.GetString(DEPRECATEDTreePathFileKey),
		DEPRECATEDDockerCredentialsDir: loader.GetString(DEPRECATEDDockerCredentialsDirKey),
		EnableSemanticVersionTags:      loader.GetBool(EnableSemanticVersionTagsKey),
		ImagesPath:                     loader.GetString(ImagesPathKey),
		LogPathFile:                    loader.GetString(LogPathFileKey),
		LogWriter:                      logWriter,
		PushImages:                     loader.GetBool(PushImagesKey),
		SemanticVersionTagsTemplates:   loader.GetStringSlice(SemanticVersionTagsTemplatesKey),

		compatibility: compatibility,
		configFile:    file,
		fs:            fs,
		loader:        loader,
	}

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	if config.BuildersPath == "" {
		config.BuildersPath = DefaultBuildersPath
	}

	if config.Concurrency < 1 {
		config.Concurrency = concurrencyValue()
	}

	if config.Credentials.StorageType == "" {
		config.Credentials.StorageType = DefaultCredentialsStorage
	}

	// DEPRECATEDDockerCredentialsDir must be check to avoid deprecation warning
	if config.Credentials.StorageType == DefaultCredentialsStorage && config.Credentials.LocalStoragePath == "" && config.DEPRECATEDDockerCredentialsDir == "" {
		config.Credentials.LocalStoragePath = DefaultCredentialsLocalStoragePath
	}

	if config.Credentials.Format == "" {
		config.Credentials.Format = DefaultCredentialsFormat
	}

	if !config.EnableSemanticVersionTags {
		config.EnableSemanticVersionTags = DefaultEnableSemanticVersionTags
	}

	if config.ImagesPath == "" {
		config.ImagesPath = DefaultImagesPath
	}

	if config.LogPathFile == "" {
		config.LogPathFile = DefaultLogPathFile
	}

	logWriter, err = createLogWriter(fs, config.LogPathFile)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	config.LogWriter = logWriter

	if config.PushImages == false {
		config.PushImages = DefaultPushImages
	}

	if len(config.SemanticVersionTagsTemplates) == 0 {
		config.SemanticVersionTagsTemplates = append([]string{}, DefaultSemanticVersionTagsTemplates)
	}

	err = config.ValidateConfiguration()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return config, nil
}

// ConfigFileUsed return which is the config file used to load the configuration
func (c *Configuration) ConfigFileUsed() string {
	return c.configFile
}

// ReloadConfigurationFromFile
func (c *Configuration) ReloadConfigurationFromFile(file string) error {
	errContext := "(Configuration::ReloadConfigurationFromFile)"

	if file == "" {
		return errors.New(errContext, "Configuration file must be provided to reload configuration from file")
	}

	newConfig, err := LoadFromFile(c.fs, c.loader, file, c.compatibility)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = newConfig.ValidateConfiguration()
	if err != nil {
		return errors.New(errContext, "", err)
	}

	newConfig.configFile = file

	*c = *newConfig
	return nil
}

// ValidateConfiguration method validates the configuration
func (c *Configuration) ValidateConfiguration() error {

	errContext := "(Configuration::ValidateConfiguration)"

	// Note: It is not validated if fs is defined
	// if c.fs == nil {
	// 	return errors.New(errContext, "File system must be provided to create a new configuration")
	// }

	if c.BuildersPath == "" {
		return errors.New(errContext, "Invalid configuration, builders path must be provided")
	}

	if c.ImagesPath == "" {
		return errors.New(errContext, "Invalid configuration, images path must be provided")
	}

	if c.Concurrency < 1 {
		return errors.New(errContext, "Invalid configuration, concurrency must be greater than 0")
	}

	if c.Credentials != nil {
		if c.Credentials.StorageType == "" {
			return errors.New(errContext, "Invalid configuration, credentials storage type must be provided")
		}

		if (c.Credentials.Format != credentials.JSONFormat) && (c.Credentials.Format != credentials.YAMLFormat) {
			return errors.New(errContext, fmt.Sprintf("Invalid configuration, credentials format '%s' is not valid", c.Credentials.Format))
		}

		if c.Credentials.StorageType == credentials.LocalStore {
			if c.Credentials.LocalStoragePath == "" {
				return errors.New(errContext, "Invalid configuration, credentials local storage path must be provided")
			}
		}
	}

	return nil
}

// CheckCompatibility
func (c *Configuration) CheckCompatibility() error {

	errContext := "(Configuration::CheckCompatibility)"

	if c.compatibility == nil {
		return errors.New(errContext, "To ckeck configuration compatiblity is required a compatibilitier")
	}

	if c.DEPRECATEDTreePathFile != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDTreePathFileKey, ImagesPathKey))

		if c.ImagesPath != "" && c.ImagesPath != DefaultImagesPath {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDTreePathFileKey, ImagesPathKey, DEPRECATEDTreePathFileKey))
		}

		c.ImagesPath = c.DEPRECATEDTreePathFile
	}

	if c.DEPRECATEDBuilderPath != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDBuilderPathKey, BuildersPathKey))

		if c.BuildersPath != "" && c.BuildersPath != DefaultBuildersPath {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDBuilderPathKey, BuildersPathKey, DEPRECATEDBuilderPathKey))
		}

		c.BuildersPath = c.DEPRECATEDBuilderPath
	}

	if c.DEPRECATEDNumWorkers > 0 {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDNumWorkerKey, ConcurrencyKey))

		if c.Concurrency > 0 && c.Concurrency != concurrencyValue() {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDNumWorkerKey, ConcurrencyKey, DEPRECATEDNumWorkerKey))
		}

		c.Concurrency = c.DEPRECATEDNumWorkers
	}

	if c.DEPRECATEDBuildOnCascade == true {
		c.compatibility.AddChanged(fmt.Sprintf("'%s' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command", DEPRECATEDBuildOnCascadeKey))
		c.DEPRECATEDBuildOnCascade = false
	}

	if c.DEPRECATEDDockerCredentialsDir != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' block to configure credentials. Credentials local storage located in '%s' has precedence over '%s' block and is going to be used as default credentials store", DEPRECATEDDockerCredentialsDirKey, CredentialsKey, c.DEPRECATEDDockerCredentialsDir, CredentialsKey))

		if c.Credentials == nil {
			c.Credentials = &CredentialsConfiguration{}
		}

		c.Credentials.StorageType = DefaultCredentialsStorage
		c.Credentials.Format = DefaultCredentialsFormat
		c.Credentials.LocalStoragePath = c.DEPRECATEDDockerCredentialsDir

	}

	return nil
}

// createLogWriter return a io.Writer associated to the log file
func createLogWriter(fs afero.Fs, path string) (io.Writer, error) {

	var err error
	errContext := "(cli::stevedore)"
	writer := ioutil.Discard

	if path != "" {
		writer, err = fs.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, errors.New(errContext, "Log file can not be created", err)
		}
	}

	return writer, nil
}

// concurrencyValue returns the concurrency value from the configuration, in case of panic concurrency is set to 1
func concurrencyValue() (concurrency int) {

	defer func(v *int) {
		if err := recover(); err != nil {
			*v = 1
		}
	}(&concurrency)

	concurrency = int(math.Round(float64(runtime.NumCPU()) / 4))
	if concurrency < 1 {
		concurrency = 1
	}

	return
}
