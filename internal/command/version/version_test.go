package version

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"stevedore/internal/configuration"
	"stevedore/internal/release"
	"stevedore/internal/ui/console"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBaseDir = "test"
)

func TestVersionHandler(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	tests := []struct {
		desc      string
		ctx       context.Context
		config    *configuration.Configuration
		args      []string
		version   string
		commit    string
		buildDate string
		osarch    string
		res       string
		err       error
	}{
		{
			desc: "Testing version",
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
			args: []string{},
			res: `Stevedore 1.1.1 Commit: abcdfg linux/amd64 BuildDate: Thu Mar  3 23:05:25 2005
`,
			version:   "1.1.1",
			commit:    "abcdfg",
			buildDate: "Thu Mar  3 23:05:25 2005",
			osarch:    "linux/amd64",
			err:       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			w.Reset()

			release.Version = test.version
			release.Commit = test.commit
			release.BuildDate = test.buildDate
			release.OsArch = test.osarch

			cmd := NewCommand(test.ctx, test.config)
			err := cmd.Command.RunE(cmd.Command, test.args)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, w.String(), "Received unexpected version")
			}

		})
	}

}
