package filter

import (
	"sort"

	"github.com/gostevedore/stevedore/internal/builders/store"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
)

// BuildersFilter is a filter that filters the builders
type BuildersFilter struct {
	*store.BuildersStore
}

// NewBuildersFilter returns a new BuildersFilter
func NewBuildersFilter(builders *store.BuildersStore) *BuildersFilter {
	return &BuildersFilter{builders}
}

// All return all builders
func (f *BuildersFilter) All() []*builder.Builder {
	var filtered []*builder.Builder
	for _, builder := range f.BuildersStore.Builders {
		filtered = append(filtered, builder)
	}

	sort.Sort(SortedBuilders(filtered))

	return filtered
}

// FilterByName return the builder that match to a gived name
func (f *BuildersFilter) FilterByName(name string) *builder.Builder {
	for _, builder := range f.BuildersStore.Builders {
		if builder.Name == name {
			return builder
		}
	}
	return nil
}

// FilterByDriver return a list builders that match to a gived driver
func (f *BuildersFilter) FilterByDriver(driver string) []*builder.Builder {
	var filtered []*builder.Builder
	for _, builder := range f.BuildersStore.Builders {
		if builder.Driver == driver {
			filtered = append(filtered, builder)
		}
	}

	sort.Sort(SortedBuilders(filtered))

	return filtered
}
