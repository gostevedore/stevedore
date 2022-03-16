package build

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestVarListToMap(t *testing.T) {
	// TO DO
}

const (
	testBaseDir = "test"
)

func TestBuildHandler(t *testing.T) {
	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	globalSkip := false
	globalVerbose := false

	tests := []struct {
		desc    string
		skip    bool
		verbose bool
		ctx     context.Context
		config  *configuration.Configuration
		args    []string
		res     map[string]int8
		err     error
	}{
		{
			desc:    "Testing to build an image on cascade with persistent vars on the root",
			err:     &errors.Error{},
			skip:    globalSkip,
			verbose: globalVerbose,
			res: map[string]int8{
				"{BuilderName:builder_ns_ubuntu_18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:ubuntu ImageVersion:18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost:registry PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                             int8(0),
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registry ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":            int8(0),
				"{BuilderName:builder_ns_php-fpm-dev_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:php-fpm ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion:7.4-ubuntu18.04 ImageName:php-fpm-dev ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}": int8(0),
			},
			ctx: ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			args: []string{
				"ubuntu",
				"--image-version",
				"18.04",
				"--namespace",
				"ns",
				"--connection-local",
				"--cascade",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build an image with semantic version tags enabled and no push enabled",
			err:     &errors.Error{},
			skip:    globalSkip,
			verbose: globalVerbose,
			res: map[string]int8{
				"{BuilderName:builder_ns_semver-app_1.2.3-rc0.0 BuilderOptions:map[] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image_name image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:true ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:semver-app ImageVersion:1.2.3-rc0.0 NumWorkers:0 OutputPrefix: PersistentVars:map[] RegistryNamespace:ns RegistryHost: PushImages:false SemanticVersionTagsTemplate:[{{ .Major }}.{{ .Minor }}.{{ .Patch }}] Tags:[1.2.3] Vars:map[]}": int8(0),
			},
			ctx: ctx,
			config: &configuration.Configuration{
				TreePathFile:                 filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:              filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:                  "/dev/null",
				NumWorkers:                   2,
				PushImages:                   true,
				BuildOnCascade:               false,
				DockerCredentialsDir:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				SemanticVersionTagsTemplates: []string{configuration.DefaultSemanticVersionTagsTemplates},
			},
			args: []string{
				"semver-app",
				"--namespace",
				"ns",
				"--connection-local",
				"--enable-semver-tags",
				"--no-push",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build all image version on cascade",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_ubuntu_18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:ubuntu ImageVersion:18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost:registry PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                             int8(0),
				"{BuilderName:builder_ns_ubuntu_16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:ubuntu ImageVersion:16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost:registryX PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                            int8(0),
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registry ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":            int8(0),
				"{BuilderName:builder_ns_nginx_1.15-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:16.04 ImageName:nginx ImageVersion:1.15-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":             int8(0),
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:16.04 ImageName:php-fpm ImageVersion:7.4-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":           int8(0),
				"{BuilderName:builder_ns_php-cli_7.4-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:16.04 ImageName:php-cli ImageVersion:7.4-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":           int8(0),
				"{BuilderName:builder_ns_php-fpm-dev_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:php-fpm ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion:7.4-ubuntu18.04 ImageName:php-fpm-dev ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}": int8(0),
				"{BuilderName:builder_ns_php-fpm-dev_7.4-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:php-fpm ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion:7.4-ubuntu16.04 ImageName:php-fpm-dev ImageVersion:7.4-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}": int8(0),
				"{BuilderName:builder_ns_php-cli-dev_7.4-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:php-cli ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion:7.4-ubuntu16.04 ImageName:php-cli-dev ImageVersion:7.4-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}": int8(0),
			},
			args: []string{
				"ubuntu",
				"--namespace",
				"ns",
				"--connection-local",
				"--cascade",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build all image on cascade with depth",
			err:     &errors.Error{},
			skip:    globalSkip,
			verbose: globalVerbose,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_ubuntu_18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:ubuntu ImageVersion:18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost:registry PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                   int8(0),
				"{BuilderName:builder_ns_ubuntu_16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:ubuntu ImageVersion:16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost:registryX PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                  int8(0),
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registry ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":  int8(0),
				"{BuilderName:builder_ns_nginx_1.15-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:16.04 ImageName:nginx ImageVersion:1.15-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":   int8(0),
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:16.04 ImageName:php-fpm ImageVersion:7.4-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}": int8(0),
				"{BuilderName:builder_ns_php-cli_7.4-ubuntu16.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:16.04 ImageName:php-cli ImageVersion:7.4-ubuntu16.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:16.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}": int8(0),
			},
			args: []string{
				"ubuntu",
				"--namespace",
				"ns",
				"--connection-local",
				"--cascade",
				"--cascade-depth",
				"1",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build all images on cascade with wildcard images defined",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_ubuntu_18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:ubuntu ImageVersion:18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[ubuntu_version:18.04] RegistryNamespace:ns RegistryHost:registryX PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                                                                                                   int8(0),
				"{BuilderName:builder_ns_php-cli_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-cli ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:7.4 ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[]}":                                                                  int8(0),
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:7.4 pvar1:pvar1 ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm]}":                                     int8(0),
				"{BuilderName:builder_ns_php-fpm-dev_7.4-php-fpm7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:php-fpm ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion:7.4-ubuntu18.04 ImageName:php-fpm-dev ImageVersion:7.4-php-fpm7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:7.4 pvar1:pvar1 ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm-dev]}": int8(0),
			},
			args: []string{
				"ubuntu",
				"--namespace",
				"ns",
				"--cascade",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build an image using a wildcard version",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_php-fpm_wildcard2 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:wildcard2 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:wildcard2 pvar1:pvar1-wildcard] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm-wildcard]}": int8(0),
				"{BuilderName:builder_ns_php-fpm_wildcard1 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:wildcard1 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:wildcard1 pvar1:pvar1-wildcard] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm-wildcard]}": int8(0),
			},
			args: []string{
				"php-fpm",
				"--namespace",
				"ns",
				"-v",
				"wildcard1",
				"-v",
				"wildcard2",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing build error when version is not specified on a wildcard image",
			verbose: true,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{},
			err: errors.New("", "Error building image 'any-version-app'",
				errors.New("", "Error finding image 'any-version-app' and versions '[]'",
					errors.New("", "No image 'any-version-app' found on images tree"))),
			args: []string{
				"any-version-app",
			},
		},
		{
			desc:    "Testing to build an image with unexisting version and a wildcarded tree",
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			err: errors.New("(command::buildHandler)", "Error building image 'php-cli'", errors.New("(ImagesEngine::Build)", "Error finding image 'php-cli' and versions '[wildcard1]'",
				errors.New("", "No matching images to be build found for 'php-cli' (versions: [wildcard1])"))),
			args: []string{
				"php-cli",
				"--namespace",
				"ns",
				"-v",
				"wildcard1",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build an image setting variables from cli flags",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:7.4-flag pvar1:pvar1 ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm var_from_flag:7.3-ubuntu16.04-flag]}": int8(0),
			},
			args: []string{
				"php-fpm",
				"--namespace",
				"ns",
				"-v",
				"7.4",
				"--set",
				"var_from_flag=7.3-ubuntu16.04-flag",
				"--set-persistent",
				"php_version=7.4-flag",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build an image overriding the image name",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_myphp-fpm-name_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:myphp-fpm-name ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:7.4-flag pvar1:pvar1 ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm var_from_flag:7.3-ubuntu16.04-flag]}": int8(0),
			},
			args: []string{
				"php-fpm",
				"--namespace",
				"ns",
				"--image-name",
				"myphp-fpm-name",
				"-v",
				"7.4",
				"--set",
				"var_from_flag=7.3-ubuntu16.04-flag",
				"--set-persistent",
				"php_version=7.4-flag",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing error building giving an image name and cascade flags",
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: nil,
			args: []string{
				"php-fpm",
				"--namespace",
				"ns",
				"--image-name",
				"myphp-fpm-name",
				"--cascade",
				"--connection-local",
				"--num-workers",
				"1",
			},
			err: errors.New("(command::buildHandler)", "Could not override image name with build on cascade"),
		},
		{
			desc:    "Testing to build all versions from an image with a wildcard version defined",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_php-fpm_7.4-ubuntu18.04 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:7.4-ubuntu18.04 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:7.4 pvar1:pvar1 ubuntu_version:18.04] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm]}": int8(0),
			},
			args: []string{
				"php-fpm",
				"--namespace",
				"ns",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build wildcard version in cascade",
			err:     &errors.Error{},
			verbose: globalVerbose,
			skip:    globalSkip,
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:         filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
				LogPathFile:          "/dev/null",
				NumWorkers:           2,
				PushImages:           true,
				BuildOnCascade:       false,
				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_wildcard_image_tree_config.yml"),
			},
			res: map[string]int8{
				"{BuilderName:builder_ns_php-fpm_wildcard1 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:ubuntu ImageFromRegistryNamespace: ImageFromRegistryHost:registryX ImageFromVersion:18.04 ImageName:php-fpm ImageVersion:wildcard1 NumWorkers:0 OutputPrefix: PersistentVars:map[php_version:wildcard1 pvar1:pvar1-wildcard] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm-wildcard]}":             int8(0),
				"{BuilderName:builder_ns_php-fpm-dev_wildcard1 BuilderOptions:map[inventory:inventory/all playbook:site.yml] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:true ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName:php-fpm ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion:*-ubuntu18.04 ImageName:php-fpm-dev ImageVersion:wildcard1 NumWorkers:0 OutputPrefix: PersistentVars:map[pvar1:pvar1-overwrite-php-fpm-dev-wildcard] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[] Tags:[] Vars:map[var1:var1-php-fpm-dev-wildcard]}": int8(0),
			},
			args: []string{
				"php-fpm",
				"--namespace",
				"ns",
				"-v",
				"wildcard1",
				"--cascade",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build an image with inline builder",
			err:     &errors.Error{},
			skip:    globalSkip,
			verbose: globalVerbose,
			res: map[string]int8{
				"{BuilderName:builder_ns_app-inline-builder_1.2.3 BuilderOptions:map[] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:tag push_image_key:push_image] Cascade:false ConnectionLocal:true Dockerfile: DryRun:false EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:app-inline-builder ImageVersion:1.2.3 NumWorkers:0 OutputPrefix: PersistentVars:map[] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[{{ .Major }}.{{ .Minor }}.{{ .Patch }}] Tags:[] Vars:map[]}": int8(0),
			},
			ctx: ctx,
			config: &configuration.Configuration{
				TreePathFile:                 filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:              filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:                  "/dev/null",
				NumWorkers:                   2,
				PushImages:                   true,
				BuildOnCascade:               false,
				DockerCredentialsDir:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				SemanticVersionTagsTemplates: []string{configuration.DefaultSemanticVersionTagsTemplates},
			},
			args: []string{
				"app-inline-builder",
				"--namespace",
				"ns",
				"--connection-local",
				"--num-workers",
				"1",
			},
		},
		{
			desc:    "Testing to build an image with dry-run",
			err:     &errors.Error{},
			skip:    true,
			verbose: globalVerbose,
			res: map[string]int8{
				"{BuilderName:builder_ns_app-dry-run_1.2.3 BuilderOptions:map[context:map[path:.] variables_mapping:map[image_name_key:image image_tag_key:tag]] BuilderVarMappings:map[image_builder_label_key:image_builder_label image_builder_name_key:image_builder_name image_builder_registry_host_key:image_builder_registry_host image_builder_registry_namespace_key:image_builder_registry_namespace image_builder_tag_key:image_builder_tag image_extra_tags_key:image_extra_tags image_from_name_key:image_from_name image_from_registry_host_key:image_from_registry_host image_from_registry_namespace_key:image_from_registry_namespace image_from_tag_key:image_from_tag image_name_key:image_name image_registry_host_key:image_registry_host image_registry_namespace_key:image_registry_namespace image_tag_key:image_tag push_image_key:push_image] Cascade:false ConnectionLocal:false Dockerfile: DryRun:true EnableSemanticVersionTags:false ImageFromName: ImageFromRegistryNamespace: ImageFromRegistryHost: ImageFromVersion: ImageName:app-dry-run ImageVersion:1.2.3 NumWorkers:0 OutputPrefix: PersistentVars:map[] RegistryNamespace:ns RegistryHost: PushImages:true SemanticVersionTagsTemplate:[{{ .Major }}.{{ .Minor }}.{{ .Patch }}] Tags:[] Vars:map[]}": int8(0),
			},
			ctx: ctx,
			config: &configuration.Configuration{
				TreePathFile:                 filepath.Join(testBaseDir, "stevedore_config.yml"),
				BuilderPathFile:              filepath.Join(testBaseDir, "stevedore_config.yml"),
				LogPathFile:                  "/dev/null",
				NumWorkers:                   2,
				PushImages:                   true,
				BuildOnCascade:               false,
				DockerCredentialsDir:         filepath.Join(testBaseDir, "stevedore_config.yml"),
				SemanticVersionTagsTemplates: []string{configuration.DefaultSemanticVersionTagsTemplates},
			},
			args: []string{
				"app-dry-run",
				"--namespace",
				"ns",
				"--dry-run",
				"--num-workers",
				"1",
			},
		},
		{
			desc: "Testing to load an image tree with a cyclic dependency",
			err: errors.New("", "Error creating new image engine",
				errors.New("", "Error creating images tree"),
				errors.New("", "Error generating images graph"),
				errors.New("", "Error generating template tree from 'cyclic2' to 'cyclic1:prod'"),
				errors.New("", "Error generating template tree from 'cyclic3' to 'cyclic2:prod'"),
				errors.New("", "Error generating template tree from 'cyclic1' to 'cyclic3:prod'"),
				errors.New("", "Relationship from 'cyclic3:prod' to 'cyclic1:prod' could not be created"),
				errors.New("", "Cycle detected adding relationship from 'cyclic3:prod' to 'cyclic1:prod'")),
			skip:    true, // skipped due the test output is not deterministic
			verbose: globalVerbose,
			res:     map[string]int8{},
			ctx:     ctx,
			config: &configuration.Configuration{
				TreePathFile:                 filepath.Join(testBaseDir, "stevedore_config_cyclic.yml"),
				BuilderPathFile:              filepath.Join(testBaseDir, "stevedore_config_cyclic.yml"),
				LogPathFile:                  "/dev/null",
				NumWorkers:                   2,
				PushImages:                   true,
				BuildOnCascade:               false,
				DockerCredentialsDir:         filepath.Join(testBaseDir, "stevedore_config_cyclic.yml"),
				SemanticVersionTagsTemplates: []string{configuration.DefaultSemanticVersionTagsTemplates},
			},
			args: []string{
				"cyclic1",
				"--namespace",
				"ns",
				"--no-push",
				"--num-workers",
				"1",
				"--cascade",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.skip {
				t.Skip(test.desc)
			}

			w.Reset()

			cmd := NewCommand(test.ctx, test.config)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				if test.verbose {
					t.Log("\n verbose:\n", w.String())
				}

				wSplit := strings.Split(w.String(), "\n")
				assert.Equal(t, len(test.res), len(wSplit)-1, "Unexpected number of lines")
				for i := 0; i < len(wSplit)-1; i++ {

					line := wSplit[i]
					_, ok := test.res[line]

					assert.True(t, ok)
					delete(test.res, line)
				}
				assert.Equal(t, len(test.res), 0, "Not all expected lines has appeared")
			}
		})
	}
}
