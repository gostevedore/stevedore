package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestImageNameFilterSelect(t *testing.T) {

	tests := []struct {
		desc   string
		filter ImageNameFilter
		images []*image.Image
		item   string
		res    []*image.Image
		err    error
	}{
		{
			desc:   "Testing select by image name from a nil input slice",
			filter: NewImageNameFilter(),
			item:   "image-name",
			res:    nil,
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image name from an empty input slice",
			filter: NewImageNameFilter(),
			item:   "image-name",
			res:    []*image.Image{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image name",
			filter: NewImageNameFilter(),
			item:   "image-name",
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
			res: []*image.Image{
				{
					Name:    "image-name",
					Version: "v1",
				},
				{
					Name:    "image-name",
					Version: "v2",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing select an unexisting image by image name",
			filter: NewImageNameFilter(),
			item:   "image-unexisting",
			res:    []*image.Image{},
			err:    &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			list, err := test.filter.Select(test.images, "", test.item)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.ElementsMatch(t, test.res, list)
			}
		})
	}
}
