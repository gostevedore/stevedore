package docker

import (
	"context"
	"io"
	"os"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/docker/godockerbuilder"
	"github.com/stretchr/testify/assert"
)

func TestNewDockerDriver(t *testing.T) {
	errContext := "(dockerdriver::NewDockerDriver)"

	tests := []struct {
		desc   string
		driver DockerDriverer
		writer io.Writer
		res    *DockerDriver
		err    error
	}{
		{
			desc:   "Testing error creating a docker driver with nil driver",
			driver: nil,
			writer: nil,
			err:    errors.New(errContext, "To create a DockerDriver is expected a driver"),
		},
		{
			desc:   "Testing create a docker driver",
			driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			writer: nil,
			err:    &errors.Error{},
			res: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
				writer: os.Stdout,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := NewDockerDriver(test.driver, test.writer)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}

		})
	}

}

func TestBuild(t *testing.T) {

	errContext := "(dockerdriver::Build)"

	tests := []struct {
		desc              string
		driver            *DockerDriver
		ctx               context.Context
		image             *image.Image
		options           *image.BuildDriverOptions
		prepareAssertFunc func(DockerDriverer)
		assertFunc        func(*testing.T, DockerDriverer)
		err               error
	}{
		{
			desc: "Testing error building a docker image with nil driver",
			driver: &DockerDriver{
				driver: nil,
			},
			err: errors.New(errContext, "To build an image is required a driver"),
		},
		{
			desc: "Testing error building a docker image with nil image",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			},
			err: errors.New(errContext, "To build an image is required a image"),
		},
		{
			desc:  "Testing error building a docker image with nil options",
			image: &image.Image{},
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			},
			err: errors.New(errContext, "To build an image is required a build options"),
		},
		{
			desc:  "Testing error building a docker image with nil golang context",
			image: &image.Image{},
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			},
			options: &image.BuildDriverOptions{},
			err:     errors.New(errContext, "To build an image is required a golang context"),
		},
		{
			desc:  "Testing error building a docker image with not defined image name",
			image: &image.Image{},
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			},
			ctx:     context.TODO(),
			options: &image.BuildDriverOptions{},
			err:     errors.New(errContext, "To build an image is required an image name"),
		},
		{
			desc: "Testing building a docker image",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
				writer: os.Stdout,
			},
			ctx: context.TODO(),
			image: &image.Image{
				Name:              "image",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "myregistry.test",
				Parent: &image.Image{
					Name:              "image-from-name",
					Version:           "image-from-version",
					RegistryNamespace: "image-from-registry-namespace",
					RegistryHost:      "image-from-registry-host.test",
				},
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
					"pvar2": "pvalue2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
				},
				Tags:   []string{"tag1", "tag2"},
				Labels: map[string]string{"label1": "value1", "label2": "value2"},
			},
			options: &image.BuildDriverOptions{
				PushImageAfterBuild:   true,
				PullParentImage:       true,
				RemoveImageAfterBuild: true,
				OutputPrefix:          "output-prefix",
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
				},
				PullAuthUsername: "pull-user",
				PullAuthPassword: "pull-pass",
				PushAuthUsername: "push-user",
				PushAuthPassword: "push-pass",
				BuilderOptions: &builder.BuilderOptions{
					Dockerfile: "Dockerfile.test",
					Context: []*builder.DockerDriverContextOptions{
						{Path: "/path/to/file"},
						{Path: "/path/to/file2"},
						{Git: &builder.DockerDriverGitContextOptions{
							Repository: "repo",
							Reference:  "main",
							Path:       "path",
							Auth: &builder.DockerDriverGitContextAuthOptions{
								Username: "user",
								Password: "pass",
							},
						}},
					},
				},
			},
			prepareAssertFunc: func(driver DockerDriverer) {
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithImageName", "myregistry.test/namespace/image:version")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithDockerfile", "Dockerfile.test")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "pvar1", "pvalue1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "pvar2", "pvalue2").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "var1", "value1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag1"}).Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag2"}).Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddLabel", "label1", "value1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddLabel", "label2", "value2").Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_registry_namespace", "image-from-registry-namespace").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_name", "image-from-name").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_tag", "image-from-version").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_registry_host", "image-from-registry-host.test").Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddAuth", "pull-user", "pull-pass", "image-from-registry-host.test").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddAuth", "push-user", "push-pass", "myregistry.test").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddPushAuth", "push-user", "push-pass").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithPushAfterBuild")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithPullParentImage")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithRemoveAfterPush")

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildContext", []*builder.DockerDriverContextOptions{
					{Path: "/path/to/file"},
					{Path: "/path/to/file2"},
					{
						Git: &builder.DockerDriverGitContextOptions{
							Repository: "repo",
							Reference:  "main",
							Path:       "path",
							Auth: &builder.DockerDriverGitContextAuthOptions{
								Username: "user",
								Password: "pass",
							},
						},
					},
				}).Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithResponse", os.Stdout, "output-prefix")

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithUseNormalizedNamed")

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("Run", context.TODO()).Return(nil)
			},
			assertFunc: func(t *testing.T, driver DockerDriverer) {
				driver.(*godockerbuilder.MockGoDockerBuildDriver).AssertExpectations(t)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error context not defined on build options",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			},
			ctx: context.TODO(),
			image: &image.Image{
				Name:              "image",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "myregistry.test",
				Parent: &image.Image{
					Name:              "image-from-name",
					Version:           "image-from-version",
					RegistryNamespace: "image-from-registry-namespace",
					RegistryHost:      "image-from-registry-host.test",
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
						"pvar2": "pvalue2",
					},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					Tags: []string{"tag1", "tag2"},
				},
			},
			options: &image.BuildDriverOptions{
				PushImageAfterBuild: true,
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
				},
				PullAuthUsername: "pull-user",
				PullAuthPassword: "pull-pass",
				PushAuthUsername: "push-user",
				PushAuthPassword: "push-pass",
				BuilderOptions: &builder.BuilderOptions{
					Dockerfile: "Dockerfile.test",
				},
			},
			prepareAssertFunc: func(driver DockerDriverer) {
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithImageName", "myregistry.test/namespace/image:version")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithDockerfile", "Dockerfile.test")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "pvar1", "pvalue1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "pvar2", "pvalue2").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "var1", "value1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag1"}).Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag2"}).Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_registry_namespace", "image-from-registry-namespace").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_name", "image-from-name").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_tag", "image-from-version").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_registry_host", "image-from-registry-host.test").Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddAuth", "pull-user", "pull-pass", "image-from-registry-host.test").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddAuth", "push-user", "push-pass", "myregistry.test").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddPushAuth", "push-user", "push-pass").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithPushAfterBuild")
			},
			err: errors.New(errContext, "Docker building context has not been defined on build options"),
		},
		{
			desc: "Testing error when there is not found any docker build context definition",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockGoDockerBuildDriver(),
			},
			ctx: context.TODO(),
			image: &image.Image{
				Name:              "image",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "myregistry.test",
				Parent: &image.Image{
					Name:              "image-from-name",
					Version:           "image-from-version",
					RegistryNamespace: "image-from-registry-namespace",
					RegistryHost:      "image-from-registry-host.test",
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
						"pvar2": "pvalue2",
					},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					Tags: []string{"tag1", "tag2"},
				},
			},
			options: &image.BuildDriverOptions{
				PushImageAfterBuild: true,
				BuilderVarMappings: map[string]string{
					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
				},
				PullAuthUsername: "pull-user",
				PullAuthPassword: "pull-pass",
				PushAuthUsername: "push-user",
				PushAuthPassword: "push-pass",
				BuilderOptions: &builder.BuilderOptions{
					Dockerfile: "Dockerfile.test",
					Context:    []*builder.DockerDriverContextOptions{},
				},
			},
			prepareAssertFunc: func(driver DockerDriverer) {
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithImageName", "myregistry.test/namespace/image:version")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithDockerfile", "Dockerfile.test")
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "pvar1", "pvalue1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "pvar2", "pvalue2").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "var1", "value1").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag1"}).Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag2"}).Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_registry_namespace", "image-from-registry-namespace").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_name", "image-from-name").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_tag", "image-from-version").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddBuildArgs", "image_from_registry_host", "image-from-registry-host.test").Return(nil)

				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddAuth", "pull-user", "pull-pass", "image-from-registry-host.test").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddAuth", "push-user", "push-pass", "myregistry.test").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("AddPushAuth", "push-user", "push-pass").Return(nil)
				driver.(*godockerbuilder.MockGoDockerBuildDriver).On("WithPushAfterBuild")
			},
			err: errors.New(errContext, "Docker building context has not been defined on build options"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.driver.driver)
			}

			err := test.driver.Build(test.ctx, test.image, test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				if test.assertFunc != nil {
					test.assertFunc(t, test.driver.driver)
				} else {
					t.Error(test.desc, "missing assertFunc")
				}
			}

		})
	}
}
