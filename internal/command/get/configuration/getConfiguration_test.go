package getconfiguration

import (
	"bytes"
	"context"
	"io"
	"stevedore/internal/configuration"
	"stevedore/internal/ui/console"
	"testing"

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
				TreePathFile:         "/treepathfile",
				BuilderPathFile:      "/builderpathfile",
				LogPathFile:          "/logpathfile",
				NumWorkers:           5,
				PushImages:           false,
				BuildOnCascade:       true,
				DockerCredentialsDir: "/dockercredentialsdir",
			},
			args: []string{
				"get",
				"config",
				"-c",
				"files/config/stevedore_reload.yaml",
			},
			res: `PARAMETER                       VALUE
tree_path                       /treepathfile
builder_path                    /builderpathfile
log_path                        /logpathfile
num_workers                     5
push_images                     false
build_on_cascade                true
docker_registry_credentials_dir /dockercredentialsdir
semantic_version_tags_enabled   false
semantic_version_tags_templates []
`,
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			w.Reset()

			cmd := NewCommand(test.ctx, test.config)
			err := cmd.Command.RunE(cmd.Command, test.args)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, w.String(), "Unexpected value")
			}
		})
	}
}
