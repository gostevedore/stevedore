package job

import (
	"context"
	"errors"
	"testing"

	"github.com/gostevedore/stevedore/internal/service/build/command"
	"github.com/stretchr/testify/assert"
)

func TestRunDone(t *testing.T) {
	t.Log("Testing run a job that finishes properly")

	build := command.NewMockBuildCommand()
	build.Mock.On("Execute", context.TODO()).Return(nil)

	job := NewJob(build)

	go job.Run(context.TODO())

	select {
	case <-job.Done():
		build.Mock.AssertExpectations(t)
	case <-job.Err():
		assert.Fail(t, "Job should not return an error")
	}

}

func TestRunErro(t *testing.T) {
	t.Log("Testing run a job that returns an error")

	build := command.NewMockBuildCommand()
	build.Mock.On("Execute", context.TODO()).Return(errors.New("Test error"))

	job := NewJob(build)

	go job.Run(context.TODO())

	select {
	case <-job.Done():
		assert.Fail(t, "Job should not finish properly")
	case <-job.Err():
		build.Mock.AssertExpectations(t)
	}

}
