package builders

import (
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
)

type BuilderNameFilter struct{}

func NewBuilderNameFilter() BuilderNameFilter {
	filter := BuilderNameFilter{}
	return filter
}

// Select return a sublist of images that its name value is item. operation is not used
func (f BuilderNameFilter) Select(builders []*builder.Builder, operation string, item string) ([]*builder.Builder, error) {
	list := []*builder.Builder{}

	for _, b := range builders {
		if b.Name == item {
			list = append(list, b)
		}
	}

	return list, nil
}
