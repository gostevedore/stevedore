package builders

import (
	"sort"

	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
)

// Filter is a filter that filters the builders
type Filter struct {
	*builders.Store
}

// NewFilter returns a new Filter
func NewFilter(builders *builders.Store) *Filter {
	return &Filter{builders}
}

// All return all builders
func (f *Filter) All() []*builder.Builder {
	var filtered []*builder.Builder
	for _, builder := range f.Store.Builders {
		filtered = append(filtered, builder)
	}

	sort.Sort(SortedBuilders(filtered))

	return filtered
}

// FilterByName return the builder that match to a gived name
func (f *Filter) FilterByName(name string) *builder.Builder {
	for _, builder := range f.Store.Builders {
		if builder.Name == name {
			return builder
		}
	}
	return nil
}

// FilterByDriver return a list builders that match to a gived driver
func (f *Filter) FilterByDriver(driver string) []*builder.Builder {
	var filtered []*builder.Builder
	for _, builder := range f.Store.Builders {
		if builder.Driver == driver {
			filtered = append(filtered, builder)
		}
	}

	sort.Sort(SortedBuilders(filtered))

	return filtered
}
