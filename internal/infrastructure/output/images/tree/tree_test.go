package images

import (
	"os"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	tests := []struct {
		desc   string
		output *TreeOutput
		list   []*image.Image
		err    error
	}{
		{
			desc: "Testing",
			output: NewTreeOutput(
				WithWriter(os.Stdout),
				WithGraph(
					&gdsexttree.Graph{},
				),
			),
			list: []*image.Image{
				{
					Name:              "root",
					Version:           "v1",
					RegistryHost:      "registry.test",
					RegistryNamespace: "test",
					Children: []*image.Image{
						{
							Name:              "l1",
							Version:           "v1",
							RegistryHost:      "registry.test",
							RegistryNamespace: "test",
							Children: []*image.Image{
								{
									Name:              "l2",
									Version:           "v1",
									RegistryHost:      "registry.test",
									RegistryNamespace: "test",
								},
							},
						},
					},
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		err := test.output.Output(test.list)
		if err != nil {
			assert.Equal(t, test.err, err)
		}
	}
}
