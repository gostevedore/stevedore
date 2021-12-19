package job

// BuildJobFactory is a factory for creating build jobs
type BuildJobFactory struct{}

// NewBuildJobFactory returns a new build job factory
func NewBuildJobFactory() *BuildJobFactory {
	return &BuildJobFactory{}
}

// New returns a new build job constructor
func (f *BuildJobFactory) New(command BuildCommander) *BuildJob {
	return NewBuildJob(command)
}
