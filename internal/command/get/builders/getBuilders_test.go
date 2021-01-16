package getbuilders

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/stretchr/testify/assert"
)

const (
	testBaseDir = "test"
)

func TestGetBuildersHandler(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	tests := []struct {
		desc   string
		ctx    context.Context
		config *configuration.Configuration
		args   []string
		res    map[string]int8
		err    error
	}{
		{
			desc: "Testing get images tree",
			ctx:  ctx,
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
				"builders",
				"-c",
				filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				"BUILDER        DRIVER           OPTIONS\n":                                   int8(0),
				"infrastructure ansible-playbook inventory=inventory/all playbook=site.yml\n": int8(0),
				"code           docker           context=map[path:.]\n":                       int8(0),
				"dummy          default\n":                                                    int8(0),
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			w.Reset()

			cmd := NewCommand(test.ctx, test.config)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
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
