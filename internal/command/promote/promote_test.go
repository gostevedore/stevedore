package promote

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"stevedore/internal/configuration"
	"stevedore/internal/ui/console"
	"strings"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

const (
	testBaseDir = "test"
)

func TestPromoteHandler(t *testing.T) {
	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	tests := []struct {
		desc    string
		skip    bool
		verbose bool
		ctx     context.Context
		config  *configuration.Configuration
		args    []string
		res     map[string]int8
		err     error
	}{
		{
			desc:    "Testing to promote a simple image",
			verbose: false,
			skip:    false,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			err: nil,
			args: []string{
				"--dry-run",
				"myregistryhost.com/namespace/ubuntu:20.04",
			},
			res: map[string]int8{
				"{DryRun:true EnableSemanticVersionTags:false ImagePromoteName: ImagePromoteRegistryNamespace: ImagePromoteRegistryHost: ImagePromoteTags:[] RemovePromotedTags:false ImageName:myregistryhost.com/namespace/ubuntu:20.04 OutputPrefix: SemanticVersionTagsTemplate:[]}": int8(0),
			},
		},
		{
			desc:    "Testing to promote an image to a new registry host, registry namespace, with new name and multiple tags",
			verbose: false,
			skip:    false,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			err: nil,
			args: []string{
				"--dry-run",
				"myregistryhost.com/namespace/ubuntu:20.04",
				"--promote-image-name",
				"myubuntu",
				"--promote-image-namespace",
				"stable",
				"--promote-image-registry",
				"myprodregistryhost.com",
				"--promote-image-tag",
				"tag1",
				"--promote-image-tag",
				"tag2",
				"--remove-promote-tags",
			},
			res: map[string]int8{
				"{DryRun:true EnableSemanticVersionTags:false ImagePromoteName:myubuntu ImagePromoteRegistryNamespace:stable ImagePromoteRegistryHost:myprodregistryhost.com ImagePromoteTags:[tag1 tag2] RemovePromotedTags:true ImageName:myregistryhost.com/namespace/ubuntu:20.04 OutputPrefix: SemanticVersionTagsTemplate:[]}": int8(0),
			},
		},
		{
			desc:    "Testing to promote without image name",
			verbose: true,
			skip:    false,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           false,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			err: errors.New("(command::promoteHandler)", "Is required an image name"),
			args: []string{
				"--dry-run",
			},
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
				assert.Equal(t, test.err, err)
			} else {
				if test.verbose {
					t.Log("\n", w.String())
				}

				wSplit := strings.Split(w.String(), "\n")
				assert.Equal(t, len(test.res), len(wSplit)-1, "Unexpected number of lines")
				for i := 0; i < len(wSplit)-1; i++ {

					line := wSplit[i]
					_, ok := test.res[line]

					assert.True(t, ok)
					delete(test.res, line)
				}
				assert.Equal(t, len(test.res), 0, "Not all expected lines has appeared")
			}

		})
	}

}
