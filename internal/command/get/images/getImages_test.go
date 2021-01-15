package getimages

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/internal/configuration"

	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

const (
	testBaseDir = "test"
)

func TestGetImagesHandler(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	globalSkip := false
	globalVerbose := false

	tests := []struct {
		desc    string
		ctx     context.Context
		config  *configuration.Configuration
		skip    bool
		verbose bool
		args    []string
		res     map[string]int8
		err     error
	}{
		{
			desc:    "Testing get images tree",
			ctx:     ctx,
			skip:    globalSkip,
			verbose: globalVerbose,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			args: []string{
				"get",
				"images",
				"-t",
				"-c",
				filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				" \u251C\u2500\u2500\u2500 ubuntu:16.04\n":                                int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 nginx:1.15-ubuntu16.04\n":              int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 php-fpm:7.4-ubuntu16.04\n":             int8(0),
				" \u2502  \u2502  \u251C\u2500\u2500\u2500 php-fpm-dev:7.4-ubuntu16.04\n": int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 php-cli:7.4-ubuntu16.04\n":             int8(0),
				" \u2502  \u2502  \u251C\u2500\u2500\u2500 php-cli-dev:7.4-ubuntu16.04\n": int8(0),
				" \u251C\u2500\u2500\u2500 apps:master\n":                                 int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 app1:master\n":                         int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 app2:master\n":                         int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 app3:master\n":                         int8(0),
				" \u251C\u2500\u2500\u2500 ubuntu:18.04\n":                                int8(0),
				" \u2502  \u251C\u2500\u2500\u2500 php-fpm:7.4-ubuntu18.04\n":             int8(0),
				" \u2502  \u2502  \u251C\u2500\u2500\u2500 php-fpm-dev:7.4-ubuntu18.04\n": int8(0),
			},
			err: nil,
		},
		{
			desc:    "Testing get images",
			ctx:     ctx,
			skip:    globalSkip,
			verbose: globalVerbose,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			args: []string{
				"get",
				"images",
				"-c",
				filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				"NAME        VERSION          BUILDER        NAMESPACE REGISTRY  PARENT\n":                  int8(0),
				"nginx       1.15-ubuntu16.04 infrastructure                     ubuntu:16.04\n":            int8(0),
				"php-fpm-dev 7.4-ubuntu16.04  infrastructure                     php-fpm:7.4-ubuntu16.04\n": int8(0),
				"php-cli-dev 7.4-ubuntu16.04  infrastructure                     php-cli:7.4-ubuntu16.04\n": int8(0),
				"php-fpm     7.4-ubuntu18.04  infrastructure                     ubuntu:18.04\n":            int8(0),
				"app2        master           php-code                           apps:master\n":             int8(0),
				"php-fpm     7.4-ubuntu16.04  infrastructure                     ubuntu:16.04\n":            int8(0),
				"app3        master           php-code                           apps:master\n":             int8(0),
				"php-cli     7.4-ubuntu16.04  infrastructure                     ubuntu:16.04\n":            int8(0),
				"ubuntu      18.04            infrastructure           registry  -\n":                       int8(0),
				"php-fpm-dev 7.4-ubuntu18.04  infrastructure                     php-fpm:7.4-ubuntu18.04\n": int8(0),
				"apps        master           dummy                              -\n":                       int8(0),
				"app1        master           php-code                           apps:master\n":             int8(0),
				"ubuntu      16.04            infrastructure           registryX -\n":                       int8(0),
			},
			err: nil,
		},

		{
			desc:    "Testing compatibility",
			ctx:     ctx,
			skip:    true,
			verbose: true,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config_compatibility.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config_compatibility.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config_compatibility.yml"),
			},
			args: []string{
				"get",
				"images",
				"-c",
				filepath.Join(testBaseDir, "stevedore_config_compatibility.yml"),
			},
			res: map[string]int8{
				"NAME        VERSION          BUILDER        NAMESPACE REGISTRY  PARENT\n":                  int8(0),
				"nginx       1.15-ubuntu16.04 infrastructure                     ubuntu:16.04\n":            int8(0),
				"php-fpm-dev 7.4-ubuntu16.04  infrastructure                     php-fpm:7.4-ubuntu16.04\n": int8(0),
				"php-cli-dev 7.4-ubuntu16.04  infrastructure                     php-cli:7.4-ubuntu16.04\n": int8(0),
				"php-fpm     7.4-ubuntu18.04  infrastructure                     ubuntu:18.04\n":            int8(0),
				"app2        master           php-code                           apps:master\n":             int8(0),
				"php-fpm     7.4-ubuntu16.04  infrastructure                     ubuntu:16.04\n":            int8(0),
				"app3        master           php-code                           apps:master\n":             int8(0),
				"php-cli     7.4-ubuntu16.04  infrastructure                     ubuntu:16.04\n":            int8(0),
				"ubuntu      18.04            infrastructure           registry  -\n":                       int8(0),
				"php-fpm-dev 7.4-ubuntu18.04  infrastructure                     php-fpm:7.4-ubuntu18.04\n": int8(0),
				"apps        master           dummy                              -\n":                       int8(0),
				"app1        master           php-code                           apps:master\n":             int8(0),
				"ubuntu      16.04            infrastructure           registryX -\n":                       int8(0),
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.skip {
				t.Skip(test.desc)
			}

			w.Reset()

			cmd := NewCommand(test.ctx, test.config)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)

			if err != nil && assert.Error(t, err) {
				if test.verbose {
					t.Log("\n verbose:\n", err.Error())
				}
				assert.Equal(t, test.err, err)
			} else {
				if test.verbose {
					t.Log("\n verbose:\n", w.String())
				}

				wSplit := strings.Split(w.String(), "\n")
				// len(wSplit)-1 because the output finishes with a \n it generates an extra line on wSplit
				assert.Equal(t, len(test.res), len(wSplit)-1, "Unexpected number of lines")
				for i := 0; i < w.Len(); i++ {
					line, _ := w.ReadString('\n')
					_, ok := test.res[line]
					assert.True(t, ok)
					delete(test.res, line)
				}

				assert.Equal(t, len(test.res), 0, "Not all expected lines has appeared")
			}

		})
	}
}
