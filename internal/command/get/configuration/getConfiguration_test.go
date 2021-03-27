package getconfiguration

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/stretchr/testify/assert"
)

const (
	testBaseDir = "test"
)

func TestGetConfigurationHandler(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))

	ctx := context.TODO()

	tests := []struct {
		desc   string
		ctx    context.Context
		config *configuration.Configuration
		args   []string
		res    string
		err    error
	}{
		{
			desc:   "Testing get configuration from empty configuration",
			ctx:    ctx,
			config: &configuration.Configuration{},
			args:   []string{},
			res: `PARAMETER                       VALUE
tree_path                       
builder_path                    
log_path                        
num_workers                     0
push_images                     false
build_on_cascade                false
docker_registry_credentials_dir 
semantic_version_tags_enabled   false
semantic_version_tags_templates []
`,
			err: nil,
		},
		{
			desc: "Testing get configuration",
			ctx:  ctx,
			config: &configuration.Configuration{
				TreePathFile:                 "test/stevedore_config.yml",
				BuilderPathFile:              "test/stevedore_config.yml",
				LogPathFile:                  "/dev/null",
				NumWorkers:                   1,
				PushImages:                   false,
				BuildOnCascade:               false,
				DockerCredentialsDir:         "test/stevedore_config.yml",
				EnableSemanticVersionTags:    false,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}", "{{ .Major }}.{{ .Minor }}"},
			},
			args: []string{
				"get",
				"config",
				"-c",
				"test/stevedore_config.yml",
			},
			res: `PARAMETER                       VALUE
tree_path                       test/stevedore_config.yml
builder_path                    test/stevedore_config.yml
log_path                        /dev/null
num_workers                     1
push_images                     false
build_on_cascade                false
docker_registry_credentials_dir test/stevedore_config.yml
semantic_version_tags_enabled   false
semantic_version_tags_templates [{{ .Major }} {{ .Major }}.{{ .Minor }}]
`,
			err: nil,
		},
	}

	var err error
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			w.Reset()
			t.Log(test.desc)

			cmd := NewCommand(test.ctx, test.config)
			cmd.Command.ParseFlags(test.args)
			err = cmd.Command.RunE(cmd.Command, test.args)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, w.String(), "Unexpected value")
			}
		})
	}
}
