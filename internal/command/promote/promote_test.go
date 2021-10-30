package promote

import (
	"context"
	"fmt"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	promotehandler "github.com/gostevedore/stevedore/internal/command/promote/handler"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/stretchr/testify/assert"
)

func TestPromote(t *testing.T) {

	tests := []struct {
		desc            string
		handler         HandlerPromoter
		args            []string
		prepareMockFunc func(HandlerPromoter)
		err             error
	}{
		// {
		// 	desc: "",
		// 	handler: nil,
		// 	args: nil,
		// 	prepareMockFunc: func() {},
		// 	err: &errors.Error{},
		// },
		{
			desc:    "Testing to promote a simple image",
			handler: promotehandler.NewHandlerMock(),
			args: []string{
				"--dry-run",
				"myregistryhost.com/namespace/ubuntu:20.04",
			},
			prepareMockFunc: func(h HandlerPromoter) {
				options := &promotehandler.HandlerOptions{
					DryRun:                       true,
					EnableSemanticVersionTags:    false,
					SourceImageName:              "myregistryhost.com/namespace/ubuntu:20.04",
					TargetImageName:              "",
					TargetImageRegistryNamespace: "",
					TargetImageRegistryHost:      "",
					TargetImageTags:              []string{},
					RemoveTargetImageTags:        false,
					SemanticVersionTagsTemplates: []string{},
					PromoteSourceImageTag:        false,
					RemoteSourceImage:            false,
				}

				h.(*promotehandler.HandlerMock).On("Handler", context.TODO(), options).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing to promote an image to a new registry host, registry namespace, with new name and multiple tags",
			handler: promotehandler.NewHandlerMock(),
			args: []string{
				"--dry-run",
				"myregistryhost.com/namespace/ubuntu:20.04",
				"--promote-image-name",
				"myubuntu",
				"--promote-image-namespace",
				"stable",
				"--promote-image-registry",
				"myprodregistryhost.com",
				"--promote-image-tag",
				"tag1",
				"--promote-image-tag",
				"tag2",
				"--remove-local-images-after-push",
			},
			prepareMockFunc: func(h HandlerPromoter) {
				options := &promotehandler.HandlerOptions{
					DryRun:                       true,
					EnableSemanticVersionTags:    false,
					SourceImageName:              "myregistryhost.com/namespace/ubuntu:20.04",
					TargetImageName:              "myubuntu",
					TargetImageRegistryNamespace: "stable",
					TargetImageRegistryHost:      "myprodregistryhost.com",
					TargetImageTags:              []string{"tag1", "tag2"},
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{},
					PromoteSourceImageTag:        false,
					RemoteSourceImage:            false,
				}

				h.(*promotehandler.HandlerMock).On("Handler", context.TODO(), options).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing to promote image and semver tags",
			handler: promotehandler.NewHandlerMock(),
			args: []string{
				"--dry-run",
				"myregistryhost.com/namespace/ubuntu:1.2.3",
				"--enable-semver-tags",
				"--semver-tags-template",
				"{{ .Major }}",
				"--semver-tags-template",
				"{{ .Major }}.{{ .Minor }}",
				"--promote-source-tags",
				"--remove-local-images-after-push",
				"--remote-source-image",
			},
			prepareMockFunc: func(h HandlerPromoter) {
				options := &promotehandler.HandlerOptions{
					DryRun:                       true,
					EnableSemanticVersionTags:    true,
					SourceImageName:              "myregistryhost.com/namespace/ubuntu:1.2.3",
					TargetImageName:              "",
					TargetImageRegistryNamespace: "",
					TargetImageRegistryHost:      "",
					TargetImageTags:              []string{},
					RemoveTargetImageTags:        true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}", "{{ .Major }}.{{ .Minor }}"},
					PromoteSourceImageTag:        true,
					RemoteSourceImage:            true,
				}

				h.(*promotehandler.HandlerMock).On("Handler", context.TODO(), options).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing to promote without image name",
			handler: promotehandler.NewHandlerMock(),
			args: []string{
				"--dry-run",
			},
			err: errors.New("(promote::RunE)", "Source images name must be provided"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.handler)
			}

			conf := &configuration.Configuration{}
			promoteHandler = test.handler
			promoteCmd := NewCommand(context.TODO(), conf)
			//promoteCmd.Command.PreRunE = func(cmd *cobra.Command, args []string) error { return nil }
			promoteCmd.Command.ParseFlags(test.args)

			err := promoteCmd.Command.RunE(promoteCmd.Command, test.args)

			if err != nil && assert.Error(t, err) {
				fmt.Println(err.Error())
				assert.Equal(t, test.err, err)
			} else {
				test.handler.(*promotehandler.HandlerMock).AssertExpectations(t)
			}

		})
	}
}

// func TestFunctionalPromote(t *testing.T) {
// 	var w bytes.Buffer
// 	console.SetWriter(io.Writer(&w))
// 	ctx := context.TODO()

// 	dockerPromoter := dockerpromoter.NewPromoteMock()
// 	promoteRepoDryRun := promoterepodryrun.NewDryRunPromote(dockerPromoter, io.Writer(&w))
// 	promoteRepoDocker := promoterepodocker.NewDockerPromote(dockerPromoter, io.Writer(&w))
// 	promoteRepoFactory := promoterepofactory.NewPromoteFactory()
// 	promoteRepoFactory.Register("dry-run", promoteRepoDryRun)
// 	promoteRepoFactory.Register("docker", promoteRepoDocker)

// 	config := &configuration.Configuration{
// 		TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
// 		BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
// 		LogPathFile:          "/dev/null",
// 		NumWorkers:           2,
// 		PushImages:           false,
// 		BuildOnCascade:       false,
// 		DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
// 	}
// 	credentialsStore := credentials.NewCredentialsStore()
// 	semverGenerator := semver.NewSemVerGenerator()
// 	promoteService := promoteservice.NewService(promoteRepoFactory, config, credentialsStore, semverGenerator)
// 	promoteHandler := promotehandler.NewHandler(promoteService)

// 	defaultVerbose := false
// 	defaultSkip := false

// 	tests := []struct {
// 		desc    string
// 		skip    bool
// 		verbose bool
// 		ctx     context.Context
// 		config  *configuration.Configuration
// 		args    []string
// 		res     map[string]int8
// 		err     error
// 	}{
// 		{
// 			desc:    "Testing to promote a simple image",
// 			verbose: defaultVerbose,
// 			skip:    defaultSkip,
// 			ctx:     ctx,
// 			config: &configuration.Configuration{
// 				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				LogPathFile:          "/dev/null",
// 				NumWorkers:           2,
// 				PushImages:           false,
// 				BuildOnCascade:       false,
// 				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
// 			},
// 			err: nil,
// 			args: []string{
// 				"--dry-run",
// 				"myregistryhost.com/namespace/ubuntu:20.04",
// 			},
// 			res: map[string]int8{
// 				"dry_run: true":                                         int8(0),
// 				"enable_semantic_version_tags: false":                   int8(0),
// 				"image_promote_name: \"\"":                              int8(0),
// 				"image_promote_registry_namespace: \"\"":                int8(0),
// 				"image_promote_registry_host: \"\"":                     int8(0),
// 				"image_promote_tags: []":                                int8(0),
// 				"remove_promoted_tags: false":                           int8(0),
// 				"image_name: myregistryhost.com/namespace/ubuntu:20.04": int8(0),
// 				"output_prefix: \"\"":                                   int8(0),
// 				"semantic_version_tags_templates: []":                   int8(0),
// 			},
// 		},
// 		{
// 			desc:    "Testing to promote an image to a new registry host, registry namespace, with new name and multiple tags",
// 			verbose: defaultVerbose,
// 			skip:    defaultSkip,
// 			ctx:     ctx,
// 			config: &configuration.Configuration{
// 				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				LogPathFile:          "/dev/null",
// 				NumWorkers:           2,
// 				PushImages:           false,
// 				BuildOnCascade:       false,
// 				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
// 			},
// 			err: nil,
// 			args: []string{
// 				"--dry-run",
// 				"myregistryhost.com/namespace/ubuntu:20.04",
// 				"--promote-image-name",
// 				"myubuntu",
// 				"--promote-image-namespace",
// 				"stable",
// 				"--promote-image-registry",
// 				"myprodregistryhost.com",
// 				"--promote-image-tag",
// 				"tag1",
// 				"--promote-image-tag",
// 				"tag2",
// 				"--remove-promote-tags",
// 			},
// 			res: map[string]int8{
// 				"dry_run: true":                                         int8(0),
// 				"enable_semantic_version_tags: false":                   int8(0),
// 				"image_promote_name: myubuntu":                          int8(0),
// 				"image_promote_registry_namespace: stable":              int8(0),
// 				"image_promote_registry_host: myprodregistryhost.com":   int8(0),
// 				"image_promote_tags:":                                   int8(0),
// 				"- tag1":                                                int8(0),
// 				"- tag2":                                                int8(0),
// 				"remove_promoted_tags: true":                            int8(0),
// 				"image_name: myregistryhost.com/namespace/ubuntu:20.04": int8(0),
// 				"output_prefix: \"\"":                                   int8(0),
// 				"semantic_version_tags_templates: []":                   int8(0),
// 			},
// 		},
// 		{
// 			desc:    "Testing to promote image and semver tags",
// 			verbose: defaultVerbose,
// 			skip:    defaultSkip,
// 			ctx:     ctx,
// 			config: &configuration.Configuration{
// 				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				LogPathFile:          "/dev/null",
// 				NumWorkers:           2,
// 				PushImages:           false,
// 				BuildOnCascade:       false,
// 				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
// 			},
// 			err: errors.New("(command::promoteHandler)", "Is required an image name"),
// 			args: []string{
// 				"--dry-run",
// 				"myregistryhost.com/namespace/ubuntu:1.2.3",
// 				"--enable-semver-tags",
// 				"--semver-tags-template",
// 				"{{ .Major }}",
// 				"--semver-tags-template",
// 				"{{ .Major }}.{{ .Minor }}",
// 			},
// 			res: map[string]int8{
// 				"dry_run: true":                          int8(0),
// 				"enable_semantic_version_tags: true":     int8(0),
// 				"image_promote_name: \"\"":               int8(0),
// 				"image_promote_registry_namespace: \"\"": int8(0),
// 				"image_promote_registry_host: \"\"":      int8(0),
// 				"image_promote_tags:":                    int8(0),
// 				"- 1.2.3":                                int8(0),
// 				"- \"1\"":                                int8(0),
// 				"- \"1.2\"":                              int8(0),
// 				"remove_promoted_tags: false":            int8(0),
// 				"image_name: myregistryhost.com/namespace/ubuntu:1.2.3": int8(0),
// 				"output_prefix: \"\"":              int8(0),
// 				"semantic_version_tags_templates:": int8(0),
// 				"- '{{ .Major }}'":                 int8(0),
// 				"- '{{ .Major }}.{{ .Minor }}'":    int8(0),
// 			},
// 		},

// 		{
// 			desc:    "Testing to promote image and semver tags getting config from file",
// 			verbose: defaultVerbose,
// 			skip:    defaultVerbose,
// 			ctx:     ctx,
// 			config: &configuration.Configuration{
// 				TreePathFile:                 filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				BuilderPathFile:              filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				LogPathFile:                  "/dev/null",
// 				NumWorkers:                   2,
// 				PushImages:                   false,
// 				BuildOnCascade:               false,
// 				DockerCredentialsDir:         filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				EnableSemanticVersionTags:    false,
// 				SemanticVersionTagsTemplates: []string{"{{.Major}}", "{{.Major}}.{{.Minor}}"},
// 			},
// 			err: errors.New("(command::promoteHandler)", "Is required an image name"),
// 			args: []string{
// 				"--dry-run",
// 				"myregistryhost.com/namespace/ubuntu:1.2.3",
// 				"--enable-semver-tags",
// 				"--config",
// 				"test/stevedore_config.yml",
// 			},
// 			res: map[string]int8{
// 				"image_name: myregistryhost.com/namespace/ubuntu:1.2.3": int8(0),
// 				"dry_run: true":                          int8(0),
// 				"enable_semantic_version_tags: true":     int8(0),
// 				"image_promote_name: \"\"":               int8(0),
// 				"image_promote_registry_namespace: \"\"": int8(0),
// 				"image_promote_registry_host: \"\"":      int8(0),
// 				"image_promote_tags:":                    int8(0),
// 				"- 1.2.3":                                int8(0),
// 				"- \"1\"":                                int8(0),
// 				"- \"1.2\"":                              int8(0),
// 				"remove_promoted_tags: false":            int8(0),
// 				"output_prefix: \"\"":                    int8(0),
// 				"semantic_version_tags_templates:":       int8(0),
// 				"- '{{.Major}}'":                         int8(0),
// 				"- '{{.Major}}.{{.Minor}}'":              int8(0),
// 			},
// 		},

// 		{
// 			desc:    "Testing to promote without image name",
// 			verbose: defaultVerbose,
// 			skip:    defaultSkip,
// 			ctx:     ctx,
// 			config: &configuration.Configuration{
// 				TreePathFile:         filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				BuilderPathFile:      filepath.Join(testBaseDir, "stevedore_config.yml"),
// 				LogPathFile:          "/dev/null",
// 				NumWorkers:           2,
// 				PushImages:           false,
// 				BuildOnCascade:       false,
// 				DockerCredentialsDir: filepath.Join(testBaseDir, "stevedore_config.yml"),
// 			},
// 			err: errors.New("(command::promoteHandler)", "Is required an image name"),
// 			args: []string{
// 				"--dry-run",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			if test.skip {
// 				t.Skip(test.desc)
// 			}
// 			t.Log(test.desc)

// 			w.Reset()

// 			// cmd := NewCommand(test.ctx, test.config)
// 			cmd := NewCommand(test.ctx, promoteHandler)
// 			cmd.Command.ParseFlags(test.args)
// 			err := cmd.Command.RunE(cmd.Command, test.args)
// 			if err != nil && assert.Error(t, err) {
// 				fmt.Println(err.Error())
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				if test.verbose {
// 					t.Log("\n", w.String())
// 				}

// 				wSplit := strings.Split(w.String(), "\n")
// 				assert.Equal(t, len(test.res), len(wSplit)-1, "Unexpected number of lines")
// 				for i := 0; i < len(wSplit)-1; i++ {

// 					line := wSplit[i]
// 					_, ok := test.res[line]

// 					assert.True(t, ok)
// 					delete(test.res, line)
// 				}
// 				assert.Equal(t, len(test.res), 0, "Not all expected lines has appeared")
// 			}

// 		})
// 	}

// }
