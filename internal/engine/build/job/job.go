package job

import "context"

// BuildJob is a job that can be run
type BuildJob struct {
	command BuildCommander
	done    chan struct{}
	err     chan error
}

// NewBuildJob creates a new job
func NewBuildJob(command BuildCommander) *BuildJob {
	return &BuildJob{
		command: command,
		done:    make(chan struct{}),
		err:     make(chan error),
	}
}

// Run runs the job
func (j *BuildJob) Run(ctx context.Context) {

	err := j.command.Execute(ctx)
	if err != nil {
		j.err <- err
	}

	j.done <- struct{}{}
}

// Close closes the job channels
func (j *BuildJob) Close() {
	close(j.done)
	close(j.err)
}

// Done returns a channel that is closed when the job is done
func (j *BuildJob) Done() <-chan struct{} {
	return j.done
}

// Err returns a channel that is closed when the job has an error
func (j *BuildJob) Err() <-chan error {
	return j.err
}
