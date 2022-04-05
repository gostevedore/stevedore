package configuration

import (
	"fmt"
	"html/template"
	"io"
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
	DEPRECATEDTreePathFile       string
	ImagesPath                   string
	DEPRECATEDBuilderPath        string
	BuildersPath                 string
	LogPathFile                  string
	DEPRECATEDNumWorkers         int
	Concurrency                  int
	PushImages                   bool
	DEPRECATEDBuildOnCascade     bool
	DockerCredentialsDir         string
	EnableSemanticVersionTags    bool
	SemanticVersionTagsTemplates []string

	compatibility Compatibilitier
	fs            afero.Fs
}

const (
	DefaultConfigFile   = "stevedore.yaml"
	DefaultConfigFolder = "."

	SecundaryConfigFolder = "."

	DEPRECATEDDefaultTreePathFile       = "stevedore.yaml"
	DefaultImagesPath                   = "stevedore.yaml"
	DEPRECATEDDefaultBuilderPath        = "stevedore.yaml"
	DefaultBuildersPath                 = "stevedore.yaml"
	DefaultLogPathFile                  = ""
	DEPRECATEDDefaultNumWorker          = 4
	DefaultPushImages                   = true
	DEPRECATEDDefaultBuildOnCascade     = false
	DefaultDockerCredentialsDir         = "credentials"
	DefaultEnableSemanticVersionTags    = false
	DefaultSemanticVersionTagsTemplates = "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"

	DEPRECATEDTreePathFileKey       = "tree_path"
	ImagesPathKey                   = "images_path"
	DEPRECATEDBuilderPathKey        = "builder_path"
	BuildersPathKey                 = "builders_path"
	LogPathFileKey                  = "log_path"
	DEPRECATEDNumWorkerKey          = "num_workers"
	ConcurrencyKey                  = "concurrency"
	PushImagesKey                   = "push_images"
	DEPRECATEDBuildOnCascadeKey     = "build_on_cascade"
	DockerCredentialsDirKey         = "docker_registry_credentials_dir"
	EnableSemanticVersionTagsKey    = "semantic_version_tags_enabled"
	SemanticVersionTagsTemplatesKey = "semantic_version_tags_templates"
)

// New method create a new configuration object
func New(fs afero.Fs, compatibility Compatibilitier) (*Configuration, error) {

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
	defaultConcurrency := runtime.NumCPU() / 4

	viper.SetDefault(DEPRECATEDTreePathFileKey, filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile))
	viper.SetDefault(ImagesPathKey, filepath.Join(DefaultConfigFolder, DefaultImagesPath))
	viper.SetDefault(DEPRECATEDBuilderPathKey, filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath))
	viper.SetDefault(BuildersPathKey, filepath.Join(DefaultConfigFolder, DefaultBuildersPath))
	viper.SetDefault(LogPathFileKey, DefaultLogPathFile)
	viper.SetDefault(DEPRECATEDNumWorkerKey, DEPRECATEDDefaultNumWorker)
	viper.SetDefault(ConcurrencyKey, defaultConcurrency)
	viper.SetDefault(PushImagesKey, DefaultPushImages)
	viper.SetDefault(DEPRECATEDBuildOnCascadeKey, DEPRECATEDDefaultBuildOnCascade)
	viper.SetDefault(DockerCredentialsDirKey, filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir))
	viper.SetDefault(EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags)
	viper.SetDefault(SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates})

	for _, alternativeConfigFolder := range alternativesConfigFolders {
		viper.AddConfigPath(alternativeConfigFolder)
	}
	// Set the default config folder as last default option
	viper.AddConfigPath(DefaultConfigFolder)

	// when configuration is created no error is shown if readinconfig files. It will use the defaults
	viper.ReadInConfig()

	config := &Configuration{
		ImagesPath:                   viper.GetString(ImagesPathKey),
		BuildersPath:                 viper.GetString(BuildersPathKey),
		LogPathFile:                  viper.GetString(LogPathFileKey),
		Concurrency:                  viper.GetInt(ConcurrencyKey),
		PushImages:                   viper.GetBool(PushImagesKey),
		DockerCredentialsDir:         viper.GetString(DockerCredentialsDirKey),
		EnableSemanticVersionTags:    viper.GetBool(EnableSemanticVersionTagsKey),
		SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),

		compatibility: compatibility,
		fs:            fs,
	}

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return config, nil
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

// LoadFromFile method returns a configuration object loaded from a file
func LoadFromFile(fs afero.Fs, file string, compatibility Compatibilitier) (*Configuration, error) {

	errContext := "(configuration::LoadFromFile)"

	viper.SetFs(fs)
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.New(errContext, "Configuration file could be loaded", err)
	}

	config := &Configuration{
		DEPRECATEDTreePathFile:       viper.GetString(DEPRECATEDTreePathFileKey),
		DEPRECATEDBuilderPath:        viper.GetString(DEPRECATEDBuilderPathKey),
		ImagesPath:                   viper.GetString(ImagesPathKey),
		BuildersPath:                 viper.GetString(BuildersPathKey),
		LogPathFile:                  viper.GetString(LogPathFileKey),
		DEPRECATEDNumWorkers:         viper.GetInt(DEPRECATEDNumWorkerKey),
		Concurrency:                  viper.GetInt(ConcurrencyKey),
		PushImages:                   viper.GetBool(PushImagesKey),
		DEPRECATEDBuildOnCascade:     viper.GetBool(DEPRECATEDBuildOnCascadeKey),
		DockerCredentialsDir:         viper.GetString(DockerCredentialsDirKey),
		EnableSemanticVersionTags:    viper.GetBool(EnableSemanticVersionTagsKey),
		SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),

		compatibility: compatibility,
	}

	err = config.CheckCompatibility()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return config, nil
}

// ReloadConfigurationFromFile
func (c *Configuration) ReloadConfigurationFromFile(fs afero.Fs, file string, compatibility Compatibilitier) error {
	errContext := "(Configuration::ReloadConfigurationFromFile)"
	newConfig, err := LoadFromFile(fs, file, compatibility)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	*c = *newConfig
	return nil
}

func (c *Configuration) String() string {
	str := ""

	str = fmt.Sprintln()

	str = fmt.Sprintln(str, ImagesPathKey, ": ", c.ImagesPath)
	str = fmt.Sprintln(str, BuildersPathKey, ": ", c.BuildersPath)
	str = fmt.Sprintln(str, LogPathFileKey, ": ", c.LogPathFile)
	str = fmt.Sprintln(str, ConcurrencyKey, ": ", c.Concurrency)
	str = fmt.Sprintln(str, PushImagesKey, ": ", c.PushImages)
	str = fmt.Sprintln(str, DockerCredentialsDirKey, ": ", c.DockerCredentialsDir)
	str = fmt.Sprintln(str, EnableSemanticVersionTagsKey, ":", c.EnableSemanticVersionTags)
	str = fmt.Sprintln(str, SemanticVersionTagsTemplatesKey, ":", c.SemanticVersionTagsTemplates)

	return str
}

func (c *Configuration) ToArray() ([][]string, error) {

	if c == nil {
		return nil, errors.New("(Configuration::ToArray)", "Configuration is nil")
	}

	arrayConfig := [][]string{}
	arrayConfig = append(arrayConfig, []string{ImagesPathKey, c.ImagesPath})
	arrayConfig = append(arrayConfig, []string{BuildersPathKey, c.BuildersPath})
	arrayConfig = append(arrayConfig, []string{LogPathFileKey, c.LogPathFile})
	arrayConfig = append(arrayConfig, []string{ConcurrencyKey, fmt.Sprint(c.Concurrency)})
	arrayConfig = append(arrayConfig, []string{PushImagesKey, fmt.Sprint(c.PushImages)})
	arrayConfig = append(arrayConfig, []string{DockerCredentialsDirKey, fmt.Sprint(c.DockerCredentialsDir)})
	arrayConfig = append(arrayConfig, []string{EnableSemanticVersionTagsKey, fmt.Sprint(c.EnableSemanticVersionTags)})
	semanticVersionTagsTemplatesValue := fmt.Sprint(c.SemanticVersionTagsTemplates)
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

		if c.ImagesPath == "" {
			c.ImagesPath = c.DEPRECATEDTreePathFile
		} else {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDTreePathFileKey, ImagesPathKey, ImagesPathKey))
		}
	}
	if c.DEPRECATEDBuilderPath != "" {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDBuilderPathKey, BuildersPathKey))

		if c.BuildersPath == "" {
			c.BuildersPath = c.DEPRECATEDBuilderPath
		} else {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDBuilderPathKey, BuildersPathKey, BuildersPathKey))
		}
	}
	if c.DEPRECATEDNumWorkers > 0 {
		c.compatibility.AddDeprecated(fmt.Sprintf("'%s' is deprecated and will be removed on v0.12.0, please use '%s' instead", DEPRECATEDNumWorkerKey, ConcurrencyKey))

		if c.Concurrency <= 0 {
			c.Concurrency = c.DEPRECATEDNumWorkers
		} else {
			c.compatibility.AddDeprecated(fmt.Sprintf("'%s' and '%s' are both defined, '%s' will be used", DEPRECATEDNumWorkerKey, ConcurrencyKey, ConcurrencyKey))
		}
	}
	if c.DEPRECATEDBuildOnCascade == true {
		c.compatibility.AddChanged(fmt.Sprintf("'%s' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command", DEPRECATEDBuildOnCascadeKey))
		c.DEPRECATEDBuildOnCascade = false
	}

	return nil
}
