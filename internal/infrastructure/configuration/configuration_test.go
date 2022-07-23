package configuration

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/loader"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
				DEPRECATEDTreePathFile:         filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile),
				DEPRECATEDBuilderPath:          filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath),
				LogPathFile:                    DefaultLogPathFile,
				DEPRECATEDNumWorkers:           DEPRECATEDDefaultNumWorker,
				Concurrency:                    concurrencyValue(),
				PushImages:                     DefaultPushImages,
				DEPRECATEDBuildOnCascade:       DEPRECATEDDefaultBuildOnCascade,
				DEPRECATEDDockerCredentialsDir: DEPRECATEDDefaultDockerCredentialsDir,
				EnableSemanticVersionTags:      DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates:   []string{DefaultSemanticVersionTagsTemplates},
			},
		},
		{
			desc: "Testing overwrite log_path on configuration file",
			config: []byte(`
log_path: "/var/log/stevedore/stevedore.log"
num_workers: 5 
`),
			res: &Configuration{
				DEPRECATEDTreePathFile:         filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultTreePathFile),
				DEPRECATEDBuilderPath:          filepath.Join(DefaultConfigFolder, DEPRECATEDDefaultBuilderPath),
				LogPathFile:                    "/var/log/stevedore/stevedore.log",
				DEPRECATEDNumWorkers:           5,
				Concurrency:                    concurrencyValue(),
				PushImages:                     DefaultPushImages,
				DEPRECATEDBuildOnCascade:       DEPRECATEDDefaultBuildOnCascade,
				DEPRECATEDDockerCredentialsDir: DEPRECATEDDefaultDockerCredentialsDir,
				EnableSemanticVersionTags:      DefaultEnableSemanticVersionTags,
				SemanticVersionTagsTemplates:   []string{DefaultSemanticVersionTagsTemplates},
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
		viper.SetDefault(DEPRECATEDDockerCredentialsDirKey, DEPRECATEDDefaultDockerCredentialsDir)
		viper.SetDefault(EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags)
		viper.SetDefault(SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates})

		viper.ReadConfig(bytes.NewBuffer(test.config))

		c := &Configuration{
			DEPRECATEDTreePathFile:         viper.GetString(DEPRECATEDTreePathFileKey),
			DEPRECATEDBuilderPath:          viper.GetString(DEPRECATEDBuilderPathKey),
			LogPathFile:                    viper.GetString(LogPathFileKey),
			DEPRECATEDNumWorkers:           viper.GetInt(DEPRECATEDNumWorkerKey),
			Concurrency:                    concurrencyValue(),
			PushImages:                     viper.GetBool(PushImagesKey),
			DEPRECATEDBuildOnCascade:       viper.GetBool(DEPRECATEDBuildOnCascadeKey),
			DEPRECATEDDockerCredentialsDir: viper.GetString(DEPRECATEDDockerCredentialsDirKey),
			EnableSemanticVersionTags:      viper.GetBool(EnableSemanticVersionTagsKey),
			SemanticVersionTagsTemplates:   viper.GetStringSlice(SemanticVersionTagsTemplatesKey),
		}

		assert.Equal(t, test.res, c, "Unpexpected configuration value")

	}
}

func TestNew(t *testing.T) {

	errContext := "(Configuration::New)"

	comp := &compatibility.MockCompatibility{}

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tests := []struct {
		desc              string
		res               *Configuration
		compatibility     Compatibilitier
		loader            ConfigurationLoader
		fs                afero.Fs
		prepareAssertFunc func(l ConfigurationLoader, c Compatibilitier)
		err               error
	}{
		{
			desc: "Testing error when loading configuration from file and compatibilitier is not provided",
			err:  errors.New(errContext, "Comptabilitier must be provided to create a new configuration"),
		},
		{
			desc:          "Testing error when loading configuration from file and filesystem is not provided",
			compatibility: compatibility.NewMockCompatibility(),
			err:           errors.New(errContext, "File system must be provided to create a new configuration"),
		},
		{
			desc:          "Testing error when loading configuration from file and configuration loader is not provided",
			compatibility: compatibility.NewMockCompatibility(),
			fs:            afero.NewMemMapFs(),
			err:           errors.New(errContext, "Configuration loader must be provided to create a new configuration"),
		},
		{
			desc:          "Testing create new configuration",
			fs:            afero.NewMemMapFs(),
			loader:        loader.NewMockConfigurationLoader(),
			compatibility: comp,
			prepareAssertFunc: func(l ConfigurationLoader, c Compatibilitier) {

				l.(*loader.MockConfigurationLoader).On("SetFs", afero.NewMemMapFs()).Return()
				l.(*loader.MockConfigurationLoader).On("AutomaticEnv").Return()
				l.(*loader.MockConfigurationLoader).On("SetEnvPrefix", "stevedore").Return()
				l.(*loader.MockConfigurationLoader).On("SetConfigName", DefaultConfigFile).Return()
				l.(*loader.MockConfigurationLoader).On("SetConfigType", DefaultConfigFileExtention).Return()

				l.(*loader.MockConfigurationLoader).On("SetDefault", BuildersPathKey, filepath.Join(DefaultConfigFolder, DefaultBuildersPath)).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", ConcurrencyKey, concurrencyValue()).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", EnableSemanticVersionTagsKey, DefaultEnableSemanticVersionTags).Return()

				l.(*loader.MockConfigurationLoader).On("SetDefault", ImagesPathKey, filepath.Join(DefaultConfigFolder, DefaultImagesPath)).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", LogPathFileKey, DefaultLogPathFile).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", PushImagesKey, DefaultPushImages).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", SemanticVersionTagsTemplatesKey, []string{DefaultSemanticVersionTagsTemplates}).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", strings.Join([]string{CredentialsKey, CredentialsStorageTypeKey}, "."), DefaultCredentialsStorage).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", strings.Join([]string{CredentialsKey, CredentialsLocalStoragePathKey}, "."), DefaultCredentialsLocalStoragePath).Return()
				l.(*loader.MockConfigurationLoader).On("SetDefault", strings.Join([]string{CredentialsKey, CredentialsFormatKey}, "."), DefaultCredentialsFormat).Return()

				l.(*loader.MockConfigurationLoader).On("AddConfigPath", filepath.Join(user.HomeDir, ".config", "stevedore")).Return()
				l.(*loader.MockConfigurationLoader).On("AddConfigPath", user.HomeDir).Return()
				l.(*loader.MockConfigurationLoader).On("AddConfigPath", DefaultConfigFolder).Return()

				l.(*loader.MockConfigurationLoader).On("ReadInConfig").Return(nil)

				l.(*loader.MockConfigurationLoader).On("GetString", LogPathFileKey).Return(DefaultLogPathFile)

				l.(*loader.MockConfigurationLoader).On("GetString", BuildersPathKey).Return(filepath.Join(DefaultConfigFolder, DefaultBuildersPath))
				l.(*loader.MockConfigurationLoader).On("GetInt", ConcurrencyKey).Return(concurrencyValue())
				l.(*loader.MockConfigurationLoader).On("GetBool", EnableSemanticVersionTagsKey).Return(DefaultEnableSemanticVersionTags)
				l.(*loader.MockConfigurationLoader).On("GetString", ImagesPathKey).Return(filepath.Join(DefaultConfigFolder, DefaultImagesPath))
				l.(*loader.MockConfigurationLoader).On("GetBool", PushImagesKey).Return(DefaultPushImages)
				l.(*loader.MockConfigurationLoader).On("GetStringSlice", SemanticVersionTagsTemplatesKey).Return([]string{DefaultSemanticVersionTagsTemplates})
				l.(*loader.MockConfigurationLoader).On("GetString", strings.Join([]string{CredentialsKey, CredentialsStorageTypeKey}, ".")).Return(DefaultCredentialsStorage)
				l.(*loader.MockConfigurationLoader).On("GetString", strings.Join([]string{CredentialsKey, CredentialsLocalStoragePathKey}, ".")).Return(DefaultCredentialsLocalStoragePath)
				l.(*loader.MockConfigurationLoader).On("GetString", strings.Join([]string{CredentialsKey, CredentialsFormatKey}, ".")).Return(DefaultCredentialsFormat)
			},
			res: &Configuration{
				ImagesPath:                   filepath.Join(".", "stevedore.yaml"),
				BuildersPath:                 filepath.Join(".", "stevedore.yaml"),
				LogPathFile:                  "",
				Concurrency:                  concurrencyValue(),
				PushImages:                   false,
				LogWriter:                    ioutil.Discard,
				EnableSemanticVersionTags:    false,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}.{{ .Minor }}.{{ .Patch }}"},
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "credentials",
					Format:           "json",
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.loader, test.compatibility)
			}

			c, err := New(test.fs, test.loader, test.compatibility)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {

				assert.Equal(t, test.res.BuildersPath, c.BuildersPath, "assert BuildersPath")
				assert.Equal(t, test.res.Concurrency, c.Concurrency, "assert Concurrency")
				assert.Equal(t, test.res.Credentials, c.Credentials, "assert Credentials")
				assert.Equal(t, test.res.EnableSemanticVersionTags, c.EnableSemanticVersionTags, "assert EnableSemanticVersionTags")
				assert.Equal(t, test.res.ImagesPath, c.ImagesPath, "assert ImagesPath")
				assert.Equal(t, test.res.LogWriter, c.LogWriter, "assert LogWriter")
				assert.Equal(t, test.res.LogWriter, c.LogWriter, "assert LogWriter")
				assert.Equal(t, test.res.PushImages, c.PushImages, "assert PushImages")
				assert.Equal(t, test.res.SemanticVersionTagsTemplates, c.SemanticVersionTagsTemplates, "assert SemanticVersionTagsTemplates")

				c.loader.(*loader.MockConfigurationLoader).AssertExpectations(t)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {

	var err error
	errContext := "(Configuration::LoadFromFile)"

	baseDir := "/config"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "stevedore.yaml"), []byte(`
builders_path: /config/stevedore.yaml
concurrency: 10
credentials:
  storage_type: local
  local_storage_path: mycredentials
  format: yaml
semantic_version_tags_enabled: true
images_path: /config/stevedore.yaml
log_path: mystevedore.log
push_images: false
semantic_version_tags_templates:
  - "{{ -Major }}"
  - "{{ -Major }}.{{ .Minor }}"
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "stevedore_emtpy.yaml"), []byte(`
# empty file
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "stevedore_deprecated.yaml"), []byte(`
builder_path: /config/stevedore.yaml
num_workers: 10
docker_registry_credentials_dir: mycredentials
build_on_cascade: true
semantic_version_tags_enabled: true
tree_path: /config/stevedore.yaml
log_path: mystevedore.log
push_images: false
semantic_version_tags_templates:
  - "{{ -Major }}"
  - "{{ -Major }}.{{ .Minor }}"
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		fs                afero.Fs
		file              string
		err               error
		res               *Configuration
		compatibility     Compatibilitier
		loader            ConfigurationLoader
		prepareAssertFunc func(l ConfigurationLoader, c Compatibilitier)
	}{
		{
			desc: "Testing error when loading configuration from file and compatibilitier is not provided",
			err:  errors.New(errContext, "Comptabilitier must be provided to create a new configuration"),
		},
		{
			desc:          "Testing error when loading configuration from file and filesystem is not provided",
			compatibility: compatibility.NewMockCompatibility(),
			err:           errors.New(errContext, "File system must be provided to create a new configuration"),
		},
		{
			desc:          "Testing error when loading configuration from file and configuration loader is not provided",
			compatibility: compatibility.NewMockCompatibility(),
			fs:            testFs,
			err:           errors.New(errContext, "Configuration loader must be provided to create a new configuration"),
		},
		{
			desc:          "Testing error when loading configuration from file and file does not exists",
			compatibility: compatibility.NewMockCompatibility(),
			fs:            testFs,
			loader:        &loader.MockConfigurationLoader{},
			file:          "unknown",
			err:           errors.New(errContext, "Configuration file could be loaded", errors.New(errContext, "testing error")),
			prepareAssertFunc: func(l ConfigurationLoader, c Compatibilitier) {
				l.(*loader.MockConfigurationLoader).On("SetFs", testFs).Return()
				l.(*loader.MockConfigurationLoader).On("SetConfigFile", "unknown").Return()
				l.(*loader.MockConfigurationLoader).On("ReadInConfig").Return(errors.New(errContext, "testing error"))
			},
		},
		{
			desc:   "Testing create new configuration from file",
			fs:     testFs,
			loader: loader.NewConfigurationLoader(viper.New()),
			file:   filepath.Join(baseDir, "stevedore.yaml"),
			err:    &errors.Error{},
			res: &Configuration{
				BuildersPath: "/config/stevedore.yaml",
				Concurrency:  10,
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "mycredentials",
					Format:           "yaml",
				},
				EnableSemanticVersionTags: true,
				ImagesPath:                "/config/stevedore.yaml",
				LogPathFile:               "mystevedore.log",
				PushImages:                false,
				SemanticVersionTagsTemplates: []string{
					"{{ -Major }}",
					"{{ -Major }}.{{ .Minor }}",
				},
			},
			compatibility: compatibility.NewMockCompatibility(),
		},
		{
			desc:   "Testing create new configuration from an empty file",
			fs:     testFs,
			loader: loader.NewConfigurationLoader(viper.New()),
			file:   filepath.Join(baseDir, "stevedore_emtpy.yaml"),
			err:    &errors.Error{},
			res: &Configuration{
				BuildersPath: "stevedore.yaml",
				Concurrency:  concurrencyValue(),
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "credentials",
					Format:           "json",
				},
				EnableSemanticVersionTags: false,
				ImagesPath:                "stevedore.yaml",
				LogPathFile:               "",
				PushImages:                false,
				SemanticVersionTagsTemplates: []string{
					"{{ .Major }}.{{ .Minor }}.{{ .Patch }}",
				},
			},
			compatibility: compatibility.NewMockCompatibility(),
		},
		{
			desc:   "Testing create new configuration from file with deprecated configuration",
			fs:     testFs,
			loader: loader.NewConfigurationLoader(viper.New()),
			file:   filepath.Join(baseDir, "stevedore_deprecated.yaml"),
			err:    &errors.Error{},
			res: &Configuration{
				BuildersPath: "/config/stevedore.yaml",
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "mycredentials",
					Format:           "json",
				},
				Concurrency:               10,
				EnableSemanticVersionTags: true,
				ImagesPath:                "/config/stevedore.yaml",
				LogPathFile:               "mystevedore.log",
				PushImages:                false,
				SemanticVersionTagsTemplates: []string{
					"{{ -Major }}",
					"{{ -Major }}.{{ .Minor }}",
				},
			},
			compatibility: compatibility.NewMockCompatibility(),
			prepareAssertFunc: func(l ConfigurationLoader, c Compatibilitier) {
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'tree_path' is deprecated and will be removed on v0.12.0, please use 'images_path' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'builder_path' is deprecated and will be removed on v0.12.0, please use 'builders_path' instead"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'num_workers' is deprecated and will be removed on v0.12.0, please use 'concurrency' instead"})
				c.(*compatibility.MockCompatibility).On("AddChanged", []string{"'build_on_cascade' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command"})
				c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'docker_registry_credentials_dir' is deprecated and will be removed on v0.12.0, please use 'credentials' block to configure credentials. Credentials local storage located in 'mycredentials' has precedence over 'credentials' block and is going to be used as default credentials store"})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.loader, test.compatibility)
			}

			config, err := LoadFromFile(test.fs, test.loader, test.file, test.compatibility)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res.BuildersPath, config.BuildersPath, "assert BuildersPath")
				assert.Equal(t, test.res.Concurrency, config.Concurrency, "assert Concurrency")
				assert.Equal(t, test.res.Credentials, config.Credentials, "assert Credentials")
				assert.Equal(t, test.res.EnableSemanticVersionTags, config.EnableSemanticVersionTags, "assert EnableSemanticVersionTags")
				assert.Equal(t, test.res.ImagesPath, config.ImagesPath, "assert ImagesPath")
				assert.Equal(t, test.res.LogPathFile, config.LogPathFile, "assert LogPathFile")
				assert.Equal(t, test.res.PushImages, config.PushImages, "assert PushImages")
				assert.Equal(t, test.res.SemanticVersionTagsTemplates, config.SemanticVersionTagsTemplates, "assert SemanticVersionTagsTemplates")
			}
		})
	}
}

func TestReloadConfigurationFromFile(t *testing.T) {
	var err error

	errContext := "(Configuration::ReloadConfigurationFromFile)"

	baseDir := "/config"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)
	err = afero.WriteFile(testFs, filepath.Join(baseDir, "stevedore.yaml"), []byte(`
builders_path: /config/mystevedore.yaml
concurrency: 10
credentials:
  storage_type: local
  local_storage_path: mycredentials
semantic_version_tags_enabled: true
images_path: /config/mystevedore.yaml
log_path: mystevedore.log
push_images: true
semantic_version_tags_templates:
  - "{{ .Major }}"
`), 0644)
	if err != nil {
		t.Log(err)
	}
	err = afero.WriteFile(testFs, filepath.Join(baseDir, "mystevedore.yaml"), []byte(`
	`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		config            *Configuration
		file              string
		res               *Configuration
		prepareAssertFunc func(c *Configuration)
		err               error
	}{
		{
			desc:   "Testing error when reloading configuration from file with missing file",
			config: &Configuration{},
			err:    errors.New(errContext, "Configuration file must be provided to reload configuration from file"),
		},
		{
			desc: "Testing reloading configuration from file",
			config: &Configuration{
				fs:            testFs,
				loader:        loader.NewConfigurationLoader(viper.New()),
				compatibility: compatibility.NewMockCompatibility(),
			},

			file: filepath.Join(baseDir, "stevedore.yaml"),
			res: &Configuration{
				BuildersPath: "/config/mystevedore.yaml",
				Concurrency:  10,
				ImagesPath:   "/config/mystevedore.yaml",
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "mycredentials",
					Format:           "json",
				},
				LogPathFile:                  "mystevedore.log",
				PushImages:                   true,
				EnableSemanticVersionTags:    true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.config)
			}

			err := test.config.ReloadConfigurationFromFile(test.file)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res.BuildersPath, test.config.BuildersPath, "assert BuildersPath")
				assert.Equal(t, test.res.Concurrency, test.config.Concurrency, "assert Concurrency")
				assert.Equal(t, test.res.Credentials, test.config.Credentials, "assert Credentials")
				assert.Equal(t, test.res.EnableSemanticVersionTags, test.config.EnableSemanticVersionTags, "assert EnableSemanticVersionTags")
				assert.Equal(t, test.res.ImagesPath, test.config.ImagesPath, "assert ImagesPath")
				assert.Equal(t, test.res.PushImages, test.config.PushImages, "assert PushImages")
				assert.Equal(t, test.res.SemanticVersionTagsTemplates, test.config.SemanticVersionTagsTemplates, "assert SemanticVersionTagsTemplates")
			}
		})
	}
}

func TestCheckCompatibility(t *testing.T) {

	errContext := "(Configuration::CheckCompatibility)"

	tests := []struct {
		desc              string
		config            *Configuration
		res               *Configuration
		prepareAssertFunc func(c *Configuration)
		err               error
	}{
		{
			desc:   "Testing error checking configuration compatibility when compatibilitier is not defined",
			config: &Configuration{},
			err:    errors.New(errContext, "To ckeck configuration compatiblity is required a compatibilitier"),
		},
		{
			desc: "Testing check configuration with deprecated configuration",
			config: &Configuration{
				DEPRECATEDBuilderPath:          "mystevedore.yaml",
				DEPRECATEDBuildOnCascade:       true,
				DEPRECATEDNumWorkers:           10,
				DEPRECATEDTreePathFile:         "mystevedore.yaml",
				DEPRECATEDDockerCredentialsDir: "mycredentials",

				LogPathFile:                  "mystevedore.log",
				PushImages:                   true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},

				compatibility: compatibility.NewMockCompatibility(),
			},
			res: &Configuration{
				BuildersPath: "mystevedore.yaml",
				Concurrency:  10,
				ImagesPath:   "mystevedore.yaml",
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "mycredentials",
					Format:           "json",
				},
				LogPathFile:                  "mystevedore.log",
				PushImages:                   true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},

				compatibility: compatibility.NewMockCompatibility(),
			},
			prepareAssertFunc: func(c *Configuration) {
				c.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'tree_path' is deprecated and will be removed on v0.12.0, please use 'images_path' instead"})
				c.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'builder_path' is deprecated and will be removed on v0.12.0, please use 'builders_path' instead"})
				c.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'num_workers' is deprecated and will be removed on v0.12.0, please use 'concurrency' instead"})
				c.compatibility.(*compatibility.MockCompatibility).On("AddChanged", []string{"'build_on_cascade' is not available anymore as a configuration parameter. Cascade execution plan is only enabled by '--cascade' flag on build command"})
				c.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'docker_registry_credentials_dir' is deprecated and will be removed on v0.12.0, please use 'credentials' block to configure credentials. Credentials local storage located in 'mycredentials' has precedence over 'credentials' block and is going to be used as default credentials store"})
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.config)
			}

			err := test.config.CheckCompatibility()

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res.BuildersPath, test.config.BuildersPath, "assert BuildersPath")
				assert.Equal(t, test.res.Concurrency, test.config.Concurrency, "assert Concurrency")
				assert.Equal(t, test.res.Credentials, test.config.Credentials, "assert Credentials")
				assert.Equal(t, test.res.EnableSemanticVersionTags, test.config.EnableSemanticVersionTags, "assert EnableSemanticVersionTags")
				assert.Equal(t, test.res.ImagesPath, test.config.ImagesPath, "assert ImagesPath")
				assert.Equal(t, test.res.LogWriter, test.config.LogWriter, "assert LogWriter")
				assert.Equal(t, test.res.LogWriter, test.config.LogWriter, "assert LogWriter")
				assert.Equal(t, test.res.PushImages, test.config.PushImages, "assert PushImages")
				assert.Equal(t, test.res.SemanticVersionTagsTemplates, test.config.SemanticVersionTagsTemplates, "assert SemanticVersionTagsTemplates")
			}
		})
	}
}

func TestValidateConfiguration(t *testing.T) {
	var err error
	errContext := "(Configuration::ValidateConfiguration)"

	baseDir := "/config"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)
	err = afero.WriteFile(testFs, filepath.Join(baseDir, "mystevedore.yaml"), []byte(`
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc   string
		config *Configuration
		err    error
	}{
		{
			desc: "Testing validate configuration when everything is ok",
			config: &Configuration{
				BuildersPath: filepath.Join(baseDir, "mystevedore.yaml"),
				Concurrency:  10,
				ImagesPath:   filepath.Join(baseDir, "mystevedore.yaml"),
				Credentials: &CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "mycredentials",
					Format:           "json",
				},
				LogPathFile:                  "mystevedore.log",
				PushImages:                   true,
				EnableSemanticVersionTags:    true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
				fs:                           testFs,
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing error when file system is not defined on configuration",
			config: &Configuration{},
			err:    errors.New(errContext, "File system must be provided to create a new configuration"),
		},
		{
			desc: "Testing error when builders path is not defined",
			config: &Configuration{
				fs: testFs,
			},
			err: errors.New(errContext, "Invalid configuration, builders path must be provided"),
		},
		{
			desc: "Testing error when images path is not defined",
			config: &Configuration{
				BuildersPath: filepath.Join(baseDir, "mystevedore.yaml"),
				fs:           testFs,
			},
			err: errors.New(errContext, "Invalid configuration, images path must be provided"),
		},
		{
			desc: "Testing error when currency is lower than 1",
			config: &Configuration{
				BuildersPath: filepath.Join(baseDir, "mystevedore.yaml"),
				ImagesPath:   filepath.Join(baseDir, "mystevedore.yaml"),
				Concurrency:  0,
				fs:           testFs,
			},
			err: errors.New(errContext, "Invalid configuration, concurrency must be greater than 0"),
		},
		{
			desc: "Testing error when credentials storage type is not defined",
			config: &Configuration{
				BuildersPath: filepath.Join(baseDir, "mystevedore.yaml"),
				ImagesPath:   filepath.Join(baseDir, "mystevedore.yaml"),
				Concurrency:  1,
				Credentials:  &CredentialsConfiguration{},
				fs:           testFs,
			},
			err: errors.New(errContext, "Invalid configuration, credentials storage type must be provided"),
		},
		{
			desc: "Testing error when credentials format is not defined",
			config: &Configuration{
				BuildersPath: filepath.Join(baseDir, "mystevedore.yaml"),
				ImagesPath:   filepath.Join(baseDir, "mystevedore.yaml"),
				Concurrency:  1,
				Credentials: &CredentialsConfiguration{
					StorageType: "local",
				},
				fs: testFs,
			},
			err: errors.New(errContext, "Invalid configuration, credentials format must be provided"),
		},
		{
			desc: "Testing error when credentials local storage path is not defined when is used credentials local storage type",
			config: &Configuration{
				BuildersPath: filepath.Join(baseDir, "mystevedore.yaml"),
				ImagesPath:   filepath.Join(baseDir, "mystevedore.yaml"),
				Concurrency:  1,
				Credentials: &CredentialsConfiguration{
					StorageType: "local",
					Format:      "json",
				},
				fs: testFs,
			},
			err: errors.New(errContext, "Invalid configuration, credentials local storage path must be provided"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.config.ValidateConfiguration()

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}
