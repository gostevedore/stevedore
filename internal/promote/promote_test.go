package promote

import (
	"bytes"
	"context"
	goerrors "errors"
	"io"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/stretchr/testify/assert"
)

func TestPromote(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))
	ctx := context.TODO()

	tests := []struct {
		desc    string
		skip    bool
		ctx     context.Context
		options *types.PromoteOptions
		err     error
		res     string
	}{
		{
			desc: "Testing a error promotion",
			skip: false,
			ctx:  ctx,
			options: &types.PromoteOptions{
				ImageName: "unexisting",
			},
			err: errors.New("(promote::Promote)", "Error promoting 'unexisting' image",
				errors.New("", "Error promoting 'unexisting' to 'unexisting'",
					errors.New("", "Error pushing 'docker.io/library/unexisting'",
						errors.New("", "Error tagging image 'unexisting' to 'docker.io/library/unexisting'",
							goerrors.New("Error response from daemon: No such image: unexisting:latest"))))),
		},
		{
			desc: "Testing a error promotion",
			skip: false,
			ctx:  ctx,
			options: &types.PromoteOptions{
				DryRun:    true,
				ImageName: "unexisting",
			},
			res: `dry_run: true
enable_semantic_version_tags: false
image_promote_name: ""
image_promote_registry_namespace: ""
image_promote_registry_host: ""
image_promote_tags: []
remove_promoted_tags: false
image_name: unexisting
output_prefix: ""
semantic_version_tags_templates: []
`,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if test.skip {
				t.Skip(test.desc)
			}

			err := Promote(test.ctx, test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, w.String())
			}
		})
	}
}
