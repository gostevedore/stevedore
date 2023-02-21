package images

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

type ImageNameFilter struct{}

func NewImageNameFilter() ImageNameFilter {
	filter := ImageNameFilter{}
	return filter
}

// Select return a sublist of images that its name value is item. operation is not used
func (f ImageNameFilter) Select(images []*image.Image, operation string, item string) ([]*image.Image, error) {
	list := []*image.Image{}

	for _, i := range images {
		if i.Name == item {
			list = append(list, i)
		}
	}

	return list, nil
}
