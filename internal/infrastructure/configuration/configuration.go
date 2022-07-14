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
	StorageType      string
	LocalStoragePath string
	Format           string
}

type Configuration struct {
	BuildersPath                   string
	Concurrency                    int
	DEPRECATEDBuilderPath          string
	DEPRECATEDBuildOnCascade       bool
	DEPRECATEDNumWorkers           int
	DEPRECATEDTreePathFile         string
	DEPRECATEDDockerCredentialsDir string

	Credentials *CredentialsConfiguration

	EnableSemanticVersionTags    bool
	ImagesPath                   string
	LogPathFile                  string
	LogWriter                    io.Writer
	PushImages                   bool
	SemanticVersionTagsTemplates []string

	compatibility Compatibilitier
	fs            afero.Fs
	loader        ConfigurationLoader
}

const (
	DefaultConfigFile          = "stevedore"
	DefaultConfigFileExtention = "yaml"
	DefaultConfigFolder        = "."

	DefaultBuildersPath                   = "stevedore.yaml"
	DefaultCredentialsFormat              = credentials.JSONFormat
	DefaultCredentialsLocalStoragePath    = "credentials"
	DefaultCredentialsStorage             = credentials.LocalStore
	DefaultEnableSemanticVersionTags      = false
	DefaultImagesPath                     = "stevedore.yaml"
	DefaultLogPathFile                    = ""
	DefaultPushImages                     = true
	DefaultSemanticVersionTagsTemplates   = "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
	DEPRECATEDDefaultBuilderPath          = "stevedore.yaml"
	DEPRECATEDDefaultBuildOnCascade       = false
	DEPRECATEDDefaultDockerCredentialsDir = "credentials"
	DEPRECATEDDefaultNumWorker            = 4
	DEPRECATEDDefaultTreePathFile         = "stevedore.yaml"

	BuildersPathKey                   = "builders_path"
	ConcurrencyKey                    = "concurrency"
	CredentialsFormatKey              = "format"
	CredentialsKey                    = "credentials"
	CredentialsLocalStoragePathKey    = "local_storage_path"
	CredentialsStorageTypeKey         = "storage_type"
	DEPRECATEDBuilderPathKey          = "builder_path"
	DEPRECATEDBuildOnCascadeKey       = "build_on_cascade"
	DEPRECATEDDockerCredentialsDirKey = "docker_registry_credentials_dir"
	DEPRECATEDNumWorkerKey            = "num_workers"
	DEPRECATEDTreePathFileKey         = "tree_path"
	EnableSemanticVersionTagsKey      = "semantic_version_tags_enabled"
	ImagesPathKey                     = "images_path"
	LogPathFileKey                    = "log_path"
	PushImagesKey                     = "push_images"
	SemanticVersionTagsTemplatesKey   = "semantic_version_tags_templates"
)

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
	//loader.SetDefault(DEPRECATEDDockerCredentialsDirKey, filepath.Join(user.HomeDir, ".config", "stevedore", DEPRECATEDDefaultDockerCredentialsDir))
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

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return config, nil
}

// LoadFromFile method returns a configuration object loaded from a file
func LoadFromFile(fs afero.Fs, loader ConfigurationLoader, file string, compatibility Compatibilitier) (*Configuration, error) {

	var err error
	var logWriter io.Writer

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
		return nil, errors.New(errContext, "Configuration file could be loaded", err)
	}

	logWriter, err = createLogWriter(fs, loader.GetString(LogPathFileKey))
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	config := &Configuration{
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
		fs:            fs,
		loader:        loader,
	}

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return config, nil
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

	*c = *newConfig
	return nil
}

func (c *Configuration) isValid() bool {
	return true
}

// func (c *Configuration) String() string {
// 	str := ""

// 	str = fmt.Sprintln()

// 	str = fmt.Sprintln(str, BuildersPathKey, ": ", c.BuildersPath)
// 	str = fmt.Sprintln(str, ConcurrencyKey, ": ", c.Concurrency)
// 	str = fmt.Sprintln(str, DEPRECATEDDockerCredentialsDirKey, ": ", c.DEPRECATEDDockerCredentialsDir)
// 	str = fmt.Sprintln(str, EnableSemanticVersionTagsKey, ":", c.EnableSemanticVersionTags)
// 	str = fmt.Sprintln(str, ImagesPathKey, ": ", c.ImagesPath)
// 	str = fmt.Sprintln(str, LogPathFileKey, ": ", c.LogPathFile)
// 	str = fmt.Sprintln(str, PushImagesKey, ": ", c.PushImages)
// 	str = fmt.Sprintln(str, SemanticVersionTagsTemplatesKey, ":", c.SemanticVersionTagsTemplates)

// 	return str
// }

// func (c *Configuration) ToArray() ([][]string, error) {

// 	if c == nil {
// 		return nil, errors.New("(Configuration::ToArray)", "Configuration is nil")
// 	}

// 	arrayConfig := [][]string{}
// 	semanticVersionTagsTemplatesValue := fmt.Sprint(c.SemanticVersionTagsTemplates)
// 	arrayConfig = append(arrayConfig, []string{BuildersPathKey, c.BuildersPath})
// 	arrayConfig = append(arrayConfig, []string{ConcurrencyKey, fmt.Sprint(c.Concurrency)})
// 	arrayConfig = append(arrayConfig, []string{DEPRECATEDDockerCredentialsDirKey, fmt.Sprint(c.DEPRECATEDDockerCredentialsDir)})
// 	arrayConfig = append(arrayConfig, []string{EnableSemanticVersionTagsKey, fmt.Sprint(c.EnableSemanticVersionTags)})
// 	arrayConfig = append(arrayConfig, []string{ImagesPathKey, c.ImagesPath})
// 	arrayConfig = append(arrayConfig, []string{LogPathFileKey, c.LogPathFile})
// 	arrayConfig = append(arrayConfig, []string{PushImagesKey, fmt.Sprint(c.PushImages)})
// 	arrayConfig = append(arrayConfig, []string{SemanticVersionTagsTemplatesKey, semanticVersionTagsTemplatesValue})

// 	return arrayConfig, nil
// }

// func ConfigurationHeaders() []string {
// 	h := []string{
// 		"PARAMETER",
// 		"VALUE",
// 	}

// 	return h
// }

// CheckCompatibility
func (c *Configuration) CheckCompatibility() error {

	errContext := "(Configuration::CheckCompatibility)"

	if c.compatibility == nil {
		return errors.New(errContext, "To ckeck configuration compatiblity is required a compatibilitier")
	}

	if c.DEPRECATEDTreePathFile != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDTreePathFileKey, ImagesPathKey))

		c.ImagesPath = c.DEPRECATEDTreePathFile
		if c.ImagesPath == "" {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDTreePathFileKey, ImagesPathKey, DEPRECATEDTreePathFileKey))
		}
	}

	if c.DEPRECATEDBuilderPath != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDBuilderPathKey, BuildersPathKey))

		c.BuildersPath = c.DEPRECATEDBuilderPath

		if c.BuildersPath == "" {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDBuilderPathKey, BuildersPathKey, DEPRECATEDBuilderPathKey))
		}
	}

	if c.DEPRECATEDNumWorkers > 0 {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDNumWorkerKey, ConcurrencyKey))

		c.Concurrency = c.DEPRECATEDNumWorkers

		if c.Concurrency <= 0 {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDNumWorkerKey, ConcurrencyKey, DEPRECATEDNumWorkerKey))
		}
	}

	if c.DEPRECATEDBuildOnCascade == true {
		c.compatibility.AddChanged(fmt.Sprintf("'%s' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command", DEPRECATEDBuildOnCascadeKey))
		c.DEPRECATEDBuildOnCascade = false
	}

	if c.DEPRECATEDDockerCredentialsDir != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' block to configure credentials. Credentials local storage located in '%s' is going to be used as default", DEPRECATEDDockerCredentialsDirKey, CredentialsKey, DefaultCredentialsLocalStoragePath))

		if c.Credentials == nil {
			c.Credentials = &CredentialsConfiguration{}
		}

		if c.Credentials.StorageType == "" {
			c.Credentials.StorageType = DefaultCredentialsStorage
			c.Credentials.LocalStoragePath = c.DEPRECATEDDockerCredentialsDir
		} else {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and 'credentials' block are both defined, local credentials storage on '%s' will be used", DEPRECATEDDockerCredentialsDirKey, c.DEPRECATEDDockerCredentialsDir))
		}

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
