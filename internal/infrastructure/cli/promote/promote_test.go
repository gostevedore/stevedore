package promote

import (
	"context"
	"fmt"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/promote"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
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
				"--use-source-image-from-remote",
				"--use-docker-normalized-name",
			},
			prepareMockFunc: func(compatibility Compatibilitier, promote Entrypointer, config *configuration.Configuration) {

				entrypointOptions := &entrypoint.Options{
					UseDockerNormalizedName: true,
				}
				handlerOptions := &handler.Options{
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
					entrypointOptions,
					handlerOptions,
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
				"--remove-local-images-after-push",
				"--force-promote-source-image",
				"--use-source-image-from-remote",
			},
			prepareMockFunc: func(comp Compatibilitier, promote Entrypointer, config *configuration.Configuration) {

				entrypointOptions := &entrypoint.Options{
					UseDockerNormalizedName: false,
				}
				handlerOptions := &handler.Options{
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
					entrypointOptions,
					handlerOptions,
				).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageRemoveTargetImageTags}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:          "Testing to promote an image to a new registry host, registry namespace, with new name and multiple tags",
			handler:       handler.NewHandlerMock(),
			compatibility: compatibility.NewMockCompatibility(),
			config:        &configuration.Configuration{},
			entrypoint:    entrypoint.NewMockEntrypoint(),
			args: []string{
				"myregistryhost.com/namespace/ubuntu:20.04",
				"--dry-run",
				"--promote-image-name",
				"myubuntu",
				"--promote-image-registry-namespace",
				"stable",
				"--promote-image-registry-host",
				"myprodregistryhost.com",
				"--promote-image-tag",
				"tag1",
				"--promote-image-tag",
				"tag2",
				"--remove-local-images-after-push",
			},
			prepareMockFunc: func(compatibility Compatibilitier, promote Entrypointer, config *configuration.Configuration) {
				entrypointOptions := &entrypoint.Options{}
				handlerOptions := &handler.Options{
					DryRun:                       true,
					EnableSemanticVersionTags:    false,
					TargetImageName:              "myubuntu",
					TargetImageRegistryNamespace: "stable",
					TargetImageRegistryHost:      "myprodregistryhost.com",
					TargetImageTags:              []string{"tag1", "tag2"},
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{},
					PromoteSourceImageTag:        false,
					RemoteSourceImage:            false,
				}

				promote.(*entrypoint.MockEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{
						"myregistryhost.com/namespace/ubuntu:20.04",
					},
					config,
					entrypointOptions,
					handlerOptions,
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:       "Testing to promote image and semver tags",
			handler:    handler.NewHandlerMock(),
			entrypoint: entrypoint.NewMockEntrypoint(),
			args: []string{
				"myregistryhost.com/namespace/ubuntu:1.2.3",
				"--dry-run",
				"--enable-semver-tags",
				"--semver-tags-template",
				"{{ .Major }}",
				"--semver-tags-template",
				"{{ .Major }}.{{ .Minor }}",
				"--use-source-image-from-remote",
				"--remove-local-images-after-push",
			},
			prepareMockFunc: func(compatibility Compatibilitier, promote Entrypointer, config *configuration.Configuration) {

				entrypointOptions := &entrypoint.Options{}
				handlerOptions := &handler.Options{
					DryRun:                       true,
					EnableSemanticVersionTags:    true,
					TargetImageName:              image.UndefinedStringValue,
					TargetImageRegistryNamespace: image.UndefinedStringValue,
					TargetImageRegistryHost:      image.UndefinedStringValue,
					TargetImageTags:              []string{},
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}", "{{ .Major }}.{{ .Minor }}"},
					PromoteSourceImageTag:        false,
					RemoteSourceImage:            true,
				}

				promote.(*entrypoint.MockEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{
						"myregistryhost.com/namespace/ubuntu:1.2.3",
					},
					config,
					entrypointOptions,
					handlerOptions,
				).Return(nil)
			},
			err: &errors.Error{},
		},
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
