package builders

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/stretchr/testify/assert"
)

// func TestPrintAll(t *testing.T) {
// 	builders := filter.NewFilter(&builders.Store{
// 		Builders: map[string]*builder.Builder{
// 			"builder1": {
// 				Name:   "builder1",
// 				Driver: "docker",
// 			},
// 			"builder2": {
// 				Name:   "builder2",
// 				Driver: "docker",
// 			},
// 		},
// 	})

// 	console := console.NewMockConsole()
// 	output := NewOutput(console, builders)

// 	console.On("PrintTable", [][]string{
// 		{"NAME", "DRIVER"},
// 		{"builder1", "docker"},
// 		{"builder2", "docker"},
// 	}).Return(nil)

// 	output.PrintAll()

// 	assert.True(t, console.AssertExpectations(t))
// }

func TestOutput(t *testing.T) {
	errContext := "(output::builders::Output::PlainOutput)"

	tests := []struct {
		desc            string
		output          *PlainOutput
		list            []*builder.Builder
		prepareMockFunc func(*PlainOutput)
		err             error
	}{
		{
			desc:   "Testing error on image plain output when writer is not define",
			output: NewPlainOutput(),
			err:    errors.New(errContext, "Builders output requires a writer"),
		},
		{
			desc: "Testing get builder in plain output",
			output: NewPlainOutput(
				WithWriter(console.NewMockConsole()),
			),
			list: []*builder.Builder{
				{
					Name:   "builder-1",
					Driver: "driver",
				},
				{
					Name:   "builder-2",
					Driver: "driver-2",
				},
				{},
			},
			prepareMockFunc: func(o *PlainOutput) {
				o.writer.(*console.MockConsole).On("PrintTable", [][]string{
					{"NAME", "DRIVER"},
					{"builder-1", "driver"},
					{"builder-2", "driver-2"},
				}).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

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
