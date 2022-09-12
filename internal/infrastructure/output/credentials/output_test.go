package credentials

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	write "github.com/gostevedore/stevedore/internal/infrastructure/console"
	output "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials/types/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPrintAll(t *testing.T) {

	// mockOutput := output.NewMockOutput()
	// mockOutput.On("Output", mock.Anything).Return("", "", nil)

	tests := []struct {
		desc              string
		output            *Output
		badges            []*credentials.Badge
		prepareAssertFunc func(output *Output)
		err               error
	}{
		{
			desc: "Testing output for all credentials",
			output: &Output{
				methods: []Outputter{
					output.NewMockOutput(),
				},
				write: write.NewMockConsole(),
			},
			badges: []*credentials.Badge{
				{
					ID:       "id",
					Username: "username",
					Password: "password",
				},
			},
			prepareAssertFunc: func(o *Output) {
				method := o.methods[0]
				method.(*output.MockOutput).On("Output", mock.Anything).Return("type", "details", nil)
				o.write.(*write.MockConsole).On("PrintTable", [][]string{
					{"ID", "TYPE", "CRENDENTIALS"},
					{"id", "type", "details"},
				}).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.output)
			}

			err := test.output.Print(test.badges)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {

			}
		})
	}
}
