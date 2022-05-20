package context

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gitcontextbasicauth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	gitcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git"
	gitauth "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git/auth"
	pathcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/path"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDockerBuildContext(t *testing.T) {
	errContext := "(DockerBuildContextFactory::GenerateDockerBuildContext)"

	pathContext := pathcontext.NewPathBuildContext()
	pathContext.WithPath("context_path")

	gitContext := gitcontext.NewGitBuildContext()
	gitContext.WithRepository("repository")
	gitContext.WithReference("main")
	gitContext.WithPath("docker-context")
	gitContext.WithAuth(&gitcontextbasicauth.BasicAuth{
		Username: "user",
		Password: "password",
	})

	tests := []struct {
		desc    string
		options *builder.DockerDriverContextOptions
		factory *DockerBuildContextFactory
		context DockerBuildContexter
		err     error
	}{
		{
			desc:    "Testing error when options is nil",
			options: nil,
			factory: &DockerBuildContextFactory{},
			context: nil,
			err:     errors.New(errContext, "Docker build context options are required to generate a build context"),
		},
		{
			desc: "Testing to generate a docker build context from path",
			options: &builder.DockerDriverContextOptions{
				Path: "context_path",
			},
			factory: &DockerBuildContextFactory{},
			context: pathContext,
			err:     errors.New(errContext, "Docker build context options are required to generate a build context"),
		},
		{
			desc: "Testing error when git auth generator is not defined",
			options: &builder.DockerDriverContextOptions{
				Git: &builder.DockerDriverGitContextOptions{
					Repository: "repository",
					Reference:  "main",
					Path:       "docker-context",
					Auth: &builder.DockerDriverGitContextAuthOptions{
						Username: "user",
						Password: "password",
					},
				},
			},
			factory: &DockerBuildContextFactory{},
			context: nil,
			err:     errors.New(errContext, "Git auth generator is required to generate a git build context"),
		},
		{
			desc: "Testing to generate a docker build context from git repository",
			options: &builder.DockerDriverContextOptions{
				Git: &builder.DockerDriverGitContextOptions{
					Repository: "repository",
					Reference:  "main",
					Path:       "docker-context",
					Auth: &builder.DockerDriverGitContextAuthOptions{
						Username: "user",
						Password: "password",
					},
				},
			},
			factory: &DockerBuildContextFactory{
				gitAuth: gitauth.NewGitAuthFactory(nil),
			},
			context: gitContext,
			err:     &errors.Error{},
		},
		{
			desc: "Testing error generating docker build context from git repository without specifing a repository",
			options: &builder.DockerDriverContextOptions{
				Git: &builder.DockerDriverGitContextOptions{
					Repository: "",
				},
			},
			factory: &DockerBuildContextFactory{
				gitAuth: gitauth.NewGitAuthFactory(nil),
			},
			context: nil,
			err:     errors.New(errContext, "A repository must be specified on git build docker context"),
		},
		{
			desc: "Testing error creating docker build context git auth method",
			options: &builder.DockerDriverContextOptions{
				Git: &builder.DockerDriverGitContextOptions{
					Repository: "my-test-repository",
					Reference:  "main",
					Path:       "docker-context",
					Auth: &builder.DockerDriverGitContextAuthOptions{
						CredentialsID: "id",
					},
				},
			},
			factory: &DockerBuildContextFactory{
				gitAuth: gitauth.NewGitAuthFactory(nil),
			},
			context: nil,
			err:     errors.New(errContext, "\n\tCredentials store is expected when a credentials id is configured"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			context, err := test.factory.GenerateDockerBuildContext(test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.context, context)
			}
		})
	}
}
