package ansibledriver

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/gostevedore/stevedore/internal/types"

	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/ui/console"

	ansibler "github.com/apenella/go-ansible"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestNewAnsiblePlaybookDriver(t *testing.T) {

	var w bytes.Buffer
	ctx := context.TODO()

	cons := &console.Console{
		Writer: io.Writer(&w),
	}

	tests := []struct {
		desc    string
		options *types.BuildOptions
		context context.Context
		err     error
		res     *ansibler.AnsiblePlaybookCmd
	}{
		{
			desc:    "Testing new ansiblePlaybookBuilder with nil options",
			options: nil,
			context: nil,
			err:     errors.New("(build::NewAnsiblePlaybookDriver)", "Build options are nil"),
			res:     nil,
		},
		// {
		// 	desc: "Testing new ansiblePlaybookBuilder with a nil context",
		// 	options: &types.BuildOptions{
		// 		BuilderOptions: map[string]interface{}{},
		// 	},
		// 	context: nil,
		// 	err:     errors.New("(build::NewAnsiblePlaybookDriver)", "Context is nil"),
		// 	res:     nil,
		// },
		{
			desc: "Testing options without a playbook defined",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{},
			},
			context: ctx,
			err:     errors.New("(build::NewAnsiblePlaybookDriver)", "playbook has not been defined on build options"),
			res:     nil,
		},
		{
			desc: "Testing an image with undefined image name",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
			},
			context: ctx,
			err:     errors.New("(build::NewAnsiblePlaybookDriver)", "Image name is not set"),
			res:     nil,
		},
		{
			desc: "Testing options without an inventory defined",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				BuilderOptions: map[string]interface{}{
					"playbook": "playbook",
				},
			},
			context: ctx,
			err:     errors.New("(build::NewAnsiblePlaybookDriver)", "inventory has not been defined on build options"),
			res:     nil,
		},
		{
			desc: "Testing an image with a registry defined",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"image_registry_host":      "registry",
						"image_builder_label":      "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
		{
			desc: "Testing an image with a main version defined",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				ImageVersion:      "version",
				RegistryNamespace: "namespace",
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName:version",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"image_tag":                "version",
						"image_builder_label":      "builder_namespace_imageName_version",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
		{
			desc: "Testing an image with a vars defined",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
				},
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"var1":                     "value1",
						"var2":                     "value2",
						"image_builder_label":      "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
		{
			desc: "Testing an image with persistent vars defined",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
					"pvar2": "pvalue2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
				},
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"pvar1":                    "pvalue1",
						"pvar2":                    "pvalue2",
						"var1":                     "value1",
						"var2":                     "value2",
						"image_builder_label":      "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
		{
			desc: "Testing an image with persistent vars defined avoiding an overwrite",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
				},
				Vars: map[string]interface{}{
					"pvar1": "newvalue1",
				},
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"pvar1":                    "pvalue1",
						"image_builder_label":      "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
		{
			desc: "Testing a build skipping image push",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
				},
				PushImages: false,
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"var1":                     "value1",
						"var2":                     "value2",
						"push_image":               false,
						"image_builder_label":      "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
		{
			desc: "Testing a build with ansible local connection",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
				},
				ConnectionLocal: true,
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":               "imageName",
						"image_registry_namespace": "namespace",
						"var1":                     "value1",
						"var2":                     "value2",
						"image_builder_label":      "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{
					Connection: "local",
				},
			},
		},
		{
			desc: "Testing a build image giving image from specs",
			options: &types.BuildOptions{
				ImageName:         "imageName",
				RegistryNamespace: "namespace",
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
				},
				BuilderOptions: map[string]interface{}{
					"playbook":  "playbook",
					"inventory": "inventory",
				},
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
				ImageFromName:              "parent",
				ImageFromRegistryNamespace: "parentNamespace",
				ImageFromRegistryHost:      "parentRegistry",
				ImageFromVersion:           "parentVersion",
				PushImages:                 true,
			},
			context: ctx,
			err:     nil,
			res: &ansibler.AnsiblePlaybookCmd{
				Playbook:   "playbook",
				Writer:     cons,
				ExecPrefix: "imageName",
				Options: &ansibler.AnsiblePlaybookOptions{
					Inventory: "inventory",
					ExtraVars: map[string]interface{}{
						"image_name":                    "imageName",
						"image_registry_namespace":      "namespace",
						"var1":                          "value1",
						"var2":                          "value2",
						"image_from_name":               "parent",
						"image_from_registry_namespace": "parentNamespace",
						"image_from_registry_host":      "parentRegistry",
						"image_from_tag":                "parentVersion",
						"image_builder_label":           "builder_namespace_imageName",
					},
				},
				ConnectionOptions: &ansibler.AnsiblePlaybookConnectionOptions{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			builderer, err := NewAnsiblePlaybookDriver(test.context, test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res.Playbook, builderer.(*ansibler.AnsiblePlaybookCmd).Playbook, "Unexpected Playbook")
				assert.Equal(t, test.res.ExecPrefix, builderer.(*ansibler.AnsiblePlaybookCmd).ExecPrefix, "Unexpected ExecPrefix")
				assert.Equal(t, test.res.Options, builderer.(*ansibler.AnsiblePlaybookCmd).Options, "Unexpected Options")
				assert.Equal(t, test.res.ConnectionOptions, builderer.(*ansibler.AnsiblePlaybookCmd).ConnectionOptions, "Unexpected ConnectionOptions")
				assert.Equal(t, test.res.PrivilegeEscalationOptions, builderer.(*ansibler.AnsiblePlaybookCmd).PrivilegeEscalationOptions, "Unexpected PrivilegeEscalationOptions")
				assert.Equal(t, test.res.StdoutCallback, builderer.(*ansibler.AnsiblePlaybookCmd).StdoutCallback, "Unexpected StdoutCallback")
			}
		})
	}

}
