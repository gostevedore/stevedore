package handler

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/service/promote"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	errContext := "(promote::Handler)"

	tests := []struct {
		desc            string
		handler         *Handler
		options         *Options
		cmd             *cobra.Command
		cmdArgs         []string
		prepareMockFunc func(*Handler)
		err             error
	}{
		{
			desc: "Testing promote handler error when no source image is provided",
			handler: &Handler{
				service: promote.NewMockService(),
			},
			cmd:     &cobra.Command{},
			cmdArgs: []string{},
			err:     errors.New(errContext, "Source images name must be provided"),
			options: &Options{

				SourceImageName:              "",
				TargetImageName:              "target_name",
				TargetImageRegistryNamespace: "target_registry_namespace",
				TargetImageRegistryHost:      "target_registry_host",
			},
		},
		{
			desc: "Testing promote handler passing all options",
			handler: &Handler{
				service: promote.NewMockService(),
			},
			cmd: &cobra.Command{},
			cmdArgs: []string{
				"source_name",
			},
			err: &errors.Error{},
			options: &Options{
				DryRun:                       true,
				EnableSemanticVersionTags:    true,
				SourceImageName:              "source_name",
				TargetImageName:              "target_name",
				TargetImageRegistryNamespace: "target_registry_namespace",
				TargetImageRegistryHost:      "target_registry_host",
				TargetImageTags:              []string{"tag"},
				RemoveTargetImageTags:        true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
				PromoteSourceImageTag:        true,
				RemoteSourceImage:            true,
			},
			prepareMockFunc: func(h *Handler) {

				options := &promote.ServiceOptions{
					EnableSemanticVersionTags:    true,
					TargetImageName:              "target_name",
					TargetImageRegistryNamespace: "target_registry_namespace",
					TargetImageRegistryHost:      "target_registry_host",
					TargetImageTags:              []string{"tag"},
					PromoteSourceImageTag:        true,
					RemoveTargetImageTags:        true,
					RemoteSourceImage:            true,
					SourceImageName:              "source_name",
					SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
				}

				h.service.(*promote.MockService).On("Promote", context.TODO(), options, "dry-run").Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.handler)
			}

			err := test.handler.Handler(context.TODO(), test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.handler.service.(*promote.MockService).AssertExpectations(t)
			}

		})
	}
}
