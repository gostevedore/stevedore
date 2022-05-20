package filter

import "github.com/gostevedore/stevedore/internal/core/domain/builder"

// SortedBuilders implements sort.Interface based on the Builder Name field
type SortedBuilders []*builder.Builder

// Len returns the length of the SortedBuilders
func (b SortedBuilders) Len() int {
	return len(b)
}

// Less returns true if the Builder Name of the first Builder is less than the second
func (b SortedBuilders) Less(i, j int) bool {
	return b[i].Name < b[j].Name
}

// Swap swaps the two Builders at the given indices
func (b SortedBuilders) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
