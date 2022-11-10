package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestImageNamespaceFilterSelect(t *testing.T) {

	tests := []struct {
		desc   string
		filter ImageNamespaceFilter
		images []*image.Image
		item   string
		res    []*image.Image
		err    error
	}{
		{
			desc:   "Testing select by image namespace from a nil input slice",
			filter: NewImageNamespaceFilter(),
			item:   "image-name",
			res:    nil,
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image namespace from an empty input slice",
			filter: NewImageNamespaceFilter(),
			item:   "image-name",
			res:    []*image.Image{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image namespace",
			filter: NewImageNamespaceFilter(),
			item:   "ns1",
			images: []*image.Image{
				{
					Name:              "image-name",
					Version:           "v1",
					RegistryNamespace: "ns1",
				},
				{
					Name:              "image-name",
					Version:           "v2",
					RegistryNamespace: "ns2",
				},
			},
			res: []*image.Image{
				{
					Name:              "image-name",
					Version:           "v1",
					RegistryNamespace: "ns1",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing select an unexisting image by image namespace",
			filter: NewImageNamespaceFilter(),
			item:   "ns-unexisting",
			res:    []*image.Image{},
			err:    &errors.Error{},
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
