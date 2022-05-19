package filter

import "github.com/gostevedore/stevedore/internal/images/image"

// SortedImages implements sort.Interface based on the image Name field and then image version field
type SortedImages []*image.Image

// Len returns the length of the SortedImages
func (images SortedImages) Len() int {
	return len(images)
}

// Less returns true if the Builder Name of the first Builder is less than the second
func (images SortedImages) Less(i, j int) bool {

	if images[i].Name == images[j].Name {
		return images[i].Version < images[j].Version
	}

	return images[i].Name < images[j].Name
}

// Swap swaps the two Builders at the given indices
func (images SortedImages) Swap(i, j int) {
	images[i], images[j] = images[j], images[i]
}
