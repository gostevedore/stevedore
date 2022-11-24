package builders

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/stretchr/testify/assert"
)

func TestBuilderNameFilterSelect(t *testing.T) {

	tests := []struct {
		desc     string
		filter   BuilderNameFilter
		builders []*builder.Builder
		item     string
		res      []*builder.Builder
		err      error
	}{
		{
			desc:   "Testing select by builder name from a nil input slice",
			filter: NewBuilderNameFilter(),
			item:   "builder-name",
			res:    nil,
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by builder name from an empty input slice",
			filter: NewBuilderNameFilter(),
			item:   "builder-name",
			res:    []*builder.Builder{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by builder name",
			filter: NewBuilderNameFilter(),
			item:   "builder-name",
			builders: []*builder.Builder{
				{
					Name: "builder-name",
				},
				{
					Name: "builder-name",
				},
			},
			res: []*builder.Builder{
				{
					Name: "builder-name",
				},
				{
					Name: "builder-name",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing select an unexisting builder by builder name",
			filter: NewBuilderNameFilter(),
			item:   "builder-unexisting",
			res:    []*builder.Builder{},
			err:    &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			list, err := test.filter.Select(test.builders, "", test.item)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.ElementsMatch(t, test.res, list)
			}
		})
	}
}
