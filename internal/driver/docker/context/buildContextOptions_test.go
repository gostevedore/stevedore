package context

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestGenerateBuildContextOptions(t *testing.T) {

	errContext := "(DockerBuildContextFactory::GenerateDockerBuildContext)"

	tests := []struct {
		desc    string
		context interface{}
		res     *DockerBuildContextOptions
		err     error
	}{
		{
			desc:    "Testing error nil context",
			context: nil,
			res:     nil,
			err:     errors.New(errContext, "Docker build context options are expected to build an image"),
		},
		{
			desc: "Testing error for invalid context definition",
			context: `
- path: path
`,
			res: nil,
			err: errors.New(errContext, "Docker build context options are not properly configured\n found:\n\n- path: path\n\n", &yaml.TypeError{[]string{"line 2: cannot unmarshal !!seq into context.DockerBuildContextOptions"}}),
		},
		{
			desc: "Testing generate build context options",
			context: `
path: path1
`,
			res: &DockerBuildContextOptions{

				Path: "path1",
			},
			err: &errors.Error{},
		},

		{
			desc: "Testing generate docker build context from git",
			context: `
git:
    repository: repo
    reference: main
    path: path
    auth:
      username: user
      password: pass
`,
			res: &DockerBuildContextOptions{

				Git: &GitContextOptions{
					Repository: "repo",
					Reference:  "main",
					Path:       "path",
					Auth: &GitContextAuthOptions{
						Username: "user",
						Password: "pass",
					},
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			res, err := GenerateBuildContextOptions(test.context)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}

		})
	}
}
