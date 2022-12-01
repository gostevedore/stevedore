package images

// Options is the options for the get images command entrypoint
type Options struct {
	// Tree enables the output in tree format
	Tree bool
	// UserDockerNormalizedName when is true are used Docker normalized name references
	UseDockerNormalizedName bool
}
