package getcredentials

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"stevedore/internal/configuration"
	"stevedore/internal/credentials"
	"stevedore/internal/ui/console"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBaseDir = "test"
)

func TestGetCredentialsHandler(t *testing.T) {

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
			desc: "Testing get credentials",
			ctx:  ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "credentials"),
			},
			args: []string{
				"get",
				"credentials",
				"--config",
				filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				"CREDENTIAL ID                    USERNAME\n":  int8(0),
				"91fc14ad02afd60985bb8165bda320a6 username1\n": int8(0),
				"b1946ac92492d2347c6235b4d2611184 username2\n": int8(0),
			},
			err: nil,
		},
		{
			desc: "Testing get credentials with wide options",
			ctx:  ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "credentials"),
			},
			args: []string{
				"get",
				"credentials",
				"--wide",
				"--config",
				filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				"CREDENTIAL ID                    USERNAME  PASSWORD\n": int8(0),
				"b1946ac92492d2347c6235b4d2611184 username2 password\n": int8(0),
				"91fc14ad02afd60985bb8165bda320a6 username1 password\n": int8(0),
			},
			err: nil,
		},
	}

	var err error
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			w.Reset()
			credentials.ClearCredentials()
			err = credentials.LoadCredentials(test.config.DockerCredentialsDir)
			if err != nil {
				t.Error(err.Error())
			}
			cmd := NewCommand(test.ctx, test.config)
			cmd.Command.ParseFlags(test.args)
			err = cmd.Command.RunE(cmd.Command, test.args)

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
