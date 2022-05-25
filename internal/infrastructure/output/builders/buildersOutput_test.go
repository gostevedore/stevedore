package builders

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	filter "github.com/gostevedore/stevedore/internal/infrastructure/filters/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	"github.com/stretchr/testify/assert"
)

func TestPrintAll(t *testing.T) {
	builders := filter.NewFilter(&builders.Store{
		Builders: map[string]*builder.Builder{
			"builder1": {
				Name:   "builder1",
				Driver: "docker",
			},
			"builder2": {
				Name:   "builder2",
				Driver: "docker",
			},
		},
	})

	console := console.NewMockConsole()
	output := NewOutput(console, builders)

	console.On("PrintTable", [][]string{
		{"NAME", "DRIVER"},
		{"builder1", "docker"},
		{"builder2", "docker"},
	}).Return(nil)

	output.PrintAll()

	assert.True(t, console.AssertExpectations(t))
}
