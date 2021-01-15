package context

import (
	"testing"

	gitcontext "github.com/gostevedore/stevedore/internal/driver/docker/context/git"

	errors "github.com/apenella/go-common-utils/error"
	dockercontext "github.com/apenella/go-docker-builder/pkg/build/context"
	dockercontextgit "github.com/apenella/go-docker-builder/pkg/build/context/git"
	dockercontextpath "github.com/apenella/go-docker-builder/pkg/build/context/path"
	"github.com/stretchr/testify/assert"
)

func TestIsGitContext(t *testing.T) {

	tests := []struct {
		desc    string
		context *DockerBuildContext
		res     bool
	}{
		{
			desc: "Testing whether a path context is docker build context is a git context",
			context: &DockerBuildContext{
				Path: "path",
			},
			res: false,
		},
		{
			desc: "Testing whether a git context docker build context is a git context",
			context: &DockerBuildContext{
				Git: &gitcontext.GitContext{
					Repository: "repository",
				},
			},
			res: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			isit := test.context.IsGitContext()
			assert.Equal(t, test.res, isit)
		})
	}
}

func TestIsPathContext(t *testing.T) {

	tests := []struct {
		desc    string
		context *DockerBuildContext
		res     bool
	}{
		{
			desc: "Testing whether a path context is docker build context is a path context",
			context: &DockerBuildContext{
				Path: "path",
			},
			res: true,
		},
		{
			desc: "Testing whether a git context docker build context is a path context",
			context: &DockerBuildContext{
				Git: &gitcontext.GitContext{
					Repository: "repository",
				},
			},
			res: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			isit := test.context.IsPathContext()
			assert.Equal(t, test.res, isit)
		})
	}
}

func TestGetContextType(t *testing.T) {

	tests := []struct {
		desc    string
		context *DockerBuildContext
		res     uint8
	}{
		{
			desc: "Testing get a path docker build context",
			context: &DockerBuildContext{
				Path: "path",
			},
			res: PathContextType,
		},
		{
			desc: "Testing get a git docker build context",
			context: &DockerBuildContext{
				Git: &gitcontext.GitContext{
					Repository: "repository",
				},
			},
			res: GitContextType,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			ctype := test.context.GetContextType()
			assert.Equal(t, test.res, ctype)
		})
	}
}

func TestGenerateDockerBuildContext(t *testing.T) {
	tests := []struct {
		desc    string
		context map[string]interface{}
		res     dockercontext.DockerBuildContexter
		err     error
	}{
		{
			desc:    "Testing generate docker build context from a nil context",
			context: nil,
			res:     nil,
			err:     errors.New("(build::docker::context::GenerateDockerBuildContext)", "Unknown context type"),
		},
		{
			desc: "Testing generate docker build context from a path context",
			context: map[string]interface{}{
				"path": "mypath",
			},
			res: &dockercontextpath.PathBuildContext{
				Path: "mypath",
			},
			err: nil,
		},
		{
			desc: "Testing generate docker build context from a blank path context",
			context: map[string]interface{}{
				"path": "",
			},
			res: nil,
			err: errors.New("(build::docker::context::GenerateDockerBuildContext)", "Unknown context type"),
		},
		{
			desc: "Testing error when generate docker build context from a wrong git context",
			context: map[string]interface{}{
				"git": "mypath",
			},
			res: nil,
			err: errors.New("(build::docker::context::GenerateDockerBuildContext)", "Docker build context could not be unmarshalled to DockerBuildContext",
				errors.New("", "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `mypath` into gitcontext.GitContext")),
		},
		{
			desc: "Testing generate docker build context from a nil git context",
			context: map[string]interface{}{
				"git": nil,
			},
			res: nil,
			err: errors.New("(build::docker::context::GenerateDockerBuildContext)", "Unknown context type"),
		},
		{
			desc: "Testing generate docker build context from a blank repository on git context",
			context: map[string]interface{}{
				"git": &gitcontext.GitContext{
					Repository: "",
				},
			},
			res: nil,
			err: errors.New("(build::docker::context::GenerateDockerBuildContext)", "A repository must be specified on git build docker context"),
		},
		{
			desc: "Testing generate docker build context from a git context",
			context: map[string]interface{}{
				"git": &gitcontext.GitContext{
					Repository: "repository",
					Reference:  "reference",
				},
			},
			res: &dockercontextgit.GitBuildContext{
				Repository: "repository",
				Reference:  "reference",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			context, err := GenerateDockerBuildContext(test.context)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, context, test.res, "Unexpected context value")
			}
		})
	}
}
