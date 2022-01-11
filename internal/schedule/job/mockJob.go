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

// Close closes the job channels
func (j *MockJob) Close() {
	j.Called()
}

// Done returns a channel that is closed when the job is done
func (j *MockJob) Done() <-chan struct{} {
	j.Mock.Called()
	return make(<-chan struct{})
}

// Err returns a channel that is closed when the job has an error
func (j *MockJob) Err() <-chan error {
	j.Mock.Called()
	return make(<-chan error)
}
