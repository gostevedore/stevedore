package images

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

type ImageVersionFilter struct{}

func NewImageVersionFilter() ImageVersionFilter {
	filter := ImageVersionFilter{}
	return filter
}

// Select return a sublist of images that its version value is item. operation is not used
func (f ImageVersionFilter) Select(images []*image.Image, operation string, item string) ([]*image.Image, error) {
	list := []*image.Image{}

	for _, i := range images {
		if i.Version == item {
			list = append(list, i)
		}
	}

	return list, nil
}
