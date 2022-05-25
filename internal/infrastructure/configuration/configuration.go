package configuration

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap/buffer"
)

type Configuration struct {
	BuildersPath                 string
	Concurrency                  int
	DEPRECATEDBuilderPath        string
	DEPRECATEDBuildOnCascade     bool
	DEPRECATEDNumWorkers         int
	DEPRECATEDTreePathFile       string
	DockerCredentialsDir         string
	EnableSemanticVersionTags    bool
	ImagesPath                   string
	LogPathFile                  string
	LogWriter                    io.Writer
	PushImages                   bool
	SemanticVersionTagsTemplates []string

	compatibility Compatibilitier
	fs            afero.Fs
}

const (
	DefaultConfigFile   = "stevedore.yaml"
	DefaultConfigFolder = "."

	DefaultBuildersPath                 = "stevedore.yaml"
	DefaultDockerCredentialsDir         = "credentials"
	DefaultEnableSemanticVersionTags    = false
	DefaultImagesPath                   = "stevedore.yaml"
	DefaultLogPathFile                  = ""
	DefaultPushImages                   = true
	DefaultSemanticVersionTagsTemplates = "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
	DEPRECATEDDefaultBuilderPath        = "stevedore.yaml"
	DEPRECATEDDefaultBuildOnCascade     = false
	DEPRECATEDDefaultNumWorker          = 4
	DEPRECATEDDefaultTreePathFile       = "stevedore.yaml"

	BuildersPathKey                 = "builders_path"
	ConcurrencyKey                  = "concurrency"
	DEPRECATEDBuilderPathKey        = "builder_path"
	DEPRECATEDBuildOnCascadeKey     = "build_on_cascade"
	DEPRECATEDNumWorkerKey          = "num_workers"
	DEPRECATEDTreePathFileKey       = "tree_path"
	DockerCredentialsDirKey         = "docker_registry_credentials_dir"
	EnableSemanticVersionTagsKey    = "semantic_version_tags_enabled"
	ImagesPathKey                   = "images_path"
	LogPathFileKey                  = "log_path"
	PushImagesKey                   = "push_images"
	SemanticVersionTagsTemplatesKey = "semantic_version_tags_templates"
)

// New method create a new configuration object
func New(fs afero.Fs, compatibility Compatibilitier) (*Configuration, error) {

	var logWriter io.Writer

	errContext := "(Configuration::New)"

	user, err := user.Current()
	if err != nil {
		return nil, errors.New(errContext, "Current user information can not be cached", err)
	}

	alternativesConfigFolders := []string{
		filepath.Join(user.HomeDir, ".config", "stevedore"),
		user.HomeDir,
	}

	viper.SetFs(fs)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("stevedore")

	viper.SetConfigName(DefaultConfigFile)
	viper.SetConfigType("yaml")

	// dynamic default values
	defaultConcurrency := int(math.Round(float64(runtime.NumCPU()) / 4))

	viper.SetDefault(BuildersPathKey, filepath.Join(DefaultConfigFolder, DefaultBuildersPath))
	viper.SetDefault(ConcurrencyKey, defaultConcurrency)
	// viper.SetDefault(DEPRECATEDBuilderPathKey, filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath))
	// viper.SetDefault(DEPRECATEDBuildOnCascadeKey, DEPRECATEDDefaultBuildOnCascade)
	// viper.SetDefault(DEPRECATEDNumWorkerKey, DEPRECATEDDefaultNumWorker)
	// viper.SetDefault(DEPRECATEDTreePathFileKey, filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile))
	viper.SetDefault(DockerCredentialsDirKey, filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir))
	viper.SetDefault(EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags)
	viper.SetDefault(ImagesPathKey, filepath.Join(DefaultConfigFolder, DefaultImagesPath))
	viper.SetDefault(LogPathFileKey, DefaultLogPathFile)
	viper.SetDefault(PushImagesKey, DefaultPushImages)
	viper.SetDefault(SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates})

	for _, alternativeConfigFolder := range alternativesConfigFolders {
		viper.AddConfigPath(alternativeConfigFolder)
	}
	// Set the default config folder as last default option
	viper.AddConfigPath(DefaultConfigFolder)

	// when configuration is created no error is shown if readinconfig files. It will use the defaults
	viper.ReadInConfig()

	logWriter, err = createLogWriter(fs, viper.GetString(LogPathFileKey))
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	config := &Configuration{
		BuildersPath:              viper.GetString(BuildersPathKey),
		Concurrency:               viper.GetInt(ConcurrencyKey),
		DEPRECATEDBuilderPath:     viper.GetString(DEPRECATEDBuilderPathKey),
		DEPRECATEDBuildOnCascade:  viper.GetBool(DEPRECATEDBuildOnCascadeKey),
		DEPRECATEDNumWorkers:      viper.GetInt(DEPRECATEDNumWorkerKey),
		DEPRECATEDTreePathFile:    viper.GetString(DEPRECATEDTreePathFileKey),
		DockerCredentialsDir:      viper.GetString(DockerCredentialsDirKey),
		EnableSemanticVersionTags: viper.GetBool(EnableSemanticVersionTagsKey),
		ImagesPath:                viper.GetString(ImagesPathKey),
		// LogPathFile:                  viper.GetString(LogPathFileKey),
		LogWriter:                    logWriter,
		PushImages:                   viper.GetBool(PushImagesKey),
		SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),

		compatibility: compatibility,
		fs:            fs,
	}

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return config, nil
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

// LoadFromFile method returns a configuration object loaded from a file
func LoadFromFile(fs afero.Fs, file string, compatibility Compatibilitier) (*Configuration, error) {

	var err error
	var logWriter io.Writer

	errContext := "(configuration::LoadFromFile)"

	viper.SetFs(fs)
	viper.SetConfigFile(file)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.New(errContext, "Configuration file could be loaded", err)
	}

	logWriter, err = createLogWriter(fs, viper.GetString(LogPathFileKey))
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	config := &Configuration{
		BuildersPath:                 viper.GetString(BuildersPathKey),
		Concurrency:                  viper.GetInt(ConcurrencyKey),
		DEPRECATEDBuilderPath:        viper.GetString(DEPRECATEDBuilderPathKey),
		DEPRECATEDBuildOnCascade:     viper.GetBool(DEPRECATEDBuildOnCascadeKey),
		DEPRECATEDNumWorkers:         viper.GetInt(DEPRECATEDNumWorkerKey),
		DEPRECATEDTreePathFile:       viper.GetString(DEPRECATEDTreePathFileKey),
		DockerCredentialsDir:         viper.GetString(DockerCredentialsDirKey),
		EnableSemanticVersionTags:    viper.GetBool(EnableSemanticVersionTagsKey),
		ImagesPath:                   viper.GetString(ImagesPathKey),
		LogPathFile:                  viper.GetString(LogPathFileKey),
		LogWriter:                    logWriter,
		PushImages:                   viper.GetBool(PushImagesKey),
		SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),

		compatibility: compatibility,
		fs:            fs,
	}

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return config, nil
}

// ReloadConfigurationFromFile
func (c *Configuration) ReloadConfigurationFromFile(fs afero.Fs, file string, compatibility Compatibilitier) error {
	errContext := "(Configuration::ReloadConfigurationFromFile)"

	newConfig, err := LoadFromFile(fs, file, compatibility)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	*c = *newConfig
	return nil
}

func (c *Configuration) String() string {
	str := ""

	str = fmt.Sprintln()

	str = fmt.Sprintln(str, BuildersPathKey, ": ", c.BuildersPath)
	str = fmt.Sprintln(str, ConcurrencyKey, ": ", c.Concurrency)
	str = fmt.Sprintln(str, DockerCredentialsDirKey, ": ", c.DockerCredentialsDir)
	str = fmt.Sprintln(str, EnableSemanticVersionTagsKey, ":", c.EnableSemanticVersionTags)
	str = fmt.Sprintln(str, ImagesPathKey, ": ", c.ImagesPath)
	str = fmt.Sprintln(str, LogPathFileKey, ": ", c.LogPathFile)
	str = fmt.Sprintln(str, PushImagesKey, ": ", c.PushImages)
	str = fmt.Sprintln(str, SemanticVersionTagsTemplatesKey, ":", c.SemanticVersionTagsTemplates)

	return str
}

func (c *Configuration) ToArray() ([][]string, error) {

	if c == nil {
		return nil, errors.New("(Configuration::ToArray)", "Configuration is nil")
	}

	arrayConfig := [][]string{}
	semanticVersionTagsTemplatesValue := fmt.Sprint(c.SemanticVersionTagsTemplates)
	arrayConfig = append(arrayConfig, []string{BuildersPathKey, c.BuildersPath})
	arrayConfig = append(arrayConfig, []string{ConcurrencyKey, fmt.Sprint(c.Concurrency)})
	arrayConfig = append(arrayConfig, []string{DockerCredentialsDirKey, fmt.Sprint(c.DockerCredentialsDir)})
	arrayConfig = append(arrayConfig, []string{EnableSemanticVersionTagsKey, fmt.Sprint(c.EnableSemanticVersionTags)})
	arrayConfig = append(arrayConfig, []string{ImagesPathKey, c.ImagesPath})
	arrayConfig = append(arrayConfig, []string{LogPathFileKey, c.LogPathFile})
	arrayConfig = append(arrayConfig, []string{PushImagesKey, fmt.Sprint(c.PushImages)})
	arrayConfig = append(arrayConfig, []string{SemanticVersionTagsTemplatesKey, semanticVersionTagsTemplatesValue})

	return arrayConfig, nil
}

func ConfigurationHeaders() []string {
	h := []string{
		"PARAMETER",
		"VALUE",
	}

	return h
}

// CreateConfigurationFile
func (c *Configuration) WriteConfigurationFile(w io.Writer) error {

	var buff buffer.Buffer

	tmpl, err := template.New("configuration").Parse(configurationTemplate)
	if err != nil {
		return errors.New("(configuration::CreateConfigurationFile)", "Configuration template could not be parsed", err)
	}

	err = tmpl.Execute(&buff, c)
	if err != nil {
		return errors.New("(configuration::CreateConfigurationFile)", "Error applying variables to configuration template", err)

	}

	// golang does not support some charaters on raw strings and must be reprecented by another symbols
	// "`" is reprecented by "#u0060" and must be replaced to all its occurrences
	// Though there are some templating variables which must not be replaced by parser symbols "{" and "}" are also represented by "#u007b" and "#u007b"
	replacer := strings.NewReplacer("#u0060", "`", "#u007b", "{", "#u007d", "}")
	config := replacer.Replace(buff.String())

	fmt.Fprintln(w, config)

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

	return nil
}

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