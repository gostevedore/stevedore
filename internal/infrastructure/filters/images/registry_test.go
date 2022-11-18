package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestImageRegistryFilterSelect(t *testing.T) {

	tests := []struct {
		desc   string
		filter ImageRegistryFilter
		images []*image.Image
		item   string
		res    []*image.Image
		err    error
	}{
		{
			desc:   "Testing select by image registry from a nil input slice",
			filter: NewImageRegistryFilter(),
			item:   "image-name",
			res:    nil,
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image registry from an empty input slice",
			filter: NewImageRegistryFilter(),
			item:   "image-name",
			res:    []*image.Image{},
			err:    &errors.Error{},
		},
		{
			desc:   "Testing select by image registry",
			filter: NewImageRegistryFilter(),
			item:   "registry1",
			images: []*image.Image{
				{
					Name:              "image-name",
					Version:           "v1",
					RegistryNamespace: "ns1",
					RegistryHost:      "registry1",
				},
				{
					Name:              "image-name",
					Version:           "v2",
					RegistryNamespace: "ns2",
					RegistryHost:      "registry2",
				},
			},
			res: []*image.Image{
				{
					Name:              "image-name",
					Version:           "v1",
					RegistryNamespace: "ns1",
					RegistryHost:      "registry1",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing select an unexisting image by image registry",
			filter: NewImageRegistryFilter(),
			item:   "ns-unexisting",
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
