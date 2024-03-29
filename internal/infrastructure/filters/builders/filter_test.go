package builders

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	tests := []struct {
		desc     string
		Builders *Filter
		res      []*builder.Builder
	}{
		{
			desc: "Testing filter to get all builders",
			Builders: NewFilter(
				&builders.Store{
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
				},
			),
			res: []*builder.Builder{
				{
					Name:   "builder1",
					Driver: "docker",
				},
				{
					Name:   "builder2",
					Driver: "docker",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.Builders.All()
			assert.Equal(t, test.res, res)
		})
	}
}

func TestFilterByName(t *testing.T) {
	tests := []struct {
		desc     string
		name     string
		Builders *Filter
		res      *builder.Builder
	}{
		{
			desc: "Testing filter by name",
			name: "test",
			Builders: NewFilter(
				&builders.Store{
					Builders: map[string]*builder.Builder{
						"test": {
							Name: "test",
						},
					},
				},
			),
			res: &builder.Builder{
				Name: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.Builders.FilterByName(test.name)
			assert.Equal(t, test.res, res)
		})
	}
}

func TestFilterByDriver(t *testing.T) {
	tests := []struct {
		desc     string
		driver   string
		Builders *Filter
		res      []*builder.Builder
	}{
		{
			desc:   "Testing filter by driver",
			driver: "driver1",
			Builders: NewFilter(
				&builders.Store{
					Builders: map[string]*builder.Builder{
						"driver1": {
							Driver: "driver1",
						},
						"driver2": {
							Driver: "driver2",
						},
					},
				},
			),
			res: []*builder.Builder{
				{Driver: "driver1"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.Builders.FilterByDriver(test.driver)
			assert.Equal(t, test.res, res)
		})
	}
}
