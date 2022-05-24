package builders

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	errContext := "(builders::Store)"

	tests := []struct {
		desc     string
		err      error
		builders *Store
		builder  *builder.Builder
		res      map[string]*builder.Builder
	}{
		{
			desc:     "Testing add a builder",
			builders: NewStore(),
			err:      &errors.Error{},
			builder: &builder.Builder{
				Name: "first",
			},
			res: map[string]*builder.Builder{
				"first": {Name: "first"},
			},
		},

		{
			desc: "Testing error adding already existing builder",
			builders: &Store{
				Builders: map[string]*builder.Builder{
					"first": {Name: "first"},
				},
			},
			err: errors.New(errContext, "Builder 'first' already exist"),
			builder: &builder.Builder{
				Name: "first",
			},
			res: map[string]*builder.Builder{
				"first": {Name: "first"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builders.Store(test.builder)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.builders.Builders)
			}
		})
	}
}

func TestFind(t *testing.T) {
	errContext := "(builders::GetBuilder)"

	tests := []struct {
		desc     string
		err      error
		builders *Store
		builder  string
		res      *builder.Builder
	}{
		{
			desc: "Testing get a builder",
			builders: &Store{
				Builders: map[string]*builder.Builder{
					"first": {Name: "first"},
				},
			},
			err:     &errors.Error{},
			builder: "first",
			res: &builder.Builder{
				Name: "first",
			},
		},

		{
			desc: "Testing error getting an unexisting",
			builders: &Store{
				Builders: map[string]*builder.Builder{},
			},
			err:     errors.New(errContext, "Builder 'first' does not exists"),
			builder: "first",

			res: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.builders.Find(test.builder)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}
