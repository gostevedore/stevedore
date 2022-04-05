package configuration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/compatibility"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/buffer"
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
				DEPRECATEDTreePathFile:       filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile),
				DEPRECATEDBuilderPath:        filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath),
				LogPathFile:                  DefaultLogPathFile,
				DEPRECATEDNumWorkers:         DEPRECATEDDefaultNumWorker,
				Concurrency:                  4,
				PushImages:                   DefaultPushImages,
				DEPRECATEDBuildOnCascade:     DEPRECATEDDefaultBuildOnCascade,
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
				DEPRECATEDTreePathFile:       filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile),
				DEPRECATEDBuilderPath:        filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath),
				LogPathFile:                  "/var/log/stevedore/stevedore.log",
				DEPRECATEDNumWorkers:         5,
				Concurrency:                  4,
				PushImages:                   DefaultPushImages,
				DEPRECATEDBuildOnCascade:     DEPRECATEDDefaultBuildOnCascade,
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
		viper.SetDefault(DEPRECATEDTreePathFileKey, filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile))
		viper.SetDefault(DEPRECATEDBuilderPathKey, filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath))
		viper.SetDefault(LogPathFileKey, DefaultLogPathFile)
		viper.SetDefault(DEPRECATEDNumWorkerKey, DEPRECATEDDefaultNumWorker)
		viper.SetDefault(PushImagesKey, DefaultPushImages)
		viper.SetDefault(DEPRECATEDBuildOnCascadeKey, DEPRECATEDDefaultBuildOnCascade)
		viper.SetDefault(DockerCredentialsDirKey, DefaultDockerCredentialsDir)
		viper.SetDefault(EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags)
		viper.SetDefault(SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates})

		viper.ReadConfig(bytes.NewBuffer(test.config))

		c := &Configuration{
			DEPRECATEDTreePathFile:       viper.GetString(DEPRECATEDTreePathFileKey),
			DEPRECATEDBuilderPath:        viper.GetString(DEPRECATEDBuilderPathKey),
			LogPathFile:                  viper.GetString(LogPathFileKey),
			DEPRECATEDNumWorkers:         viper.GetInt(DEPRECATEDNumWorkerKey),
			Concurrency:                  4,
			PushImages:                   viper.GetBool(PushImagesKey),
			DEPRECATEDBuildOnCascade:     viper.GetBool(DEPRECATEDBuildOnCascadeKey),
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
		desc              string
		preFunc           func()
		postFunc          func()
		res               *Configuration
		compatibility     Compatibilitier
		prepareAssertFunc func(c Compatibilitier)
	}{
		{
			desc: "Testing all defaults",
			preFunc: func() {
				viper.Reset()
			},
			postFunc: nil,
			res: &Configuration{
				ImagesPath:                   filepath.Join(DefaultConfigFolder, DefaultImagesPath),
				BuildersPath:                 filepath.Join(DefaultConfigFolder, DefaultBuildersPath),
				LogPathFile:                  DefaultLogPathFile,
				Concurrency:                  4,
				PushImages:                   DefaultPushImages,
				DockerCredentialsDir:         filepath.Join(user.HomeDir, ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
			compatibility:     &compatibility.MockCompatibility{},
			prepareAssertFunc: func(c Compatibilitier) {},
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
				ImagesPath:                   filepath.Join(DefaultConfigFolder, DefaultImagesPath),
				BuildersPath:                 filepath.Join(DefaultConfigFolder, DefaultBuildersPath),
				LogPathFile:                  DefaultLogPathFile,
				Concurrency:                  4,
				PushImages:                   DefaultPushImages,
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

		if test.prepareAssertFunc != nil {
			test.prepareAssertFunc(test.compatibility)
		}

		c, err := New(afero.NewMemMapFs(), test.compatibility)
		if err != nil {
			t.Error(err.Error())
		} else {
			assert.Equal(t, test.res.DEPRECATEDTreePathFile, c.DEPRECATEDTreePathFile)
			assert.Equal(t, test.res.DEPRECATEDBuilderPath, c.DEPRECATEDBuilderPath)
			assert.Equal(t, test.res.ImagesPath, c.ImagesPath)
			assert.Equal(t, test.res.BuildersPath, c.BuildersPath)
			assert.Equal(t, test.res.LogPathFile, c.LogPathFile)
			assert.Equal(t, test.res.DEPRECATEDNumWorkers, c.DEPRECATEDNumWorkers)
			assert.Equal(t, test.res.Concurrency, c.Concurrency)
			assert.Equal(t, test.res.PushImages, c.PushImages)
			assert.Equal(t, test.res.DEPRECATEDBuildOnCascade, c.DEPRECATEDBuildOnCascade)
			assert.Equal(t, test.res.DockerCredentialsDir, c.DockerCredentialsDir)
			assert.Equal(t, test.res.EnableSemanticVersionTags, c.EnableSemanticVersionTags)
			assert.Equal(t, test.res.SemanticVersionTagsTemplates, c.SemanticVersionTagsTemplates)
		}
		if test.postFunc != nil {
			test.postFunc()
		}
	}
}

func TestLoadFromFile(t *testing.T) {

	var err error
	errContext := "(Configuration::LoadFromFile)"

	baseDir := "/config"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)
	err = afero.WriteFile(testFs, filepath.Join(baseDir, "stevedore.yaml"), []byte(`
images_path: mystevedore.yaml
builders_path: mystevedore.yaml
log_path: mystevedore.log
concurrency: 10
push_images: false
build_on_cascade: true
docker_registry_credentials_dir: mycredentials
semantic_version_tags_enabled: true
semantic_version_tags_templates:
  - "{{ -Major }}"
  - "{{ -Major }}.{{ .Minor }}"
`), 0644)

	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		file              string
		err               error
		res               *Configuration
		compatibility     Compatibilitier
		prepareAssertFunc func(c Compatibilitier)
	}{
		{
			desc: "Testing error when loading configuration from file",
			file: "unknown",
			err:  errors.New(errContext, "Configuration file could be loaded", errors.New("", "open unknown: file does not exist")),
			res:  nil,
		},
		{
			desc: "Testing loading configuration from file",
			file: filepath.Join(baseDir, "stevedore.yaml"),
			err:  &errors.Error{},
			res: &Configuration{
				ImagesPath:                "mystevedore.yaml",
				DEPRECATEDTreePathFile:    "stevedore.yaml",
				BuildersPath:              "mystevedore.yaml",
				DEPRECATEDBuilderPath:     "stevedore.yaml",
				LogPathFile:               "mystevedore.log",
				DEPRECATEDNumWorkers:      4,
				Concurrency:               10,
				PushImages:                false,
				DEPRECATEDBuildOnCascade:  false,
				DockerCredentialsDir:      "mycredentials",
				EnableSemanticVersionTags: true,
				SemanticVersionTagsTemplates: []string{
					"{{ -Major }}",
					"{{ -Major }}.{{ .Minor }}",
				},
			},
			compatibility: compatibility.NewMockCompatibility(),
			prepareAssertFunc: func(c Compatibilitier) {
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'tree_path' is deprecated and will be removed on v0.12.0, please use 'images_path' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'tree_path' and 'images_path' are both defined, 'images_path' will be used"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'builder_path' is deprecated and will be removed on v0.12.0, please use 'builders_path' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'builder_path' and 'builders_path' are both defined, 'builders_path' will be used"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'num_workers' is deprecated and will be removed on v0.12.0, please use 'concurrency' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'num_workers' and 'concurrency' are both defined, 'concurrency' will be used"})
				c.(*compatibility.MockCompatibility).On("AddChanged", []string{"'build_on_cascade' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command"})
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		if test.prepareAssertFunc != nil {
			test.prepareAssertFunc(test.compatibility)
		}

		config, err := LoadFromFile(testFs, test.file, test.compatibility)
		if err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res.BuildersPath, config.BuildersPath)
			assert.Equal(t, test.res.Concurrency, config.Concurrency)
			assert.Equal(t, test.res.DEPRECATEDBuilderPath, config.DEPRECATEDBuilderPath)
			assert.Equal(t, test.res.DEPRECATEDBuildOnCascade, config.DEPRECATEDBuildOnCascade)
			assert.Equal(t, test.res.DEPRECATEDNumWorkers, config.DEPRECATEDNumWorkers)
			assert.Equal(t, test.res.DEPRECATEDTreePathFile, config.DEPRECATEDTreePathFile)
			assert.Equal(t, test.res.DockerCredentialsDir, config.DockerCredentialsDir)
			assert.Equal(t, test.res.EnableSemanticVersionTags, config.EnableSemanticVersionTags)
			assert.Equal(t, test.res.ImagesPath, config.ImagesPath)
			assert.Equal(t, test.res.LogPathFile, config.LogPathFile)
			assert.Equal(t, test.res.PushImages, config.PushImages)
			assert.Equal(t, test.res.SemanticVersionTagsTemplates, config.SemanticVersionTagsTemplates)

		}

	}

}

func TestReloadConfigurationFromFile(t *testing.T) {
	errContext := "(Configuration::ReloadConfigurationFromFile)"
	var err error

	baseDir := "/config"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)
	err = afero.WriteFile(testFs, filepath.Join(baseDir, "stevedore.yaml"), []byte(`
images_path: mystevedore.yaml
builders_path: mystevedore.yaml
log_path: mystevedore.log
concurrency: 10
push_images: false
build_on_cascade: true
docker_registry_credentials_dir: mycredentials
semantic_version_tags_enabled: true
semantic_version_tags_templates:
  - "{{ -Major }}"
  - "{{ -Major }}.{{ .Minor }}"
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		file              string
		err               error
		res               *Configuration
		compatibility     Compatibilitier
		prepareAssertFunc func(c Compatibilitier)
	}{
		{
			desc: "Testing error when reload configuration from file",
			file: "unknown",
			err: errors.New(errContext, "Configuration file could be loaded",
				errors.New("", "open unknown: file does not exist")),
			res: nil,
		},
		{
			desc: "Testing reload configuration from file",
			file: filepath.Join(baseDir, "stevedore.yaml"),
			err:  &errors.Error{},
			res: &Configuration{
				DEPRECATEDTreePathFile:    "stevedore.yaml",
				ImagesPath:                "mystevedore.yaml",
				DEPRECATEDBuilderPath:     "stevedore.yaml",
				BuildersPath:              "mystevedore.yaml",
				LogPathFile:               "mystevedore.log",
				DEPRECATEDNumWorkers:      4,
				Concurrency:               10,
				PushImages:                false,
				DEPRECATEDBuildOnCascade:  false,
				DockerCredentialsDir:      "mycredentials",
				EnableSemanticVersionTags: true,
				SemanticVersionTagsTemplates: []string{
					"{{ -Major }}",
					"{{ -Major }}.{{ .Minor }}",
				},
			},
			compatibility: compatibility.NewMockCompatibility(),
			prepareAssertFunc: func(c Compatibilitier) {
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'tree_path' is deprecated and will be removed on v0.12.0, please use 'images_path' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'tree_path' and 'images_path' are both defined, 'images_path' will be used"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'builder_path' is deprecated and will be removed on v0.12.0, please use 'builders_path' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'builder_path' and 'builders_path' are both defined, 'builders_path' will be used"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'num_workers' is deprecated and will be removed on v0.12.0, please use 'concurrency' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'num_workers' and 'concurrency' are both defined, 'concurrency' will be used"})
				c.(*compatibility.MockCompatibility).On("AddChanged", []string{"'build_on_cascade' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command"})
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		if test.prepareAssertFunc != nil {
			test.prepareAssertFunc(test.compatibility)
		}

		config, err := New(testFs, test.compatibility)
		if err != nil {
			t.Error(err.Error())
		}
		err = config.ReloadConfigurationFromFile(testFs, test.file, test.compatibility)

		if err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res.BuildersPath, config.BuildersPath)
			assert.Equal(t, test.res.Concurrency, config.Concurrency)
			assert.Equal(t, test.res.DEPRECATEDBuilderPath, config.DEPRECATEDBuilderPath)
			assert.Equal(t, test.res.DEPRECATEDBuildOnCascade, config.DEPRECATEDBuildOnCascade)
			assert.Equal(t, test.res.DEPRECATEDNumWorkers, config.DEPRECATEDNumWorkers)
			assert.Equal(t, test.res.DEPRECATEDTreePathFile, config.DEPRECATEDTreePathFile)
			assert.Equal(t, test.res.DockerCredentialsDir, config.DockerCredentialsDir)
			assert.Equal(t, test.res.EnableSemanticVersionTags, config.EnableSemanticVersionTags)
			assert.Equal(t, test.res.ImagesPath, config.ImagesPath)
			assert.Equal(t, test.res.LogPathFile, config.LogPathFile)
			assert.Equal(t, test.res.PushImages, config.PushImages)
			assert.Equal(t, test.res.SemanticVersionTagsTemplates, config.SemanticVersionTagsTemplates)
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
				ImagesPath:                   filepath.Join(DefaultConfigFolder, DefaultImagesPath),
				BuildersPath:                 filepath.Join(DefaultConfigFolder, DefaultBuildersPath),
				LogPathFile:                  DefaultLogPathFile,
				Concurrency:                  4,
				PushImages:                   DefaultPushImages,
				DockerCredentialsDir:         filepath.Join("$HOME", ".config", "stevedore", DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
			res: [][]string{
				{BuildersPathKey, filepath.Join(DefaultConfigFolder, DefaultBuildersPath)},
				{ConcurrencyKey, fmt.Sprintf("%d", runtime.NumCPU()/4)},
				{DockerCredentialsDirKey, filepath.Join("$HOME", ".config", "stevedore", DefaultDockerCredentialsDir)},
				{EnableSemanticVersionTagsKey, fmt.Sprint(DefaultEnableSemanticVersionTags)},
				{ImagesPathKey, filepath.Join(DefaultConfigFolder, DefaultImagesPath)},
				{LogPathFileKey, DefaultLogPathFile},
				{PushImagesKey, fmt.Sprint(DefaultPushImages)},
				{SemanticVersionTagsTemplatesKey, fmt.Sprintf("[%s]", DefaultSemanticVersionTagsTemplates)},
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
				ImagesPath:                   filepath.Join(DefaultConfigFolder, DefaultImagesPath),
				BuildersPath:                 filepath.Join(DefaultConfigFolder, DefaultBuildersPath),
				LogPathFile:                  DefaultLogPathFile,
				Concurrency:                  4,
				PushImages:                   DefaultPushImages,
				DEPRECATEDBuildOnCascade:     DEPRECATEDDefaultBuildOnCascade,
				DockerCredentialsDir:         filepath.Join(DefaultDockerCredentialsDir),
				EnableSemanticVersionTags:    DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates: []string{DefaultSemanticVersionTagsTemplates},
			},
			res: `
 builders_path :  stevedore.yaml
 concurrency :  4
 docker_registry_credentials_dir :  credentials
 semantic_version_tags_enabled : false
 images_path :  stevedore.yaml
 log_path :  
 push_images :  true
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
		DEPRECATEDTreePathFile:    "test_stevedore.yaml",
		DEPRECATEDBuilderPath:     "test_stevedore.yaml",
		LogPathFile:               "test_stevedore.log",
		Concurrency:               8,
		PushImages:                true,
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
