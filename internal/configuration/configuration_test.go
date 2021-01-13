package configuration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/buffer"
)

const (
	testBaseDir = "test"
)

func TestViperBehavior(t *testing.T) {

	tests := []struct {
		desc   string
		config []byte
		res    *Configuration
	}{
		{
			desc:   "Testing empty configuration file",
			config: []byte(""),
			res: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   DefaultNumWorker,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               DefaultBuildOnCascade,
				DockerCredentialsDir:         DefaultDockerCredentialsDir,
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
		},
		{
			desc: "Testing overwrite log_path on configuration file",
			config: []byte(`
log_path: "/var/log/stevedore/stevedore.log"
num_workers: 5 
`),
			res: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  "/var/log/stevedore/stevedore.log",
				NumWorkers:                   5,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               DefaultBuildOnCascade,
				DockerCredentialsDir:         DefaultDockerCredentialsDir,
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		viper.Reset()
		viper.SetConfigType("yaml")
		viper.SetDefault(TreePathFileKey, filepath.Join(DefaultConfigFolder, DefaultTreePathFile))
		viper.SetDefault(BuilderPathFileKey, filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile))
		viper.SetDefault(LogPathFileKey, DefaultLogPathFile)
		viper.SetDefault(NumWorkerKey, DefaultNumWorker)
		viper.SetDefault(PushImagesKey, DefaultPushImages)
		viper.SetDefault(BuildOnCascadeKey, DefaultBuildOnCascade)
		viper.SetDefault(DockerCredentialsDirKey, DefaultDockerCredentialsDir)
		viper.SetDefault(EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags)
		viper.SetDefault(SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates})

		viper.ReadConfig(bytes.NewBuffer(test.config))

		c := &Configuration{
			TreePathFile:                 viper.GetString(TreePathFileKey),
			BuilderPathFile:              viper.GetString(BuilderPathFileKey),
			LogPathFile:                  viper.GetString(LogPathFileKey),
			NumWorkers:                   viper.GetInt(NumWorkerKey),
			PushImages:                   viper.GetBool(PushImagesKey),
			BuildOnCascade:               viper.GetBool(BuildOnCascadeKey),
			DockerCredentialsDir:         viper.GetString(DockerCredentialsDirKey),
			EnableSemanticVersionTags:    viper.GetBool(EnableSemanticVersionTagsKey),
			SemanticVersionTagsTemplates: viper.GetStringSlice(SemanticVersionTagsTemplatesKey),
		}

		assert.Equal(t, test.res, c, "Unpexpected configuration value")

	}
}

func TestNew(t *testing.T) {
	viper.Reset()

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tests := []struct {
		desc     string
		preFunc  func()
		postFunc func()
		res      *Configuration
	}{
		{
			desc: "Testing all defaults",
			preFunc: func() {
				viper.Reset()
			},
			postFunc: nil,
			res: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   DefaultNumWorker,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               DefaultBuildOnCascade,
				DockerCredentialsDir:         filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
		},
		{
			desc: "Testing set num_workers using environment variables",
			preFunc: func() {
				os.Setenv("STEVEDORE_NUM_WORKERS", "5")
				os.Setenv("STEVEDORE_BUILD_ON_CASCADE", "true")
				viper.Reset()
			},
			postFunc: func() {
				os.Unsetenv("STEVEDORE_NUM_WORKERS")
				os.Unsetenv("STEVEDORE_BUILD_ON_CASCADE")
			},
			res: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   5,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               true,
				DockerCredentialsDir:         filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)
		if test.preFunc != nil {
			test.preFunc()
		}

		c, err := New()
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, test.res, c, "Configuration values does not coincide")

		if test.postFunc != nil {
			test.postFunc()
		}
	}
}

func TestLoadFromFile(t *testing.T) {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tests := []struct {
		desc string
		file string
		err  error
		res  *Configuration
	}{
		{
			desc: "Testing error when loading configuration from file",
			file: "unknown",
			err:  errors.New("(Configuration::LoadFromFile)", "Configuration could be load from 'unknown'", errors.New("", "open unknown: no such file or directory")),
			res:  nil,
		},
		{
			desc: "Testing loading configuration from file",
			file: filepath.Join(testBaseDir, "stevedore.yaml"),
			err:  nil,
			res: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   10,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               true,
				DockerCredentialsDir:         filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		c, err := LoadFromFile(test.file)
		if err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res, c, "Configuration values does not coincide")
		}

	}

}

func TestReloadConfigurationFromFile(t *testing.T) {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tests := []struct {
		desc string
		file string
		err  error
		res  *Configuration
	}{
		{
			desc: "Testing error when reload configuration from file",
			file: "unknown",
			err: errors.New("(Configuration::ReloadConfigurationFromFile)", "Configuration could not be reload from file 'unknown'",
				errors.New("(Configuration::LoadFromFile)", "Configuration could be load from 'unknown'",
					errors.New("", "open unknown: no such file or directory"))),
			res: nil,
		},
		{
			desc: "Testing reload configuration from file",
			file: filepath.Join(testBaseDir, "stevedore.yaml"),
			err:  nil,
			res: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   10,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               true,
				DockerCredentialsDir:         filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		config, err := New()
		if err != nil {
			t.Error(err.Error())
		}
		err = config.ReloadConfigurationFromFile(test.file)

		if err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res, config, "Configuration values does not coincide")
		}
	}
}

func TestToArray(t *testing.T) {

	tests := []struct {
		desc   string
		config *Configuration
		res    [][]string
		err    error
	}{
		{
			desc:   "Testing transform a nil configuration to array",
			config: nil,
			res:    nil,
			err:    errors.New("(Configuration::ToArray)", "Configuration is nil"),
		},
		{
			desc: "Testing transform configuration to array",
			config: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   DefaultNumWorker,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               DefaultBuildOnCascade,
				DockerCredentialsDir:         filepath.Join("$HOME", ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
			res: [][]string{
				{TreePathFileKey, filepath.Join(DefaultConfigFolder, DefaultTreePathFile)},
				{BuilderPathFileKey, filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile)},
				{LogPathFileKey, DefaultLogPathFile},
				{NumWorkerKey, fmt.Sprint(DefaultNumWorker)},
				{PushImagesKey, fmt.Sprint(DefaultPushImages)},
				{BuildOnCascadeKey, fmt.Sprint(DefaultBuildOnCascade)},
				{DockerCredentialsDirKey, filepath.Join("$HOME", ".config", "stevedore", DefaultDockerCredentialsDir)},
				{EnableSemanticVersionTagsKey, fmt.Sprint(DefaultEnableSemanticVersionTags)},
				{SemanticVersionTagsTemplatesKey, fmt.Sprint(fmt.Sprintf("[%s]", DefaultSemanticVersionTagsTemplates))},
			},
			err: nil,
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		array, err := test.config.ToArray()

		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res, array, "Configuration values does not coincide")
		}
	}
}

func TestString(t *testing.T) {

	tests := []struct {
		desc   string
		config *Configuration
		res    string
	}{
		{
			desc: "Testing transform configuration to string",
			config: &Configuration{
				TreePathFile:                 filepath.Join(DefaultConfigFolder, DefaultTreePathFile),
				BuilderPathFile:              filepath.Join(DefaultConfigFolder, DefaultBuilderPathFile),
				LogPathFile:                  DefaultLogPathFile,
				NumWorkers:                   DefaultNumWorker,
				PushImages:                   DefaultPushImages,
				BuildOnCascade:               DefaultBuildOnCascade,
				DockerCredentialsDir:         filepath.Join(DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
			res: `
 tree_path :  stevedore.yaml
 builder_path :  stevedore.yaml
 log_path :  /dev/null
 num_workers :  4
 push_images :  true
 build_on_cascade :  false
 docker_registry_credentials_dir :  credentials
 semantic_version_tags_enabled : false
 semantic_version_tags_templates : [{{ .Major }}.{{ .Minor }}.{{ .Patch }}]
`,
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		str := test.config.String()

		assert.Equal(t, test.res, str, "Configuration values does not coincide")
	}
}

func TestConfigurationHeaders(t *testing.T) {

	t.Log("Testing list configuration header")
	expected := []string{"PARAMETER", "VALUE"}
	res := ConfigurationHeaders()

	assert.Equal(t, expected, res)
}

func TestCreateConfigurationFile(t *testing.T) {

	var err error
	var buff buffer.Buffer
	var expected []byte

	config := &Configuration{
		TreePathFile:              "test_stevedore.yaml",
		BuilderPathFile:           "test_stevedore.yaml",
		LogPathFile:               "test_stevedore.log",
		NumWorkers:                8,
		PushImages:                true,
		BuildOnCascade:            true,
		DockerCredentialsDir:      ".credentials",
		EnableSemanticVersionTags: true,
		SemanticVersionTagsTemplates: []string{
			"{{ .Major }}",
			"{{ .Major }}.{{ .Minor }}",
		},
	}

	expected, err = ioutil.ReadFile("test/stevedore.yaml.golden")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = config.WriteConfigurationFile(&buff)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, expected, buff.Bytes(), "Unexpected response")
}
