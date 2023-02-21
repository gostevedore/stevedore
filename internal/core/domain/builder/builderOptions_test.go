package builder

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestGetContext(t *testing.T) {
	errContext := "(core::domain::builder::BuilderOptions::GetContext)"
	tests := []struct {
		desc    string
		options *BuilderOptions
		res     []*DockerDriverContextOptions
		err     error
	}{
		{
			desc: "Testing error on get Docker driver build context when the context a not valid format",
			options: &BuilderOptions{
				Context: 1,
			},
			res: nil,
			err: errors.New(errContext, "Docker driver context options format is not valid"),
		},
		{
			desc: "Testing error on get Docker driver build context when the context could not be unmarshal",
			options: &BuilderOptions{
				Context: []interface{}{1},
			},
			res: nil,
			err: errors.New(errContext, "Docker driver context options could not be created.\nfound:\n'- 1\n'\n\n yaml: unmarshal errors:\n  line 1: cannot unmarshal !!int `1` into builder.DockerDriverContextOptions"),
		},
		{
			desc: "Testing get Docker driver build context when the context is a context itself",
			options: &BuilderOptions{
				Context: map[string]interface{}{
					"git": map[string]interface{}{
						"path":       "path",
						"repository": "repository",
						"reference":  "reference",
						"auth": map[string]interface{}{
							"username": "username",
							"password": "password",
						},
					},
				},
			},
			res: []*DockerDriverContextOptions{
				{
					Git: &DockerDriverGitContextOptions{
						Path:       "path",
						Repository: "repository",
						Reference:  "reference",
						Auth: &DockerDriverGitContextAuthOptions{
							Username: "username",
							Password: "password",
						},
					},
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get Docker driver build context when the context is a list of context",
			options: &BuilderOptions{
				Context: []interface{}{
					map[string]interface{}{
						"git": map[string]interface{}{
							"path":       "path",
							"repository": "repository",
							"reference":  "reference",
							"auth": map[string]interface{}{
								"username": "username",
								"password": "password",
							},
						},
					},
					map[string]interface{}{
						"path": "path",
					},
				},
			},
			res: []*DockerDriverContextOptions{
				{
					Git: &DockerDriverGitContextOptions{
						Path:       "path",
						Repository: "repository",
						Reference:  "reference",
						Auth: &DockerDriverGitContextAuthOptions{
							Username: "username",
							Password: "password",
						},
					},
				},
				{
					Path: "path",
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			contexts, err := test.options.GetContext()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.ElementsMatch(t, test.res, contexts)
			}
		})
	}
}
