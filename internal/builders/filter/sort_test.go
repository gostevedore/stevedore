package filter

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	tests := []struct {
		desc     string
		builders SortedBuilders
		res      int
	}{
		{
			desc: "Testing lenght of sorted builders",
			builders: []*builder.Builder{
				{
					Name: "builder1",
				},
				{
					Name: "builder2",
				},
			},
			res: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			l := test.builders.Len()

			assert.Equal(t, test.res, l)
		})
	}
}
func TestLess(t *testing.T) {
	tests := []struct {
		desc     string
		builders SortedBuilders
		i, j     int
		res      bool
	}{
		{
			desc: "Testing less on sorted builders",
			builders: []*builder.Builder{
				{
					Name: "builder1",
				},
				{
					Name: "builder2",
				},
			},
			i:   0,
			j:   1,
			res: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			l := test.builders.Less(test.i, test.j)

			assert.Equal(t, test.res, l)
		})
	}
}
func TestSwap(t *testing.T) {

	tests := []struct {
		desc     string
		builders SortedBuilders
		i, j     int
		res      SortedBuilders
	}{
		{
			desc: "Testing swap builders",
			builders: []*builder.Builder{
				{
					Name: "builder1",
				},
				{
					Name: "builder2",
				},
			},
			i: 0,
			j: 1,
			res: []*builder.Builder{
				{
					Name: "builder2",
				},
				{
					Name: "builder1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.builders.Swap(test.i, test.j)

			assert.Equal(t, test.res, test.builders)
		})
	}

}
