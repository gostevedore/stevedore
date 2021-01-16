package schedule

import (
	"context"
	"testing"
	"time"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestNewDispatch(t *testing.T) {
	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	tests := []struct {
		desc       string
		err        error
		context    context.Context
		numWorkers int
	}{
		{
			desc:       "Testing a new dispatcher creation with an invalid number of workers",
			err:        errors.New("(schedule::NewDispatch)", "Invalid value for number of workers, it must be greater than zero"),
			context:    cancelContext,
			numWorkers: 0,
		},
		{
			desc:       "Testing a new dispatcher creation",
			err:        nil,
			context:    cancelContext,
			numWorkers: 1,
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		d, err := NewDispatch(test.context, test.numWorkers)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.NotNil(t, d.WorkerPool)
			assert.NotNil(t, d.inputJobQueue)
		}
	}

}

func TestStartDispatch(t *testing.T) {
	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	tests := []struct {
		desc     string
		err      error
		dispatch *Dispatch
	}{
		{
			desc: "Testing start a new dispatcher",
			err:  nil,
			dispatch: &Dispatch{
				context:       cancelContext,
				WorkerPool:    make(chan chan Jobber, 1),
				inputJobQueue: make(chan Jobber),
				NumWorkers:    1,
			},
		},
	}

	for _, test := range tests {

		t.Log(test.desc)

		err := test.dispatch.Start()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}

func TestEnqueue(t *testing.T) {
	backgroundContext := context.Background()
	cancelContext, cancel := context.WithCancel(backgroundContext)
	defer cancel()

	dispatch := &Dispatch{
		context:       cancelContext,
		WorkerPool:    make(chan chan Jobber, 1),
		inputJobQueue: make(chan Jobber),
		NumWorkers:    1,
	}

	go dispatch.Start()

	job := &MockJobber{}
	dispatch.Enqueue(job)
	// sleep to wait a grace time for processing job
	time.Sleep(1 * time.Millisecond)
	assert.True(t, job.run)

}
