package job

import "context"

// Job is a job that can be run
type Job struct {
	command Commander
	done    chan struct{}
	err     chan error
}

// NewJob creates a new job
func NewJob(command Commander) *Job {
	return &Job{
		command: command,
		done:    make(chan struct{}),
		err:     make(chan error),
	}
}

// Run runs the job
func (j *Job) Run(ctx context.Context) {

	err := j.command.Execute(ctx)
	if err != nil {
		j.err <- err
	}

	j.done <- struct{}{}
}

// Close closes the job channels
func (j *Job) Close() {
	close(j.done)
	close(j.err)
}

// Done returns a channel that is closed when the job is done
func (j *Job) Done() <-chan struct{} {
	return j.done
}

// Err returns a channel that is closed when the job has an error
func (j *Job) Err() <-chan error {
	return j.err
}
