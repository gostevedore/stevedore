package promote

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/service/promote"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	errContext := "(promote::Handler)"

	tests := []struct {
		desc    string
		handler *Handler
		options *Options

		prepareAssertFunc func(*Handler)
		err               error
	}{
		{
			desc:    "Testing promote handler error when no source image is provided",
			handler: NewHandler(promote.NewMockService()),
			err:     errors.New(errContext, "Source images name must be provided"),
			options: &Options{

				SourceImageName:              "",
				TargetImageName:              "target_name",
				TargetImageRegistryNamespace: "target_registry_namespace",
				TargetImageRegistryHost:      "target_registry_host",
			},
		},
		{
			desc:    "Testing promote handler passing all options",
			handler: NewHandler(promote.NewMockService()),
			err:     &errors.Error{},
			options: &Options{
				DryRun:                       true,
				EnableSemanticVersionTags:    true,
				PromoteSourceImageTag:        true,
				RemoteSourceImage:            true,
				RemoveTargetImageTags:        true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
				SourceImageName:              "source_name",
				TargetImageName:              "target_name",
				TargetImageRegistryHost:      "target_registry_host",
				TargetImageRegistryNamespace: "target_registry_namespace",
				TargetImageTags:              []string{"tag"},
			},
			prepareAssertFunc: func(h *Handler) {

				options := &promote.ServiceOptions{
					DryRun:                       true,
					EnableSemanticVersionTags:    true,
					PromoteSourceImageTag:        true,
					RemoteSourceImage:            true,
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
					SourceImageName:              "source_name",
					TargetImageName:              "target_name",
					TargetImageRegistryHost:      "target_registry_host",
					TargetImageRegistryNamespace: "target_registry_namespace",
					TargetImageTags:              []string{"tag"},
				}

				h.service.(*promote.MockService).On("Promote", context.TODO(), options).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.handler)
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
