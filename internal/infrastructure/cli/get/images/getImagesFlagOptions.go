package images

// getImagesFlagOptions is the options for the get images command
type getImagesFlagOptions struct {

	// Tree enables the output in tree format
	Tree bool
	// Filter is a list of filters to apply to the output
	Filter []string
}
