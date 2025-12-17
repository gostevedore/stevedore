package file

import (
	"io/ioutil"
	"testing"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {

	var err error
	var content, expected []byte

	testFs := afero.NewMemMapFs()

	output := NewConfigurationFilePersist(
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
			EncryptionKey:    "encryptionkey",
		},
		LogPathFile:                  "mystevedore.log",
		PushImages:                   true,
		EnableSemanticVersionTags:    true,
		SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
	}

	expected, err = ioutil.ReadFile("test/stevedore.yaml.golden")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = output.Write(config)
	if err != nil {
		t.Errorf("%v", err)
	}

	a := afero.Afero{
		Fs: testFs,
	}
	content, err = a.ReadFile("test.yaml")
	if err != nil {
		t.Errorf("%v", err)
	}

	// Uncomment to generate new golden file
	// f, err := os.OpenFile("test/conf", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// }
	// defer f.Close()
	// f.Write(content)

	assert.Equal(t, expected, content, "Unexpected response")
}
