package console

import (
	"bytes"
	"testing"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	var buff bytes.Buffer

	config := &configuration.Configuration{
		BuildersPath: "mystevedore.yaml",
		Concurrency:  10,
		ImagesPath:   "mystevedore.yaml",
		Credentials: &configuration.CredentialsConfiguration{
			StorageType:      "local",
			LocalStoragePath: "mycredentials",
			Format:           "json",
		},
		LogPathFile:                  "/log/mystevedore.log",
		PushImages:                   true,
		EnableSemanticVersionTags:    true,
		SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
	}

	expected := ` builders_path: mystevedore.yaml
 concurrency: 10
 semantic_version_tags_enabled: true
 images_path: mystevedore.yaml
 log_path: /log/mystevedore.log
 push_images: true
 semantic_version_tags_templates:
   - {{ .Major }}
 credentials:
   storage_type: local
   format: json
   local_storage_path: mycredentials
`

	console := NewConfigurationConsoleOutput(&buff)
	console.Write(config)
	assert.Equal(t, expected, buff.String())
}
