package list

// SortedStringList implements sort.Interface based on string values
type SortedStringList []string

// Len returns the length of the SortedStringList
func (b SortedStringList) Len() int {
	return len(b)
}

// Less returns true if the Builder Name of the first Builder is less than the second
func (b SortedStringList) Less(i, j int) bool {
	return b[i] < b[j]
}

// Swap swaps the two Builders at the given indices
func (b SortedStringList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
