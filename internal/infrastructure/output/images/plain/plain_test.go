package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/stretchr/testify/assert"
)

func TestOuput(t *testing.T) {
	errContext := "(output::images::PlainOutput::Output)"

	tests := []struct {
		desc            string
		output          *PlainOutput
		list            []*image.Image
		prepareMockFunc func(*PlainOutput)
		err             error
	}{
		{
			desc:   "Testing error on image plain output when writer is not define",
			output: NewPlainOutput(),
			err:    errors.New(errContext, "Images plain text output requires a writer"),
		},
		{
			desc: "Testing get images in plain output",
			output: NewPlainOutput(
				WithWriter(console.NewMockConsole()),
			),
			list: []*image.Image{
				{
					Builder:  "builder",
					Children: []*image.Image{{Name: "child"}},
					Labels:   map[string]string{"label": "value"},
					Name:     "image1",
					Parent: &image.Image{
						Name:              "parent",
						Version:           "v1",
						RegistryHost:      "registry.test",
						RegistryNamespace: "library",
						Builder:           &builder.Builder{},
					},
					PersistentLabels:  map[string]string{"plabel": "pvalue"},
					PersistentVars:    map[string]interface{}{"pvar": "pvalue"},
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Tags:              []string{"tag"},
					Vars:              map[string]interface{}{"var": "value"},
					Version:           "v1",
				},
				{
					Name:              "parent",
					Version:           "v1",
					RegistryHost:      "registry.test",
					RegistryNamespace: "library",
					Builder:           &builder.Builder{},
				},
				{
					Name:              "",
					Version:           "",
					RegistryHost:      "",
					RegistryNamespace: "",
				},
			},
			prepareMockFunc: func(o *PlainOutput) {
				o.writer.(*console.MockConsole).On("PrintTable", [][]string{
					{"NAME", "VERSION", "REGISTRY", "NAMESPACE", "BUILDER", "PARENT"},
					{"image1", "v1", "registry.test", "namespace", "builder", "registry.test/library/parent:v1"},
					{"parent", "v1", "registry.test", "library", "<in-line>"},
					{"-", "-", "-", "-", "-"},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareMockFunc != nil && test.output != nil {
				test.prepareMockFunc(test.output)
			}

			err := test.output.Output(test.list)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.output.writer.(*console.MockConsole).AssertExpectations(t)
			}

		})
	}
}
