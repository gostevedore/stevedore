package job

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockJob is a mock of Jobber interface
type MockJob struct {
	mock.Mock
}

// NewMockJob creates a new mock job
func NewMockJob() *MockJob {
	return &MockJob{}
}

// Run is a mock implementation of Jobber.Run
func (j *MockJob) Run(ctx context.Context) {
	j.Called(ctx)
}

// Wait waits for the job to finish
func (j *MockJob) Wait() error {
	args := j.Mock.Called()
	return args.Error(0)
}

// Close closes the job channels
func (j *MockJob) Close() {
	j.Called()
}

// Done returns a channel that is closed when the job is done
func (j *MockJob) Done() <-chan struct{} {
	args := j.Mock.Called()
	return args.Get(0).(<-chan struct{})
}

// Err returns a channel that is closed when the job has an error
func (j *MockJob) Err() <-chan error {
	args := j.Mock.Called()
	return args.Get(0).(<-chan error)
}
