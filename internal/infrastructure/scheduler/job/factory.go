package job

import "github.com/gostevedore/stevedore/internal/infrastructure/scheduler"

// JobFactory is a factory for creating jobs
type JobFactory struct{}

// NewJobFactory returns a new job factory
func NewJobFactory() *JobFactory {
	return &JobFactory{}
}

// New returns a new build job constructor
func (f *JobFactory) New(command Commander) scheduler.Jobber {
	return NewJob(command)
}
