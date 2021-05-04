package dryrunpromoter

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

func TestPromote(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))

	tests := []struct {
		desc    string
		options *types.PromoteOptions
	}{
		{
			desc:    "Testing promote with nil options",
			options: nil,
		},
		{
			desc:    "Testing promote with empty options",
			options: &types.PromoteOptions{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			w.Reset()

			Promote(context.TODO(), test.options)
		})
	}
}
