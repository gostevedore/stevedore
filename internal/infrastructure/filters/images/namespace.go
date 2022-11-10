package images

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

type ImageNamespaceFilter struct{}

func NewImageNamespaceFilter() ImageNamespaceFilter {
	filter := ImageNamespaceFilter{}
	return filter
}

// Select return a sublist of images that its namespace value is item. operation is not used
func (f ImageNamespaceFilter) Select(images []*image.Image, operation string, item string) ([]*image.Image, error) {
	list := []*image.Image{}

	for _, i := range images {
		if i.RegistryNamespace == item {
			list = append(list, i)
		}
	}

	return list, nil
}
