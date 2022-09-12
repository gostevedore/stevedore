package file

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestWriteConfigurationFile(t *testing.T) {

	var err error
	var buff bytes.Buffer
	var expected []byte

	output := NewConfigurationFileOutput(&buff)

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

	// Uncomment to generate new golden file
	// f, err := os.OpenFile("test/conf", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// }
	// defer f.Close()
	// f.Write(buff.Bytes())

	assert.Equal(t, expected, buff.Bytes(), "Unexpected response")
}
