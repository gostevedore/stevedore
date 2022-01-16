package configuration

import (
	"fmt"
	"html/template"
	"io"
	"os/user"
	"path/filepath"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/viper"
	"go.uber.org/zap/buffer"
)

type Configuration struct {
	TreePathFile                 string
	BuilderPathFile              string
	LogPathFile                  string
	NumWorkers                   int
	PushImages                   bool
	BuildOnCascade               bool
	DockerCredentialsDir         string
	EnableSemanticVersionTags    bool
	SemanticVersionTagsTemplates []string

	compatibility Compatibilitier
}

const (
	DefaultConfigFile   = "stevedore.yaml"
	DefaultConfigFolder = "."

	SecundaryConfigFolder = "."

	DefaultTreePathFile                 = "stevedore.yaml"
	DefaultBuilderPathFile              = "stevedore.yaml"
	DefaultLogPathFile                  = "/dev/null"
	DefaultNumWorker                    = 4
	DefaultPushImages                   = true
	DefaultBuildOnCascade               = false
	DefaultDockerCredentialsDir         = "credentials"
	DefaultEnableSemanticVersionTags    = false
	DefaultSemanticVersionTagsTemplates = "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"

	TreePathFileKey                 = "tree_path"
	BuilderPathFileKey              = "builder_path"
	LogPathFileKey                  = "log_path"
	NumWorkerKey                    = "num_workers"
	PushImagesKey                   = "push_images"
	BuildOnCascadeKey               = "build_on_cascade"
	DockerCredentialsDirKey         = "docker_registry_credentials_dir"
	EnableSemanticVersionTagsKey    = "semantic_version_tags_enabled"
	SemanticVersionTagsTemplatesKey = "semantic_version_tags_templates"
)

// New method create a new configuration object
func New() (*Configuration, error) {

	user, err := user.Current()
	if err != nil {
		return nil, errors.New("(configuration::New)", "Current user information can not be cached", err)
	}

	alternativesConfigFolders := []string{
		filepath.Join(user.HomeDir, ".config", "stevedore"),
		user.HomeDir,
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("stevedore")

	viper.SetConfigName(DefaultConfigFile)
	viper.SetConfigType("yaml")

	viper.SetDefault(TreePathFileKey, filepath.Join(DefaultConfigFolder, DefaultTreePathFile))
	viper.SetDefault(BuilderPathFileKey, filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile))
	viper.SetDefault(LogPathFileKey, DefaultLogPathFile)
	viper.SetDefault(NumWorkerKey, DefaultNumWorker)
	viper.SetDefault(PushImagesKey, DefaultPushImages)
	viper.SetDefault(BuildOnCascadeKey, DefaultBuildOnCascade)

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

	return &Configuration{
		TreePathFile:                 viper.GetString(TreePathFileKey),
		BuilderPathFile:              viper.GetString(BuilderPathFileKey),
		LogPathFile:                  viper.GetString(LogPathFileKey),
		NumWorkers:                   viper.GetInt(NumWorkerKey),
		PushImages:                   viper.GetBool(PushImagesKey),
		BuildOnCascade:               viper.GetBool(BuildOnCascadeKey),
		DockerCredentialsDir:         viper.GetString(DockerCredentialsDirKey),
		EnableSemanticVersionTags:    viper.GetBool(EnableSemanticVersionTagsKey),
		SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),
	}, nil
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

// LoadFromFile method returns a configuration object loaded from a file
func LoadFromFile(file string) (*Configuration, error) {

	//	viper.Reset()
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.New("(Configuration::LoadFromFile)", "Configuration could be load from '"+file+"'", err)
	}

	return &Configuration{
		TreePathFile:                 viper.GetString(TreePathFileKey),
		BuilderPathFile:              viper.GetString(BuilderPathFileKey),
		LogPathFile:                  viper.GetString(LogPathFileKey),
		NumWorkers:                   viper.GetInt(NumWorkerKey),
		PushImages:                   viper.GetBool(PushImagesKey),
		BuildOnCascade:               viper.GetBool(BuildOnCascadeKey),
		DockerCredentialsDir:         viper.GetString(DockerCredentialsDirKey),
		EnableSemanticVersionTags:    viper.GetBool(EnableSemanticVersionTagsKey),
		SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),
	}, nil
}

// ReloadConfigurationFromFile
func (c *Configuration) ReloadConfigurationFromFile(file string) error {
	newConfig, err := LoadFromFile(file)
	if err != nil {
		return errors.New("(Configuration::ReloadConfigurationFromFile)", "Configuration could not be reload from file '"+file+"'", err)
	}

	*c = *newConfig
	return nil
}

func (c *Configuration) String() string {
	str := ""

	str = fmt.Sprintln()

	str = fmt.Sprintln(str, TreePathFileKey, ": ", c.TreePathFile)
	str = fmt.Sprintln(str, BuilderPathFileKey, ": ", c.BuilderPathFile)
	str = fmt.Sprintln(str, LogPathFileKey, ": ", c.LogPathFile)
	str = fmt.Sprintln(str, NumWorkerKey, ": ", c.NumWorkers)
	str = fmt.Sprintln(str, PushImagesKey, ": ", c.PushImages)
	str = fmt.Sprintln(str, BuildOnCascadeKey, ": ", c.BuildOnCascade)
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
	arrayConfig = append(arrayConfig, []string{TreePathFileKey, c.TreePathFile})
	arrayConfig = append(arrayConfig, []string{BuilderPathFileKey, c.BuilderPathFile})
	arrayConfig = append(arrayConfig, []string{LogPathFileKey, c.LogPathFile})
	arrayConfig = append(arrayConfig, []string{NumWorkerKey, fmt.Sprint(c.NumWorkers)})
	arrayConfig = append(arrayConfig, []string{PushImagesKey, fmt.Sprint(c.PushImages)})
	arrayConfig = append(arrayConfig, []string{BuildOnCascadeKey, fmt.Sprint(c.BuildOnCascade)})
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

	// configFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	// if err != nil {
	// 	return errors.New("(configuration::CreateConfigurationFile)", fmt.Sprintf("File '%s' could not be opened", file), err)
	// }
	// defer configFile.Close()

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
