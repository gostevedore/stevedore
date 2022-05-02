package build

// EntrypointOptions defines the options for the entrypoint that initialize a build application
type Options struct {
	// Concurrency is the number of images builds that can be excuted at the same time
	Concurrency int
	// Debug if is true debug mode is enabled: ???
	Debug bool
	// // DryRun is true if the build should be a dry run: ???
	// DryRun bool
}
