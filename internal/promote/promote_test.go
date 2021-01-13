package promote

import (
	"bytes"
	"context"
	goerrors "errors"
	"io"
	"stevedore/internal/types"
	"stevedore/internal/ui/console"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
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
					errors.New("", "Error tagging 'unexisting' to 'docker.io/library/unexisting'",
						goerrors.New("Error response from daemon: No such image: unexisting:latest")))),
		},
		{
			desc: "Testing a error promotion",
			skip: false,
			ctx:  ctx,
			options: &types.PromoteOptions{
				DryRun:    true,
				ImageName: "unexisting",
			},
			res: "{DryRun:true EnableSemanticVersionTags:false ImagePromoteName: ImagePromoteRegistryNamespace: ImagePromoteRegistryHost: ImagePromoteTags:[] RemovePromotedTags:false ImageName:unexisting OutputPrefix: SemanticVersionTagsTemplate:[]}\n",
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
