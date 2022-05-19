package promote

import (
	"context"
	"fmt"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/compatibility"
	"github.com/gostevedore/stevedore/internal/configuration"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {

	tests := []struct {
		desc            string
		handler         HandlerPromoter
		compatibility   Compatibilitier
		entrypoint      Entrypointer
		config          *configuration.Configuration
		args            []string
		prepareMockFunc func(Compatibilitier, Entrypointer, *configuration.Configuration)
		err             error
	}{
		{
			desc:          "Testing run promote command",
			handler:       handler.NewHandlerMock(),
			compatibility: compatibility.NewMockCompatibility(),
			config:        &configuration.Configuration{},
			entrypoint:    entrypoint.NewMockEntrypoint(),
			args: []string{
				"source-registry-host.com/source-namespace/source-image:source-tag",
				"--enable-semver-tags",
				"--semver-tags-template",
				"{{ .Major }}",
				"--dry-run",
				"--promote-image-name",
				"promote-image-name",
				"--promote-image-registry-host",
				"promote-registry-host.com",
				"--promote-image-registry-namespace",
				"promote-registry-namespace",
				"--promote-image-tag",
				"promote-image-tag",
				"--remove-local-images-after-push",
				"--force-promote-source-image",
				"--image-from-remote-source",
			},
			prepareMockFunc: func(compatibility Compatibilitier, promote Entrypointer, config *configuration.Configuration) {

				options := &handler.Options{
					DryRun:                       true,
					EnableSemanticVersionTags:    true,
					TargetImageName:              "promote-image-name",
					TargetImageRegistryNamespace: "promote-registry-namespace",
					TargetImageRegistryHost:      "promote-registry-host.com",
					TargetImageTags:              []string{"promote-image-tag"},
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
					PromoteSourceImageTag:        true,
					RemoteSourceImage:            true,
				}

				promote.(*entrypoint.MockEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{
						"source-registry-host.com/source-namespace/source-image:source-tag",
					},
					config,
					options,
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:          "Testing run promote command with deprecated commands",
			handler:       handler.NewHandlerMock(),
			compatibility: compatibility.NewMockCompatibility(),
			config:        &configuration.Configuration{},
			entrypoint:    entrypoint.NewMockEntrypoint(),
			args: []string{
				"source-registry-host.com/source-namespace/source-image:source-tag",
				"--enable-semver-tags",
				"--semver-tags-template",
				"{{ .Major }}",
				"--dry-run",
				"--promote-image-name",
				"promote-image-name",
				"--promote-image-registry-host",
				"promote-registry-host.com",
				"--promote-image-registry-namespace",
				"promote-registry-namespace",
				"--promote-image-tag",
				"promote-image-tag",
				"--remove-promote-tags",
				"--force-promote-source-image",
				"--image-from-remote-source",
			},
			prepareMockFunc: func(comp Compatibilitier, promote Entrypointer, config *configuration.Configuration) {

				options := &handler.Options{
					DryRun:                       true,
					EnableSemanticVersionTags:    true,
					TargetImageName:              "promote-image-name",
					TargetImageRegistryNamespace: "promote-registry-namespace",
					TargetImageRegistryHost:      "promote-registry-host.com",
					TargetImageTags:              []string{"promote-image-tag"},
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
					PromoteSourceImageTag:        true,
					RemoteSourceImage:            true,
				}

				promote.(*entrypoint.MockEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{
						"source-registry-host.com/source-namespace/source-image:source-tag",
					},
					config,
					options,
				).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageRemoveTargetImageTags}).Return(nil)
			},
			err: &errors.Error{},
		},

		// {
		// 	desc:          "Testing to promote an image to a new registry host, registry namespace, with new name and multiple tags",
		// 	handler:       handler.NewHandlerMock(),
		// 	compatibility: compatibility.NewMockCompatibility(),
		// 	config:        &configuration.Configuration{},
		// 	entrypoint:    entrypoint.NewMockEntrypoint(),
		// 	args: []string{
		// 		"--dry-run",
		// 		"myregistryhost.com/namespace/ubuntu:20.04",
		// 		"--promote-image-name",
		// 		"myubuntu",
		// 		"--promote-image-namespace",
		// 		"stable",
		// 		"--promote-image-registry",
		// 		"myprodregistryhost.com",
		// 		"--promote-image-tag",
		// 		"tag1",
		// 		"--promote-image-tag",
		// 		"tag2",
		// 		"--remove-local-images-after-push",
		// 	},
		// 	prepareMockFunc: func(compatibility Compatibilitier, promote Entrypointer, config *configuration.Configuration) {
		// 		args := []string{}

		// 		options := &handler.Options{
		// 			DryRun:                       true,
		// 			EnableSemanticVersionTags:    false,
		// 			SourceImageName:              "myregistryhost.com/namespace/ubuntu:20.04",
		// 			TargetImageName:              "myubuntu",
		// 			TargetImageRegistryNamespace: "stable",
		// 			TargetImageRegistryHost:      "myprodregistryhost.com",
		// 			TargetImageTags:              []string{"tag1", "tag2"},
		// 			RemoveTargetImageTags:        true,
		// 			SemanticVersionTagsTemplates: []string{},
		// 			PromoteSourceImageTag:        false,
		// 			RemoteSourceImage:            false,
		// 		}

		// 		promote.(*entrypoint.MockEntrypoint).On("Execute", context.TODO(), args, config, options).Return(nil)
		// 	},
		// 	err: &errors.Error{},
		// },
		// {
		// 	desc:       "Testing to promote image and semver tags",
		// 	handler:    handler.NewHandlerMock(),
		// 	entrypoint: entrypoint.NewMockEntrypoint(),
		// 	args: []string{
		// 		"--dry-run",
		// 		"myregistryhost.com/namespace/ubuntu:1.2.3",
		// 		"--enable-semver-tags",
		// 		"--semver-tags-template",
		// 		"{{ .Major }}",
		// 		"--semver-tags-template",
		// 		"{{ .Major }}.{{ .Minor }}",
		// 		"--promote-source-tags",
		// 		"--remove-local-images-after-push",
		// 		"--remote-source-image",
		// 	},
		// 	prepareMockFunc: func(compatibility Compatibilitier, promote Entrypointer, config *configuration.Configuration) {
		// 		args := []string{}

		// 		options := &handler.Options{
		// 			DryRun:                       true,
		// 			EnableSemanticVersionTags:    true,
		// 			SourceImageName:              "myregistryhost.com/namespace/ubuntu:1.2.3",
		// 			TargetImageName:              "",
		// 			TargetImageRegistryNamespace: "",
		// 			TargetImageRegistryHost:      "",
		// 			TargetImageTags:              []string{},
		// 			RemoveTargetImageTags:        true,
		// 			SemanticVersionTagsTemplates: []string{"{{ .Major }}", "{{ .Major }}.{{ .Minor }}"},
		// 			PromoteSourceImageTag:        true,
		// 			RemoteSourceImage:            true,
		// 		}

		// 		promote.(*entrypoint.MockEntrypoint).On("Execute", context.TODO(), args, config, options).Return(nil)
		// 	},
		// 	err: &errors.Error{},
		// },
		// {
		// 	desc:       "Testing to promote without image name",
		// 	handler:    handler.NewHandlerMock(),
		// 	entrypoint: entrypoint.NewMockEntrypoint(),
		// 	args: []string{
		// 		"--dry-run",
		// 	},
		// 	err: errors.New("(promote::RunE)", "Source images name must be provided"),
		// },
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.compatibility, test.entrypoint, test.config)
			}

			cmd := NewCommand(context.TODO(), test.compatibility, test.config, test.entrypoint)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)

			if err != nil && assert.Error(t, err) {
				fmt.Println(err.Error())
				assert.Equal(t, test.err, err)
			} else {
				test.handler.(*handler.HandlerMock).AssertExpectations(t)
			}

		})
	}
}
