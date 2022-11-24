package builders

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/stretchr/testify/assert"
)

func TestBuilderDriverFilterSelect(t *testing.T) {

	tests := []struct {
		desc     string
		filter   BuilderDriverFilter
		builders []*builder.Builder
		item     string
		res      []*builder.Builder
		err      error
	}{
		{
			desc:   "Testing select by builder driver from a nil input slice",
			filter: NewBuilderDriverFilter(),
			item:   "builder-driver",
			res:    nil,
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by builder driver from an empty input slice",
			filter: NewBuilderDriverFilter(),
			item:   "builder-driver",
			res:    []*builder.Builder{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by builder driver",
			filter: NewBuilderDriverFilter(),
			item:   "builder-driver",
			builders: []*builder.Builder{
				{
					Name:   "builder-1",
					Driver: "builder-driver",
				},
				{
					Name:   "builder-2",
					Driver: "builder-driver",
				},
			},
			res: []*builder.Builder{
				{
					Name:   "builder-1",
					Driver: "builder-driver",
				},
				{
					Name:   "builder-2",
					Driver: "builder-driver",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing select an unexisting builder by builder driver",
			filter: NewBuilderDriverFilter(),
			item:   "driver-unexisting",
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
