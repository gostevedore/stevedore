package ansible

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/ansible/goansible"
	reference "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/default"
	"github.com/stretchr/testify/assert"
)

func TestNewAnsiblePlaybookDriver(t *testing.T) {

	errContext := "(ansibledriver::NewAnsiblePlaybookDriver)"

	tests := []struct {
		desc          string
		driver        AnsibleDriverer
		writer        io.Writer
		referenceName repository.ImageReferenceNamer
		res           *AnsiblePlaybookDriver
		err           error
	}{
		{
			desc:          "Testing error creating an ansible-playbook driver with nil driver",
			driver:        nil,
			referenceName: nil,
			writer:        nil,
			err:           errors.New(errContext, "To create an AnsiblePlaybookDriver is required a driver"),
		},
		{
			desc:          "Testing error creating an ansible-playbook driver with nil reference name",
			driver:        goansible.NewMockAnsibleDriver(),
			referenceName: nil,
			writer:        nil,
			err:           errors.New(errContext, "To create an AnsiblePlaybookDriver is required a reference name"),
		},
		{
			desc:          "Testing create and ansible-playbook driver",
			driver:        goansible.NewMockAnsibleDriver(),
			writer:        nil,
			referenceName: reference.NewDefaultReferenceName(),
			res: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        os.Stdout,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := NewAnsiblePlaybookDriver(test.driver, test.referenceName, test.writer)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}

		})
	}

}

func TestBuild(t *testing.T) {

	errContext := "(ansibledriver::Build)"

	tests := []struct {
		desc              string
		driver            *AnsiblePlaybookDriver
		ref               repository.ImageReferenceNamer
		image             *image.Image
		options           *image.BuildDriverOptions
		err               error
		prepareAssertFunc func(driver AnsibleDriverer)
		assertFunc        func(driver AnsibleDriverer) bool
	}{
		{
			desc: "Testing error building an image build with nil driver",
			driver: &AnsiblePlaybookDriver{
				driver: nil,
			},
			options: nil,
			err:     errors.New(errContext, "To build an image is required a driver"),
		},
		{
			desc: "Testing error building an image with nil reference name",
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: nil,
				writer:        nil,
			},
			options: nil,
			err:     errors.New(errContext, "To build an image is required a reference name"),
		},
		{
			desc: "Testing error building an image with nil image",
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        nil,
			},
			options: nil,
			err:     errors.New(errContext, "To build an image is required a image"),
		},
		{
			desc:  "Testing error building an image with nil options",
			image: &image.Image{},
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        nil,
			},
			options: nil,
			err:     errors.New(errContext, "To build an image is required a build options"),
		},
		{
			desc:  "Testing error building without options from the builder",
			image: &image.Image{},
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        nil,
			},
			options: &image.BuildDriverOptions{},
			err:     errors.New(errContext, "To build an image are required the options from the builder"),
		},
		{
			desc:  "Testing error building without a playbook defined on builder options",
			image: &image.Image{},
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        nil,
			},
			options: &image.BuildDriverOptions{
				BuilderOptions: &builder.BuilderOptions{},
			},
			err: errors.New(errContext, "Playbook has not been defined on build options"),
		},
		{
			desc:  "Testing error building an image with undefined image name",
			image: &image.Image{},
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        nil,
			},
			options: &image.BuildDriverOptions{
				BuilderOptions: &builder.BuilderOptions{
					Playbook: "site.yml",
				},
			},
			err: errors.New(errContext, "Inventory has not been defined on build options"),
		},
		{
			desc:  "Testing error building an image with undefined image name",
			image: &image.Image{},
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        nil,
			},
			options: &image.BuildDriverOptions{
				BuilderOptions: &builder.BuilderOptions{
					Playbook:  "site.yml",
					Inventory: "inventory.yml",
				},
			},
			err: errors.New(errContext, "Image name is not defined"),
		},
		{
			desc: "Testing build an image without parent",
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        os.Stdout,
			},
			image: &image.Image{
				Name:              "image_name",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
			},
			options: &image.BuildDriverOptions{
				BuilderOptions: &builder.BuilderOptions{
					Playbook:  "site.yml",
					Inventory: "inventory.yml",
				},
				AnsibleConnectionLocal: true,
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageFromFullyQualifiedNameKey:   varsmap.VarMappingImageFromFullyQualifiedNameValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFullyQualifiedNameKey:       varsmap.VarMappingImageFullyQualifiedNameValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
				},
			},
			prepareAssertFunc: func(driver AnsibleDriverer) {

				ansibleOptions := &ansible.AnsiblePlaybookOptions{
					Inventory: "inventory.yml",
					ExtraVars: map[string]interface{}{
						"image_builder_label":        "builder_namespace_image_name_version",
						"image_fully_qualified_name": "registry/namespace/image_name:version",
						"image_name":                 "image_name",
						"image_registry_host":        "registry",
						"image_registry_namespace":   "namespace",
						"image_tag":                  "version",
						"push_image":                 false,
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
			assertFunc: func(driver AnsibleDriverer) bool {
				return driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithPlaybook", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithConnectionOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "PrepareExecutor", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "Run", 1)
			},
		},
		{
			desc: "Testing build an image with parent and all the build options defined",
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        os.Stdout,
			},
			image: &image.Image{
				Name:              "image_name",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
				Parent: &image.Image{
					Name:              "from_image",
					Version:           "from_version",
					RegistryNamespace: "from_namespace",
					RegistryHost:      "from_registry",
				},
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
			},
			options: &image.BuildDriverOptions{
				BuilderOptions: &builder.BuilderOptions{
					Playbook:  "site.yml",
					Inventory: "inventory.yml",
				},
				OutputPrefix:                     "prefix",
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "intermediate_container",
				AnsibleInventoryPath:             "override-inventory.yml",
				AnsibleLimit:                     "limit",
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageExtraTagsKey:                varsmap.VarMappingImageExtraTagsDefaultValue,
					varsmap.VarMappingImageFromFullyQualifiedNameKey:   varsmap.VarMappingImageFromFullyQualifiedNameValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFullyQualifiedNameKey:       varsmap.VarMappingImageFullyQualifiedNameValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
				},
			},
			prepareAssertFunc: func(driver AnsibleDriverer) {
				ansibleOptions := &ansible.AnsiblePlaybookOptions{
					Inventory: "override-inventory.yml",
					Limit:     "limit",
					ExtraVars: map[string]interface{}{
						"image_builder_label":             "intermediate_container",
						"image_extra_tags":                []string{"tag1", "tag2"},
						"image_from_fully_qualified_name": "from_registry/from_namespace/from_image:from_version",
						"image_from_name":                 "from_image",
						"image_from_registry_host":        "from_registry",
						"image_from_registry_namespace":   "from_namespace",
						"image_from_tag":                  "from_version",
						"image_fully_qualified_name":      "registry/namespace/image_name:version",
						"image_name":                      "image_name",
						"image_registry_host":             "registry",
						"image_registry_namespace":        "namespace",
						"image_tag":                       "version",
						"persistent_var1":                 "value1",
						"persistent_var2":                 "value2",
						"push_image":                      false,
						"var1":                            "value1",
						"var2":                            "value2",
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
			assertFunc: func(driver AnsibleDriverer) bool {
				return driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithPlaybook", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithConnectionOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "PrepareExecutor", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "Run", 1)
			},
		},
		{
			desc: "Testing build an image with same variable defined either on persistent_vars and vars",
			driver: &AnsiblePlaybookDriver{
				driver:        goansible.NewMockAnsibleDriver(),
				referenceName: reference.NewDefaultReferenceName(),
				writer:        os.Stdout,
			},
			image: &image.Image{
				Name:              "image_name",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
				PersistentVars: map[string]interface{}{
					"var1": "persistent_value1",
					"var2": "persistent_value2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
					"var2": "value2",
				},
			},
			options: &image.BuildDriverOptions{
				BuilderOptions: &builder.BuilderOptions{
					Playbook:  "site.yml",
					Inventory: "inventory.yml",
				},
				AnsibleConnectionLocal: true,
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageFromFullyQualifiedNameKey:   varsmap.VarMappingImageFromFullyQualifiedNameValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFullyQualifiedNameKey:       varsmap.VarMappingImageFullyQualifiedNameValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
				},
			},
			prepareAssertFunc: func(driver AnsibleDriverer) {

				ansibleOptions := &ansible.AnsiblePlaybookOptions{
					Inventory: "inventory.yml",
					ExtraVars: map[string]interface{}{
						"image_builder_label":        "builder_namespace_image_name_version",
						"image_fully_qualified_name": "registry/namespace/image_name:version",
						"image_name":                 "image_name",
						"image_registry_host":        "registry",
						"image_registry_namespace":   "namespace",
						"image_tag":                  "version",
						"push_image":                 false,
						"var1":                       "persistent_value1",
						"var2":                       "persistent_value2",
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
			assertFunc: func(driver AnsibleDriverer) bool {
				return driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithPlaybook", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "WithConnectionOptions", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "PrepareExecutor", 1) &&
					driver.(*goansible.MockAnsibleDriver).AssertNumberOfCalls(t, "Run", 1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.driver.driver)
			}

			err := test.driver.Build(context.TODO(), test.image, test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.True(t, test.assertFunc(test.driver.driver))
			}
		})
	}

}
