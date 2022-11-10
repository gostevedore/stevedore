package images

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

type ImageRegistryFilter struct{}

func NewImageRegistryFilter() ImageRegistryFilter {
	filter := ImageRegistryFilter{}
	return filter
}

// Select return a sublist of images that its registry value is item. operation is not used
func (f ImageRegistryFilter) Select(images []*image.Image, operation string, item string) ([]*image.Image, error) {
	list := []*image.Image{}

	for _, i := range images {
		if i.RegistryHost == item {
			list = append(list, i)
		}
	}

	return list, nil
}
