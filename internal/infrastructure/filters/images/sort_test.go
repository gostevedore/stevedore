package images

import (
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	tests := []struct {
		desc   string
		images SortedImages
		res    int
	}{
		{
			desc: "Testing lenght of sorted images",
			images: []*image.Image{
				{
					Name: "image1",
				},
				{
					Name: "image2",
				},
			},
			res: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			l := test.images.Len()

			assert.Equal(t, test.res, l)
		})
	}
}
func TestLess(t *testing.T) {
	tests := []struct {
		desc   string
		images SortedImages
		i, j   int
		res    bool
	}{
		{
			desc: "Testing less on sorted images only by name",
			images: []*image.Image{
				{
					Name: "image1",
				},
				{
					Name: "image2",
				},
			},
			i:   0,
			j:   1,
			res: true,
		},
		{
			desc: "Testing less on sorted images when name is equal",
			images: []*image.Image{
				{
					Name:    "image",
					Version: "1",
				},
				{
					Name:    "image",
					Version: "2",
				},
			},
			i:   0,
			j:   1,
			res: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			l := test.images.Less(test.i, test.j)

			assert.Equal(t, test.res, l)
		})
	}
}
func TestSwap(t *testing.T) {

	tests := []struct {
		desc   string
		images SortedImages
		i, j   int
		res    SortedImages
	}{
		{
			desc: "Testing swap images",
			images: []*image.Image{
				{
					Name: "image1",
				},
				{
					Name: "image2",
				},
			},
			i: 0,
			j: 1,
			res: []*image.Image{
				{
					Name: "image2",
				},
				{
					Name: "image1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.images.Swap(test.i, test.j)

			assert.Equal(t, test.res, test.images)
		})
	}

}
