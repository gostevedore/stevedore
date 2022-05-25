package build

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/build"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {

	tests := []struct {
		desc              string
		compatibility     Compatibilitier
		config            *configuration.Configuration
		entrypoint        Entrypointer
		args              []string
		err               error
		prepareAssertFunc func(Compatibilitier, Entrypointer, *configuration.Configuration)
	}{
		{
			desc:          "Testing run build command",
			config:        &configuration.Configuration{},
			compatibility: compatibility.NewMockCompatibility(),
			entrypoint:    entrypoint.NewMockEntrypoint(),
			args: []string{
				"my-image",
				"--ansible-connection-local",
				"--ansible-intermediate-container-name",
				"container",
				"--ansible-inventory-path",
				"inventory",
				"--ansible-limit",
				"limit",
				"--enable-semver-tags",
				"--image-from-name",
				"image-from-name",
				"--image-from-namespace",
				"image-from-namespace",
				"--image-from-registry",
				"image-from-registry",
				"--image-from-version",
				"image-from-version",
				"--image-name",
				"image-name",
				"--image-registry-host",
				"image-registry-host",
				"--image-registry-namespace",
				"image-registry-namespace",
				"--image-version",
				"image-version",
				"--persistent-variable",
				"pvar=pvalue",
				"--variable",
				"var=value",
				"--tag",
				"tag",
				"--label",
				"name=value",
				"--semver-tags-template",
				"{{ .Major }}",
				"--build-on-cascade",
				"--cascade-depth",
				"3",
				"--concurrency",
				"5",
				"--dry-run",
				"--pull-parent-image",
				"--push-after-build",
				"--remove-local-images-after-push",
			},
			prepareAssertFunc: func(compatibility Compatibilitier, build Entrypointer, config *configuration.Configuration) {
				build.(*entrypoint.MockEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{"my-image"},
					config,
					compatibility,
					&entrypoint.Options{
						Concurrency: 5,
					},
					&handler.Options{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "container",
						AnsibleInventoryPath:             "inventory",
						AnsibleLimit:                     "limit",
						BuildOnCascade:                   true,
						CascadeDepth:                     3,
						DryRun:                           true,
						EnableSemanticVersionTags:        true,
						ImageFromName:                    "image-from-name",
						ImageFromRegistryHost:            "image-from-registry",
						ImageFromRegistryNamespace:       "image-from-namespace",
						ImageFromVersion:                 "image-from-version",
						ImageName:                        "image-name",
						ImageRegistryHost:                "image-registry-host",
						ImageRegistryNamespace:           "image-registry-namespace",
						Labels:                           []string{"name=value"},
						PersistentVars: []string{
							"pvar=pvalue",
						},
						PullParentImage:       true,
						PushImagesAfterBuild:  true,
						RemoveImagesAfterPush: true,
						SemanticVersionTagsTemplates: []string{
							"{{ .Major }}",
						},
						Tags: []string{"tag"},
						Vars: []string{
							"var=value",
						},
						Versions: []string{"image-version"},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:          "Testing run build command with deprecated commands",
			config:        &configuration.Configuration{},
			compatibility: compatibility.NewMockCompatibility(),
			entrypoint:    entrypoint.NewMockEntrypoint(),
			args: []string{
				"my-image",
				"--connection-local",
				"--builder-name",
				"container",
				"--inventory",
				"inventory",
				"--limit",
				"limit",
				"--enable-semver-tags",
				"--image-from",
				"image-from-name",
				"--image-from-namespace",
				"image-from-namespace",
				"--image-from-registry",
				"image-from-registry",
				"--image-from-version",
				"image-from-version",
				"--image-name",
				"image-name",
				"--registry",
				"image-registry-host",
				"--namespace",
				"image-registry-namespace",
				"--image-version",
				"image-version",
				"--set-persistent",
				"pvar=pvalue",
				"--set",
				"var=value",
				"--tag",
				"tag",
				"--label",
				"name=value",
				"--semver-tags-template",
				"{{ .Major }}",
				"--cascade",
				"--cascade-depth",
				"3",
				"--num-workers",
				"5",
				"--dry-run",
				"--pull-parent-image",
				"--push-after-build",
				"--remove-local-images-after-push",
				"--no-push",
			},
			prepareAssertFunc: func(comp Compatibilitier, build Entrypointer, config *configuration.Configuration) {

				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageConnectionLocal}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageBuildBuilderName}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageInventory}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageLimit}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageImageFrom}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageRegistry}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageNamespace}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageSetPersistent}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageSet}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageCascade}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageNumWorkers}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessagePushImages}).Return(nil)

				build.(*entrypoint.MockEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{"my-image"},
					config,
					comp,
					&entrypoint.Options{
						Concurrency: 5,
					},
					&handler.Options{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "container",
						AnsibleInventoryPath:             "inventory",
						AnsibleLimit:                     "limit",
						BuildOnCascade:                   true,
						CascadeDepth:                     3,
						DryRun:                           true,
						EnableSemanticVersionTags:        true,
						ImageFromName:                    "image-from-name",
						ImageFromRegistryHost:            "image-from-registry",
						ImageFromRegistryNamespace:       "image-from-namespace",
						ImageFromVersion:                 "image-from-version",
						ImageName:                        "image-name",
						ImageRegistryHost:                "image-registry-host",
						ImageRegistryNamespace:           "image-registry-namespace",
						Labels:                           []string{"name=value"},
						PersistentVars: []string{
							"pvar=pvalue",
						},
						PullParentImage:       true,
						PushImagesAfterBuild:  true,
						RemoveImagesAfterPush: true,
						SemanticVersionTagsTemplates: []string{
							"{{ .Major }}",
						},
						Tags: []string{"tag"},
						Vars: []string{
							"var=value",
						},
						Versions: []string{"image-version"},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.compatibility, test.entrypoint, test.config)
			}

			cmd := NewCommand(context.TODO(), test.compatibility, test.config, test.entrypoint)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}
