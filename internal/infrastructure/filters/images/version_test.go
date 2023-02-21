package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestImageVersionFilterSelect(t *testing.T) {

	tests := []struct {
		desc   string
		filter ImageVersionFilter
		images []*image.Image
		item   string
		res    []*image.Image
		err    error
	}{
		{
			desc:   "Testing select by image version from a nil input slice",
			filter: NewImageVersionFilter(),
			images: nil,
			item:   "image-name",
			res:    []*image.Image{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image version from an empty input slice",
			filter: NewImageVersionFilter(),
			images: []*image.Image{},
			item:   "image-name",
			res:    []*image.Image{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image version",
			filter: NewImageVersionFilter(),
			images: []*image.Image{
				{
					Name:    "image-name",
					Version: "v1",
				},
				{
					Name:    "image-name",
					Version: "v2",
				},
			},
			item: "v1",
			res: []*image.Image{
				{
					Name:    "image-name",
					Version: "v1",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing select an unexisting image by image version",
			filter: NewImageVersionFilter(),
			images: []*image.Image{
				{
					Name:    "image-name",
					Version: "v1",
				},
				{
					Name:    "image-name",
					Version: "v2",
				},
			},
			item: "v3",
			res:  []*image.Image{},
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {

		list, err := test.filter.Select(test.images, "", test.item)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.ElementsMatch(t, test.res, list)
		}
	}
}
