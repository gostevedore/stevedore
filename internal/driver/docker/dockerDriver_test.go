package dockerdriver

import (
	"context"
	"io"
	"os"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/build/varsmap"
	buildcontext "github.com/gostevedore/stevedore/internal/driver/docker/context"
	"github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder"
	"github.com/gostevedore/stevedore/internal/types"
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
			driver: godockerbuilder.NewMockDockerDriver(),
			writer: nil,
			err:    &errors.Error{},
			res: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
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
		options           *types.BuildOptions
		prepareAssertFunc func(DockerDriverer)
		assertFunc        func(DockerDriverer) bool
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
			desc: "Testing error building a docker image with nil options",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
			},
			err: errors.New(errContext, "To build an image is required a build options"),
		},
		{
			desc: "Testing error building a docker image with nil golang context",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
			},
			options: &types.BuildOptions{},
			err:     errors.New(errContext, "To build an image is required a golang context"),
		},
		{
			desc: "Testing error building a docker image with not defined image name",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
			},
			ctx:     context.TODO(),
			options: &types.BuildOptions{},
			err:     errors.New(errContext, "To build an image is required an image name"),
		},
		{
			desc: "Testing building a docker image",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
			},
			ctx: context.TODO(),
			options: &types.BuildOptions{
				ImageName:                  "image",
				ImageVersion:               "version",
				RegistryNamespace:          "namespace",
				RegistryHost:               "myregistry.test",
				ImageFromName:              "image-from-name",
				ImageFromVersion:           "image-from-version",
				ImageFromRegistryNamespace: "image-from-registry-namespace",
				ImageFromRegistryHost:      "image-from-registry-host.test",
				PushImages:                 true,
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
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
					"pvar2": "pvalue2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
				},
				Tags:             []string{"tag1", "tag2"},
				PullAuthUsername: "pull-user",
				PullAuthPassword: "pull-pass",
				PushAuthUsername: "push-user",
				PushAuthPassword: "push-pass",
				BuilderOptions: map[string]interface{}{
					builderConfOptionsDockerfileKey: "Dockerfile.test",
					builderConfOptionsContextKey: `
- path: "/path/to/file"
- path: "/path/to/file2"
- git:
    repository: repo
    reference: main
    path: path
    auth:
      username: user
      password: pass
`,
				},
			},
			prepareAssertFunc: func(driver DockerDriverer) {
				driver.(*godockerbuilder.MockDockerDriver).On("WithImageName", "myregistry.test/namespace/image:version")
				driver.(*godockerbuilder.MockDockerDriver).On("WithDockerfile", "Dockerfile.test")
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "pvar1", "pvalue1").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "pvar2", "pvalue2").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "var1", "value1").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag1"}).Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag2"}).Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_registry_namespace", "image-from-registry-namespace").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_name", "image-from-name").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_tag", "image-from-version").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_registry_host", "image-from-registry-host.test").Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("AddAuth", "pull-user", "pull-pass", "image-from-registry-host.test").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddAuth", "push-user", "push-pass", "myregistry.test").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddPushAuth", "push-user", "push-pass").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("WithPushAfterBuild")

				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildContext", []*buildcontext.DockerBuildContextOptions{
					{Path: "/path/to/file"},
					{Path: "/path/to/file2"},
					{
						Git: &buildcontext.GitContextOptions{
							Repository: "repo",
							Reference:  "main",
							Path:       "path",
							Auth: &buildcontext.GitContextAuthOptions{
								Username: "user",
								Password: "pass",
							},
						},
					},
				}).Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("Run", context.TODO()).Return(nil)
			},
			assertFunc: func(driver DockerDriverer) bool {
				return driver.(*godockerbuilder.MockDockerDriver).AssertExpectations(t)
			},
			err: &errors.Error{},
		},

		{
			desc: "Testing error context not defined on build options",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
			},
			ctx: context.TODO(),
			options: &types.BuildOptions{
				ImageName:                  "image",
				ImageVersion:               "version",
				RegistryNamespace:          "namespace",
				RegistryHost:               "myregistry.test",
				ImageFromName:              "image-from-name",
				ImageFromVersion:           "image-from-version",
				ImageFromRegistryNamespace: "image-from-registry-namespace",
				ImageFromRegistryHost:      "image-from-registry-host.test",
				PushImages:                 true,
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
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
					"pvar2": "pvalue2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
				},
				Tags:             []string{"tag1", "tag2"},
				PullAuthUsername: "pull-user",
				PullAuthPassword: "pull-pass",
				PushAuthUsername: "push-user",
				PushAuthPassword: "push-pass",
				BuilderOptions: map[string]interface{}{
					builderConfOptionsDockerfileKey: "Dockerfile.test",
				},
			},
			prepareAssertFunc: func(driver DockerDriverer) {
				driver.(*godockerbuilder.MockDockerDriver).On("WithImageName", "myregistry.test/namespace/image:version")
				driver.(*godockerbuilder.MockDockerDriver).On("WithDockerfile", "Dockerfile.test")
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "pvar1", "pvalue1").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "pvar2", "pvalue2").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "var1", "value1").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag1"}).Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag2"}).Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_registry_namespace", "image-from-registry-namespace").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_name", "image-from-name").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_tag", "image-from-version").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_registry_host", "image-from-registry-host.test").Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("AddAuth", "pull-user", "pull-pass", "image-from-registry-host.test").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddAuth", "push-user", "push-pass", "myregistry.test").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddPushAuth", "push-user", "push-pass").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("WithPushAfterBuild")
			},
			err: errors.New(errContext, "Docker building context has not been defined on build options"),
		},
		{
			desc: "Testing error when there is not found any docker build context definition",
			driver: &DockerDriver{
				driver: godockerbuilder.NewMockDockerDriver(),
			},
			ctx: context.TODO(),
			options: &types.BuildOptions{
				ImageName:                  "image",
				ImageVersion:               "version",
				RegistryNamespace:          "namespace",
				RegistryHost:               "myregistry.test",
				ImageFromName:              "image-from-name",
				ImageFromVersion:           "image-from-version",
				ImageFromRegistryNamespace: "image-from-registry-namespace",
				ImageFromRegistryHost:      "image-from-registry-host.test",
				PushImages:                 true,
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
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
					"pvar2": "pvalue2",
				},
				Vars: map[string]interface{}{
					"var1": "value1",
				},
				Tags:             []string{"tag1", "tag2"},
				PullAuthUsername: "pull-user",
				PullAuthPassword: "pull-pass",
				PushAuthUsername: "push-user",
				PushAuthPassword: "push-pass",
				BuilderOptions: map[string]interface{}{
					builderConfOptionsDockerfileKey: "Dockerfile.test",
					builderConfOptionsContextKey: `
git:
  repository: repo
  reference: main
  path: path
  auth:
    username: user
    password: pass
`,
				},
			},
			prepareAssertFunc: func(driver DockerDriverer) {
				driver.(*godockerbuilder.MockDockerDriver).On("WithImageName", "myregistry.test/namespace/image:version")
				driver.(*godockerbuilder.MockDockerDriver).On("WithDockerfile", "Dockerfile.test")
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "pvar1", "pvalue1").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "pvar2", "pvalue2").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "var1", "value1").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag1"}).Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddTags", []string{"myregistry.test/namespace/image:tag2"}).Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_registry_namespace", "image-from-registry-namespace").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_name", "image-from-name").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_tag", "image-from-version").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddBuildArgs", "image_from_registry_host", "image-from-registry-host.test").Return(nil)

				driver.(*godockerbuilder.MockDockerDriver).On("AddAuth", "pull-user", "pull-pass", "image-from-registry-host.test").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddAuth", "push-user", "push-pass", "myregistry.test").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("AddPushAuth", "push-user", "push-pass").Return(nil)
				driver.(*godockerbuilder.MockDockerDriver).On("WithPushAfterBuild")
			},
			err: errors.New(errContext, "There is no docker build context definition found on:\ncontext:\ngit:\n  repository: repo\n  reference: main\n  path: path\n  auth:\n    username: user\n    password: pass\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.driver.driver)
			}

			err := test.driver.Build(test.ctx, test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(test.driver.driver))
				} else {
					t.Error(test.desc, "missing assertFunc")
				}
			}

		})
	}
}

// func TestNewDockerDriver_X(t *testing.T) {

// 	var w buffer.Buffer
// 	console.SetWriter(io.Writer(&w))
// 	ctx := context.TODO()

// 	optionPersistentVar := "pvar"
// 	optionVar := "var"

// 	tests := []struct {
// 		desc    string
// 		options *types.BuildOptions
// 		context context.Context
// 		err     error
// 		res     *dockerbuild.DockerBuildCmd
// 	}{
// 		{
// 			desc:    "Testing new dockerBuilder with nil options",
// 			options: nil,
// 			context: nil,
// 			err:     errors.New("(build::NewDockerDriver)", "Build options are nil"),
// 			res:     nil,
// 		},
// 		{
// 			desc: "Testing new dockerBuilder with a nil context",
// 			options: &types.BuildOptions{
// 				BuilderOptions: map[string]interface{}{},
// 			},
// 			context: nil,
// 			err:     errors.New("(build::NewDockerDriver)", "Context is nil"),
// 			res:     nil,
// 		},
// 		{
// 			desc:    "Testing new dockerBuilder with a non image name provided",
// 			options: &types.BuildOptions{},
// 			context: ctx,
// 			err:     errors.New("(build::NewDockerDriver)", "Image name is not set"),
// 			res:     nil,
// 		},
// 		{
// 			desc: "Testing options without a docker building context defined",
// 			options: &types.BuildOptions{
// 				ImageName:      "ubuntu",
// 				BuilderOptions: map[string]interface{}{},
// 			},
// 			context: ctx,
// 			err:     errors.New("(build::NewDockerDriver)", "Docker building context has not been defined on build options"),
// 			res:     nil,
// 		},
// 		{
// 			desc: "Testing run docker builder",
// 			options: &types.BuildOptions{
// 				BuilderOptions: map[string]interface{}{
// 					"context": map[string]string{
// 						"path": "test/ubuntu",
// 					},
// 				},
// 				Dockerfile:        "Dockerfile.test",
// 				ImageName:         "ubuntu",
// 				ImageVersion:      "16.04",
// 				RegistryNamespace: "library",
// 				RegistryHost:      "registry.host",
// 				PersistentVars: map[string]interface{}{
// 					"pvar1": optionPersistentVar,
// 				},
// 				Vars: map[string]interface{}{
// 					"var1": optionVar,
// 				},
// 				Tags: []string{
// 					"tag1",
// 				},
// 				PushImages: true,
// 			},
// 			context: ctx,
// 			err:     nil,
// 			res: &dockerbuild.DockerBuildCmd{
// 				ImageName: "registry.host/library/ubuntu:16.04",
// 				ImageBuildOptions: &dockertypes.ImageBuildOptions{
// 					Tags: []string{
// 						"registry.host/library/ubuntu:tag1",
// 					},
// 					BuildArgs: map[string]*string{
// 						"pvar1": &optionPersistentVar,
// 						"var1":  &optionVar,
// 					},
// 					Dockerfile: "Dockerfile.test",
// 				},
// 				PushAfterBuild: true,
// 			},
// 		},
// 		{
// 			desc: "Testing run docker builder with dockerfile defined in builder options",
// 			options: &types.BuildOptions{
// 				BuilderOptions: map[string]interface{}{
// 					"context": map[string]string{
// 						"path": "test/ubuntu",
// 					},
// 					"dockerfile": "./ubuntu/Dockerfile",
// 				},
// 				ImageName:         "ubuntu",
// 				ImageVersion:      "16.04",
// 				RegistryNamespace: "library",
// 				RegistryHost:      "registry.host",
// 				PersistentVars: map[string]interface{}{
// 					"pvar1": optionPersistentVar,
// 				},
// 				Vars: map[string]interface{}{
// 					"var1": optionVar,
// 				},
// 				Tags: []string{
// 					"tag1",
// 				},
// 				PushImages: true,
// 			},
// 			context: ctx,
// 			err:     nil,
// 			res: &dockerbuild.DockerBuildCmd{
// 				ImageName: "registry.host/library/ubuntu:16.04",
// 				ImageBuildOptions: &dockertypes.ImageBuildOptions{
// 					Tags: []string{
// 						"registry.host/library/ubuntu:tag1",
// 					},
// 					BuildArgs: map[string]*string{
// 						"pvar1": &optionPersistentVar,
// 						"var1":  &optionVar,
// 					},
// 					Dockerfile: "./ubuntu/Dockerfile",
// 				},
// 				PushAfterBuild: true,
// 			},
// 		},
// 		{
// 			desc: "Testing run docker builder with dockerfile defined in builder options and build options",
// 			options: &types.BuildOptions{
// 				BuilderOptions: map[string]interface{}{
// 					"context": map[string]string{
// 						"path": "test/ubuntu",
// 					},
// 					"dockerfile": "./ubuntu/Dockerfile",
// 				},
// 				Dockerfile:        "Dockerfile.test",
// 				ImageName:         "ubuntu",
// 				ImageVersion:      "16.04",
// 				RegistryNamespace: "library",
// 				RegistryHost:      "registry.host",
// 				PersistentVars: map[string]interface{}{
// 					"pvar1": optionPersistentVar,
// 				},
// 				Vars: map[string]interface{}{
// 					"var1": optionVar,
// 				},
// 				Tags: []string{
// 					"tag1",
// 				},
// 				PushImages: true,
// 			},
// 			context: ctx,
// 			err:     nil,
// 			res: &dockerbuild.DockerBuildCmd{
// 				ImageName: "registry.host/library/ubuntu:16.04",
// 				ImageBuildOptions: &dockertypes.ImageBuildOptions{
// 					Tags: []string{
// 						"registry.host/library/ubuntu:tag1",
// 					},
// 					BuildArgs: map[string]*string{
// 						"pvar1": &optionPersistentVar,
// 						"var1":  &optionVar,
// 					},
// 					Dockerfile: "Dockerfile.test",
// 				},
// 				PushAfterBuild: true,
// 			},
// 		},

// 		{
// 			desc: "Testing run docker builder with dockerfile defined in builder options and build options with image from details",
// 			options: &types.BuildOptions{
// 				BuilderOptions: map[string]interface{}{
// 					"context": map[string]string{
// 						"path": "test/ubuntu",
// 					},
// 					"dockerfile": "./ubuntu/Dockerfile",
// 				},
// 				BuilderVarMappings: map[string]string{
// 					varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
// 					varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
// 					varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
// 					varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
// 					varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
// 					varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
// 					varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
// 					varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
// 					varsmap.VarMappingImageNameKey:                     varsmap.VarMappingImageNameDefaultValue,
// 					varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
// 					varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
// 					varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
// 				},
// 				Dockerfile:        "Dockerfile.test",
// 				ImageName:         "ubuntu",
// 				ImageVersion:      "16.04",
// 				RegistryNamespace: "library",
// 				RegistryHost:      "registry.host",

// 				ImageFromRegistryHost:      "registryfrom",
// 				ImageFromRegistryNamespace: "registryfromnamespace",
// 				ImageFromName:              "imagefromname",
// 				ImageFromVersion:           "imagefromtag",

// 				PersistentVars: map[string]interface{}{
// 					"pvar1": optionPersistentVar,
// 				},
// 				Vars: map[string]interface{}{
// 					"var1": optionVar,
// 				},
// 				Tags: []string{
// 					"tag1",
// 				},
// 				PushImages: true,
// 			},
// 			context: ctx,
// 			err:     nil,
// 			res: &dockerbuild.DockerBuildCmd{
// 				ImageName: "registry.host/library/ubuntu:16.04",
// 				ImageBuildOptions: &dockertypes.ImageBuildOptions{
// 					Tags: []string{
// 						"registry.host/library/ubuntu:tag1",
// 					},
// 					BuildArgs: map[string]*string{
// 						"image_from_name":               ptr.String("imagefromname"),
// 						"image_from_tag":                ptr.String("imagefromtag"),
// 						"image_from_registry_namespace": ptr.String("registryfromnamespace"),
// 						"image_from_registry_host":      ptr.String("registryfrom"),
// 						"pvar1":                         &optionPersistentVar,
// 						"var1":                          &optionVar,
// 					},
// 					Dockerfile: "Dockerfile.test",
// 				},
// 				PushAfterBuild: true,
// 			},
// 		},

// 		{
// 			desc: "Testing run docker builder with push images as false",
// 			options: &types.BuildOptions{
// 				BuilderOptions: map[string]interface{}{
// 					"context": map[string]string{
// 						"path": "test/ubuntu",
// 					},
// 				},
// 				Dockerfile:        "Dockerfile.test",
// 				ImageName:         "ubuntu",
// 				ImageVersion:      "16.04",
// 				RegistryNamespace: "library",
// 				RegistryHost:      "registry.host",
// 				PersistentVars: map[string]interface{}{
// 					"pvar1": optionPersistentVar,
// 				},
// 				Vars: map[string]interface{}{
// 					"var1": optionVar,
// 				},
// 				Tags: []string{
// 					"tag1",
// 				},
// 				PushImages: false,
// 			},
// 			context: ctx,
// 			err:     nil,
// 			res: &dockerbuild.DockerBuildCmd{
// 				ImageName: "registry.host/library/ubuntu:16.04",
// 				ImageBuildOptions: &dockertypes.ImageBuildOptions{
// 					Tags: []string{
// 						"registry.host/library/ubuntu:tag1",
// 					},
// 					BuildArgs: map[string]*string{
// 						"pvar1": &optionPersistentVar,
// 						"var1":  &optionVar,
// 					},
// 					Dockerfile: "Dockerfile.test",
// 				},
// 				PushAfterBuild: false,
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			builderer, err := NewDockerDriver(test.context, test.options)

// 			if err != nil && assert.Error(t, err) {
// 				t.Log(err.Error())
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				dockerbuildercmd := builderer.(*dockerbuild.DockerBuildCmd)
// 				assert.Equal(t, test.res.ImageName, dockerbuildercmd.ImageName, "Unexpected value")
// 				assert.Equal(t, test.res.PushAfterBuild, dockerbuildercmd.PushAfterBuild, "Unexpected value")
// 				assert.Equal(t, test.res.ImageBuildOptions.Tags, dockerbuildercmd.ImageBuildOptions.Tags, "Unexpected value")
// 				assert.Equal(t, test.res.ImageBuildOptions.BuildArgs, dockerbuildercmd.ImageBuildOptions.BuildArgs, "Unexpected value")
// 			}
// 		})
// 	}
// }

// func TestExtractDockerBuildContext(t *testing.T) {

// 	tests := []struct {
// 		desc    string
// 		context interface{}
// 		err     error
// 		res     dockercontext.DockerBuildContexter
// 	}{
// 		{
// 			desc:    "Testing extracting a nil docker builder context",
// 			context: nil,
// 			err:     errors.New("(build::ExtractDockerBuildContext) Docker building context is nil"),
// 			res:     nil,
// 		},
// 		{
// 			desc: "Testing extracting a docker builder context defined on a path",
// 			context: map[string]string{
// 				"path": "test/ubuntu",
// 			},
// 			err: nil,
// 			res: &dockercontextpath.PathBuildContext{
// 				Path: "ubuntu",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		dockerBuilderContext, err := ExtractDockerBuildContext(test.context)
// 		if err != nil && assert.Error(t, err) {
// 			assert.Equal(t, test.err, err)
// 		} else {
// 			assert.Equal(t, test.res, dockerBuilderContext, "Unexpected value")
// 		}
// 	}
// }
