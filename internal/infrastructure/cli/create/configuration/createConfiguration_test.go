package configuration

import (
	"context"
	"testing"

	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		desc            string
		entrypoint      Entrypointer
		prepareMockFunc func(Entrypointer)
		args            []string
		err             error
	}{
		{
			desc:       "Testing run create configuration command",
			entrypoint: entrypoint.NewMockCreateConfigurationEntrypoint(),
			args: []string{
				"--builders-path",
				"/builders",
				"--concurrency",
				"4",
				"--config",
				"/stevedore-config.yaml",
				"--credentials-format",
				"json",
				"--credentials-local-storage-path",
				"/credentials",
				"--credentials-storage-type",
				"local",
				"--enable-semver-tags",
				"--force",
				"--images-path",
				"/images",
				"--log-path-file",
				"/logs",
				"--push-images",
				"--semver-tags-template",
				"{{ .Major }}",
				"--semver-tags-template",
				"{{ .Major }}_{{ .Minor }}",
			},
			prepareMockFunc: func(e Entrypointer) {
				e.(*entrypoint.MockCreateConfigurationEntrypoint).On(
					"Execute",
					context.TODO(),
					&entrypoint.Options{
						BuildersPath:                "/builders",
						Concurrency:                 4,
						ConfigurationFilePath:       "/stevedore-config.yaml",
						CredentialsFormat:           "json",
						CredentialsLocalStoragePath: "/credentials",
						CredentialsStorageType:      "local",
						EnableSemanticVersionTags:   true,
						Force:                       true,
						ImagesPath:                  "/images",
						LogPathFile:                 "/logs",
						PushImages:                  true,
						SemanticVersionTagsTemplates: []string{
							"{{ .Major }}",
							"{{ .Major }}_{{ .Minor }}",
						},
					},
				).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.entrypoint)
			}

			cmd := NewCommand(context.TODO(), test.entrypoint)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}
