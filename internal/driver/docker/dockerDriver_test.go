package dockerdriver

import (
	"context"
	"io"
	"testing"

	"github.com/gostevedore/stevedore/internal/types"

	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	dockerbuild "github.com/apenella/go-docker-builder/pkg/build"
	"github.com/apenella/go-docker-builder/pkg/build/context/path"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/zap/buffer"
)

func TestNewDockerDriver(t *testing.T) {

	var w buffer.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	optionPersistentVar := "pvar"
	optionVar := "var"

	tests := []struct {
		desc    string
		options *types.BuildOptions
		context context.Context
		err     error
		res     *dockerbuild.DockerBuildCmd
	}{
		{
			desc:    "Testing new dockerBuilder with nil options",
			options: nil,
			context: nil,
			err:     errors.New("(build::NewDockerDriver)", "Build options are nil"),
			res:     nil,
		},
		{
			desc: "Testing new dockerBuilder with a nil context",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{},
			},
			context: nil,
			err:     errors.New("(build::NewDockerDriver)", "Context is nil"),
			res:     nil,
		},
		{
			desc:    "Testing new dockerBuilder with a non image name provided",
			options: &types.BuildOptions{},
			context: ctx,
			err:     errors.New("(build::NewDockerDriver)", "Image name is not set"),
			res:     nil,
		},
		{
			desc: "Testing options without a docker building context defined",
			options: &types.BuildOptions{
				ImageName:      "ubuntu",
				BuilderOptions: map[string]interface{}{},
			},
			context: ctx,
			err:     errors.New("(build::NewDockerDriver)", "Docker building context has not been defined on build options"),
			res:     nil,
		},
		{
			desc: "Testing run docker builder",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"context": map[string]string{
						"path": "ubuntu",
					},
				},
				Dockerfile:        "Dockerfile.test",
				ImageName:         "ubuntu",
				ImageVersion:      "16.04",
				RegistryNamespace: "library",
				RegistryHost:      "registry.host",
				PersistentVars: map[string]interface{}{
					"pvar1": optionPersistentVar,
				},
				Vars: map[string]interface{}{
					"var1": optionVar,
				},
				Tags: []string{
					"tag1",
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &dockerbuild.DockerBuildCmd{
				DockerBuildOptions: &dockerbuild.DockerBuildOptions{
					ImageName: "registry.host/library/ubuntu:16.04",
					Tags: []string{
						"registry.host/library/ubuntu:tag1",
					},
					BuildArgs: map[string]*string{
						"pvar1": &optionPersistentVar,
						"var1":  &optionVar,
					},
					Dockerfile:     "Dockerfile.test",
					PushAfterBuild: true,
					DockerBuildContext: &path.PathBuildContext{
						Path: "ubuntu",
					},
				},
			},
		},
		{
			desc: "Testing run docker builder with dockerfile defined in builder options",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"context": map[string]string{
						"path": "ubuntu",
					},
					"dockerfile": "./ubuntu/Dockerfile",
				},
				ImageName:         "ubuntu",
				ImageVersion:      "16.04",
				RegistryNamespace: "library",
				RegistryHost:      "registry.host",
				PersistentVars: map[string]interface{}{
					"pvar1": optionPersistentVar,
				},
				Vars: map[string]interface{}{
					"var1": optionVar,
				},
				Tags: []string{
					"tag1",
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &dockerbuild.DockerBuildCmd{
				DockerBuildOptions: &dockerbuild.DockerBuildOptions{
					ImageName: "registry.host/library/ubuntu:16.04",
					Tags: []string{
						"registry.host/library/ubuntu:tag1",
					},
					BuildArgs: map[string]*string{
						"pvar1": &optionPersistentVar,
						"var1":  &optionVar,
					},
					Dockerfile:     "./ubuntu/Dockerfile",
					PushAfterBuild: true,
					DockerBuildContext: &path.PathBuildContext{
						Path: "ubuntu",
					},
				},
			},
		},
		{
			desc: "Testing run docker builder with dockerfile defined in builder options and build options",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"context": map[string]string{
						"path": "ubuntu",
					},
					"dockerfile": "./ubuntu/Dockerfile",
				},
				Dockerfile:        "Dockerfile.test",
				ImageName:         "ubuntu",
				ImageVersion:      "16.04",
				RegistryNamespace: "library",
				RegistryHost:      "registry.host",
				PersistentVars: map[string]interface{}{
					"pvar1": optionPersistentVar,
				},
				Vars: map[string]interface{}{
					"var1": optionVar,
				},
				Tags: []string{
					"tag1",
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &dockerbuild.DockerBuildCmd{
				DockerBuildOptions: &dockerbuild.DockerBuildOptions{
					ImageName: "registry.host/library/ubuntu:16.04",
					Tags: []string{
						"registry.host/library/ubuntu:tag1",
					},
					BuildArgs: map[string]*string{
						"pvar1": &optionPersistentVar,
						"var1":  &optionVar,
					},
					Dockerfile:     "Dockerfile.test",
					PushAfterBuild: true,
					DockerBuildContext: &path.PathBuildContext{
						Path: "ubuntu",
					},
				},
			},
		},

		{
			desc: "Testing run docker builder with dockerfile defined in builder options and build options with image from details",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"context": map[string]string{
						"path": "ubuntu",
					},
					"dockerfile": "./ubuntu/Dockerfile",
				},
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
				Dockerfile:        "Dockerfile.test",
				ImageName:         "ubuntu",
				ImageVersion:      "16.04",
				RegistryNamespace: "library",
				RegistryHost:      "registry.host",

				ImageFromRegistryHost:      "registryfrom",
				ImageFromRegistryNamespace: "registryfromnamespace",
				ImageFromName:              "imagefromname",
				ImageFromVersion:           "imagefromtag",

				PersistentVars: map[string]interface{}{
					"pvar1": optionPersistentVar,
				},
				Vars: map[string]interface{}{
					"var1": optionVar,
				},
				Tags: []string{
					"tag1",
				},
				PushImages: true,
			},
			context: ctx,
			err:     nil,
			res: &dockerbuild.DockerBuildCmd{
				DockerBuildOptions: &dockerbuild.DockerBuildOptions{
					ImageName: "registry.host/library/ubuntu:16.04",
					Tags: []string{
						"registry.host/library/ubuntu:tag1",
					},
					BuildArgs: map[string]*string{
						"image_from_name":               ptr.String("imagefromname"),
						"image_from_tag":                ptr.String("imagefromtag"),
						"image_from_registry_namespace": ptr.String("registryfromnamespace"),
						"image_from_registry_host":      ptr.String("registryfrom"),
						"pvar1":                         &optionPersistentVar,
						"var1":                          &optionVar,
					},
					Dockerfile:     "Dockerfile.test",
					PushAfterBuild: true,
					DockerBuildContext: &path.PathBuildContext{
						Path: "ubuntu",
					},
				},
			},
		},

		{
			desc: "Testing run docker builder with push images as false",
			options: &types.BuildOptions{
				BuilderOptions: map[string]interface{}{
					"context": map[string]string{
						"path": "ubuntu",
					},
				},
				Dockerfile:        "Dockerfile.test",
				ImageName:         "ubuntu",
				ImageVersion:      "16.04",
				RegistryNamespace: "library",
				RegistryHost:      "registry.host",
				PersistentVars: map[string]interface{}{
					"pvar1": optionPersistentVar,
				},
				Vars: map[string]interface{}{
					"var1": optionVar,
				},
				Tags: []string{
					"tag1",
				},
				PushImages: false,
			},
			context: ctx,
			err:     nil,
			res: &dockerbuild.DockerBuildCmd{
				DockerBuildOptions: &dockerbuild.DockerBuildOptions{
					ImageName: "registry.host/library/ubuntu:16.04",
					Tags: []string{
						"registry.host/library/ubuntu:tag1",
					},
					BuildArgs: map[string]*string{
						"pvar1": &optionPersistentVar,
						"var1":  &optionVar,
					},
					Dockerfile:     "Dockerfile.test",
					PushAfterBuild: false,
					DockerBuildContext: &path.PathBuildContext{
						Path: "ubuntu",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			builderer, err := NewDockerDriver(test.context, test.options)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				dockerbuildercmd := builderer.(*dockerbuild.DockerBuildCmd)
				assert.Equal(t, test.res.DockerBuildOptions, dockerbuildercmd.DockerBuildOptions, "Unexpected value")
			}
		})
	}
}

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
// 				"path": "ubuntu",
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
