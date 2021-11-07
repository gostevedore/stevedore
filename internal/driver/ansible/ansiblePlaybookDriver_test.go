package ansibledriver

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/driver/ansible/goansible"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestTestNewAnsiblePlaybookDriver(t *testing.T) {

	errContext := "(ansibledriver::NewAnsiblePlaybookDriver)"

	tests := []struct {
		desc   string
		driver Ansibler
		writer io.Writer
		res    *AnsiblePlaybookDriver
		err    error
	}{
		{
			desc:   "Testing error creating an ansible-playbook driver with nil driver",
			driver: nil,
			writer: nil,
			err:    errors.New(errContext, "To create an AnsiblePlaybookDriver is expected a driver"),
		},
		{
			desc:   "Testing create and ansible-playbook driver",
			driver: goansible.NewMockAnsibleDriver(),
			writer: nil,
			res: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: os.Stdout,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := NewAnsiblePlaybookDriver(test.driver, test.writer)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}

		})
	}

}

func TestBuild(t *testing.T) {

	// 	var w bytes.Buffer
	// 	ctx := context.TODO()

	// 	cons := &console.Console{
	// 		Writer: io.Writer(&w),
	// 	}

	errContext := "(ansibledriver::Build)"

	tests := []struct {
		desc              string
		driver            *AnsiblePlaybookDriver
		options           *types.BuildOptions
		err               error
		prepareAssertFunc func(driver Ansibler)
		assertFunc        func(driver Ansibler) bool
	}{
		{
			desc: "Testing error building an image build with nil driver",
			driver: &AnsiblePlaybookDriver{
				driver: nil,
			},
			options: nil,
			err:     errors.New(errContext, "Build driver is missing"),
		},
		{
			desc: "Testing error building an image with nil options",
			driver: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: nil,
			},
			options: nil,
			err:     errors.New(errContext, "Build options are nil"),
		},
		{
			desc: "Testing error building without a playbook defined on builder options",
			driver: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: nil,
			},
			options: &types.BuildOptions{},
			err:     errors.New(errContext, "Playbook has not been defined on build options"),
		},
		{
			desc: "Testing error building an image with undefined image name",
			driver: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: nil,
			},
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"playbook": "site.yml",
				},
			},
			err: errors.New(errContext, "Inventory has not been defined on build options"),
		},
		{
			desc: "Testing error building an image with undefined inventory",
			driver: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: nil,
			},
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"playbook":  "site.yml",
					"inventory": "inventory.yml",
				},
			},
			err: errors.New(errContext, "Image has not been defined on build options"),
		},
		{
			desc: "Testing build an image",
			driver: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: os.Stdout,
			},
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"playbook":  "site.yml",
					"inventory": "inventory.yml",
				},
				ImageName:         "image_name",
				ImageVersion:      "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
				ConnectionLocal:   true,
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
			prepareAssertFunc: func(driver Ansibler) {

				ansibleOptions := &ansible.AnsiblePlaybookOptions{
					Inventory: "inventory.yml",
					ExtraVars: map[string]interface{}{
						"image_builder_label":      "builder_namespace_image_name_version",
						"image_name":               "image_name",
						"image_registry_host":      "registry",
						"image_registry_namespace": "namespace",
						"image_tag":                "version",
						"push_image":               false,
					},
				}
				ansibleConnectionOptions := &options.AnsibleConnectionOptions{
					Connection: "local",
				}

				driver.(*goansible.MockAnsibleDriver).On("WithPlaybook", "site.yml")
				driver.(*goansible.MockAnsibleDriver).On("WithOptions", ansibleOptions)
				driver.(*goansible.MockAnsibleDriver).On("WithConnectionOptions", ansibleConnectionOptions)
				driver.(*goansible.MockAnsibleDriver).On("PrepareExecutor", os.Stdout, "image_name:version")
				driver.(*goansible.MockAnsibleDriver).On("Run", context.TODO()).Return(nil)
			},
			assertFunc: func(driver Ansibler) bool {
				return driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithPlaybook", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithConnectionOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "PrepareExecutor", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "Run", 1)
			},
		},
		{
			desc: "Testing build an image with all build options",
			driver: &AnsiblePlaybookDriver{
				driver: goansible.NewMockAnsibleDriver(),
				writer: os.Stdout,
			},
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"playbook":  "site.yml",
					"inventory": "inventory.yml",
				},
				BuilderName:                "builder",
				ImageName:                  "image_name",
				ImageVersion:               "version",
				RegistryNamespace:          "namespace",
				RegistryHost:               "registry",
				ImageFromName:              "from_image",
				ImageFromVersion:           "from_version",
				ImageFromRegistryNamespace: "from_namespace",
				ImageFromRegistryHost:      "from_registry",
				OutputPrefix:               "prefix",
				ConnectionLocal:            true,
				Tags: []string{
					"tag1",
					"tag2",
				},
				PersistentVars: map[string]interface{}{
					"persistent_var1": "value1",
					"persistent_var2": "value2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
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
					varsmap.VarMappingImageExtraTagsKey:                varsmap.VarMappingImageExtraTagsDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
				},
			},
			prepareAssertFunc: func(driver Ansibler) {

				ansibleOptions := &ansible.AnsiblePlaybookOptions{
					Inventory: "inventory.yml",
					ExtraVars: map[string]interface{}{
						"image_builder_label":           "builder",
						"image_name":                    "image_name",
						"image_registry_host":           "registry",
						"image_registry_namespace":      "namespace",
						"image_tag":                     "version",
						"image_from_name":               "from_image",
						"image_from_tag":                "from_version",
						"image_from_registry_host":      "from_registry",
						"image_from_registry_namespace": "from_namespace",
						"push_image":                    false,
						"persistent_var1":               "value1",
						"persistent_var2":               "value2",
						"var1":                          "value1",
						"var2":                          "value2",
						"image_extra_tags":              []string{"tag1", "tag2"},
					},
				}
				ansibleConnectionOptions := &options.AnsibleConnectionOptions{
					Connection: "local",
				}

				driver.(*goansible.MockAnsibleDriver).On("WithPlaybook", "site.yml")
				driver.(*goansible.MockAnsibleDriver).On("WithOptions", ansibleOptions)
				driver.(*goansible.MockAnsibleDriver).On("WithConnectionOptions", ansibleConnectionOptions)
				driver.(*goansible.MockAnsibleDriver).On("PrepareExecutor", os.Stdout, "prefix")
				driver.(*goansible.MockAnsibleDriver).On("Run", context.TODO()).Return(nil)
			},
			assertFunc: func(driver Ansibler) bool {
				return driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithPlaybook", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithConnectionOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "PrepareExecutor", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "Run", 1)
			},
		},

		// 		{
		// 			desc: "Testing an image with a registry defined",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				RegistryHost:      "registry",
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				PushImages: true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"image_registry_host":      "registry",
		// 						"image_builder_label":      "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing an image with a main version defined",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				ImageVersion:      "version",
		// 				RegistryNamespace: "namespace",
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				PushImages: true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"image_tag":                "version",
		// 						"image_builder_label":      "builder_namespace_imageName_version",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing an image with a vars defined",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				Vars: map[string]interface{}{
		// 					"var1": "value1",
		// 					"var2": "value2",
		// 				},
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				PushImages: true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"var1":                     "value1",
		// 						"var2":                     "value2",
		// 						"image_builder_label":      "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing an image with persistent vars defined",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				PersistentVars: map[string]interface{}{
		// 					"pvar1": "pvalue1",
		// 					"pvar2": "pvalue2",
		// 				},
		// 				Vars: map[string]interface{}{
		// 					"var1": "value1",
		// 					"var2": "value2",
		// 				},
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				PushImages: true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"pvar1":                    "pvalue1",
		// 						"pvar2":                    "pvalue2",
		// 						"var1":                     "value1",
		// 						"var2":                     "value2",
		// 						"image_builder_label":      "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing an image with persistent vars defined avoiding an overwrite",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				PersistentVars: map[string]interface{}{
		// 					"pvar1": "pvalue1",
		// 				},
		// 				Vars: map[string]interface{}{
		// 					"pvar1": "newvalue1",
		// 				},
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				PushImages: true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				), Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"pvar1":                    "pvalue1",
		// 						"image_builder_label":      "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing a build skipping image push",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				Vars: map[string]interface{}{
		// 					"var1": "value1",
		// 					"var2": "value2",
		// 				},
		// 				PushImages: false,
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"var1":                     "value1",
		// 						"var2":                     "value2",
		// 						"push_image":               false,
		// 						"image_builder_label":      "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing a build with ansible local connection",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				Vars: map[string]interface{}{
		// 					"var1": "value1",
		// 					"var2": "value2",
		// 				},
		// 				ConnectionLocal: true,
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				PushImages: true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":               "imageName",
		// 						"image_registry_namespace": "namespace",
		// 						"var1":                     "value1",
		// 						"var2":                     "value2",
		// 						"image_builder_label":      "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{
		// 					Connection: "local",
		// 				},
		// 			},
		// 		},
		// 		{
		// 			desc: "Testing a build image giving image from specs",
		// 			options: &types.BuildOptions{
		// 				ImageName:         "imageName",
		// 				RegistryNamespace: "namespace",
		// 				Vars: map[string]interface{}{
		// 					"var1": "value1",
		// 					"var2": "value2",
		// 				},
		// 				BuilderOptions: map[string]interface{}{
		// 					"playbook":  "playbook",
		// 					"inventory": "inventory",
		// 				},
		// 				BuilderVarMappings: map[string]string{
		// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
		// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
		// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
		// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
		// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
		// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
		// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
		// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
		// 					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
		// 				},
		// 				ImageFromName:              "parent",
		// 				ImageFromRegistryNamespace: "parentNamespace",
		// 				ImageFromRegistryHost:      "parentRegistry",
		// 				ImageFromVersion:           "parentVersion",
		// 				PushImages:                 true,
		// 			},
		// 			context: ctx,
		// 			err:     nil,
		// 			res: &ansible.AnsiblePlaybookCmd{
		// 				Playbooks: []string{"playbook"},
		// 				Exec: execute.NewDefaultExecute(
		// 					execute.WithWrite(cons),
		// 					execute.WithTransformers(
		// 						results.Prepend("imageName"),
		// 					),
		// 				),
		// 				Options: &ansible.AnsiblePlaybookOptions{
		// 					Inventory: "inventory",
		// 					ExtraVars: map[string]interface{}{
		// 						"image_name":                    "imageName",
		// 						"image_registry_namespace":      "namespace",
		// 						"var1":                          "value1",
		// 						"var2":                          "value2",
		// 						"image_from_name":               "parent",
		// 						"image_from_registry_namespace": "parentNamespace",
		// 						"image_from_registry_host":      "parentRegistry",
		// 						"image_from_tag":                "parentVersion",
		// 						"image_builder_label":           "builder_namespace_imageName",
		// 					},
		// 				},
		// 				ConnectionOptions: &options.AnsibleConnectionOptions{},
		// 			},
		// 		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.driver.driver)
			}

			err := test.driver.Build(context.TODO(), test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.True(t, test.assertFunc(test.driver.driver))
			}

			// 			builderer, err := NewAnsiblePlaybookDriver(test.context, test.options)
			// 			if err != nil && assert.Error(t, err) {
			// 				assert.Equal(t, test.err, err)
			// 			} else {
			// 				assert.Equal(t, test.res.Playbooks, builderer.(*ansible.AnsiblePlaybookCmd).Playbooks, "Unexpected Playbook")
			// 				assert.Equal(t, test.res.Options, builderer.(*ansible.AnsiblePlaybookCmd).Options, "Unexpected Options")
			// 				assert.Equal(t, test.res.ConnectionOptions, builderer.(*ansible.AnsiblePlaybookCmd).ConnectionOptions, "Unexpected ConnectionOptions")
			// 				assert.Equal(t, test.res.PrivilegeEscalationOptions, builderer.(*ansible.AnsiblePlaybookCmd).PrivilegeEscalationOptions, "Unexpected PrivilegeEscalationOptions")
			// 				assert.Equal(t, test.res.StdoutCallback, builderer.(*ansible.AnsiblePlaybookCmd).StdoutCallback, "Unexpected StdoutCallback")
			// 			}
		})
	}

}
