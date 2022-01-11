// package engine

// import (
// 	"context"
// 	"testing"

// 	mockdriver "github.com/gostevedore/stevedore/internal/driver/mock"

// 	errors "github.com/apenella/go-common-utils/error"

// 	"github.com/stretchr/testify/assert"
// )

// // Run
// func TestRun(t *testing.T) {
// 	doneChan := make(chan bool)
// 	errChan := make(chan error)

// 	backgroundContext := context.Background()
// 	cancelContext, cancel := context.WithCancel(backgroundContext)
// 	defer cancel()

// 	tests := []struct {
// 		desc string
// 		err  error
// 		job  *Job
// 		ctx  context.Context
// 	}{
// 		{
// 			desc: "Testing job build done properly",
// 			err:  &errors.Error{},
// 			ctx:  cancelContext,
// 			job: &Job{
// 				Driver: &mockdriver.MockDriver{},
// 				Done:   doneChan,
// 				Err:    errChan,
// 			},
// 		},
// 		{
// 			desc: "Testing job build finished with errors",
// 			err:  errors.New("(MockDriverRunErr)", "Error"),
// 			ctx:  cancelContext,
// 			job: &Job{
// 				Driver: &mockdriver.MockDriverErr{},
// 				Done:   doneChan,
// 				Err:    errChan,
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		var err error
// 		t.Log(test.desc)

// 		go func(e error) {
// 			select {
// 			case <-doneChan:
// 			case err = <-errChan:
// 			}

// 			if err != nil && assert.Error(t, err) {
// 				assert.Equal(t, e, err)
// 			}
// 		}(test.err)

// 		test.job.Run(test.ctx)
// 	}
// }
