package file

import (
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSaveWriteFileExistError(t *testing.T) {

	var err error

	errContext := "(configuration::output::ConfigurationFileSafePersist::Write)"

	testFs := afero.NewMemMapFs()
	err = afero.WriteFile(testFs, "test.yaml", []byte{}, 0644)
	if err != nil {
		t.Errorf(err.Error())
	}

	output := NewConfigurationFileSafePersist(
		WithFileSystem(testFs),
		WithFilePath("test.yaml"),
	)

	config := &configuration.Configuration{
		BuildersPath: "mystevedore.yaml",
		Concurrency:  10,
		ImagesPath:   "mystevedore.yaml",
		Credentials: &configuration.CredentialsConfiguration{
			StorageType:      "local",
			LocalStoragePath: "mycredentials",
			Format:           "json",
		},
		LogPathFile:                  "mystevedore.log",
		PushImages:                   true,
		EnableSemanticVersionTags:    true,
		SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
	}

	expected := errors.New(errContext, "Configuration file 'test.yaml' already exist and will not be created")

	err = output.Write(config)
	assert.Equal(t, expected, err)
}

func TestWriteToFile(t *testing.T) {

	var err error
	var content, expected []byte

	testFs := afero.NewMemMapFs()

	output := NewConfigurationFileSafePersist(
		WithFileSystem(testFs),
		WithFilePath("test.yaml"),
	)

	config := &configuration.Configuration{
		BuildersPath: "mystevedore.yaml",
		Concurrency:  10,
		ImagesPath:   "mystevedore.yaml",
		Credentials: &configuration.CredentialsConfiguration{
			StorageType:      "local",
			LocalStoragePath: "mycredentials",
			Format:           "json",
		},
		LogPathFile:                  "mystevedore.log",
		PushImages:                   true,
		EnableSemanticVersionTags:    true,
		SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
	}

	expected, err = ioutil.ReadFile("test/stevedore.yaml.golden")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = output.Write(config)
	if err != nil {
		t.Errorf(err.Error())
	}

	a := afero.Afero{
		Fs: testFs,
	}
	content, err = a.ReadFile("test.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Uncomment to generate new golden file
	// f, err := os.OpenFile("test/conf", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// }
	// defer f.Close()
	// f.Write(buff.Bytes())

	assert.Equal(t, expected, content, "Unexpected response")
}
