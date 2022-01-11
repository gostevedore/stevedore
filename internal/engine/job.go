// package engine

// import (
// 	"context"

// 	"github.com/gostevedore/stevedore/internal/types"

// 	errors "github.com/apenella/go-common-utils/error"
// )

// // Job contains an image element and the option requiered to build it
// type Job struct {
// 	Driver types.Driverer
// 	Done   chan bool
// 	Err    chan error
// }

// // Run
// func (j *Job) Run(ctx context.Context) {

// 	doneChan := make(chan bool)
// 	errChan := make(chan error)

// 	if j.Driver == nil {
// 		j.Err <- errors.New("(engine::Job:Run)", "Driver is nil")
// 		return
// 	}

// 	go func() {
// 		err := j.Driver.Run(ctx)
// 		if err != nil {
// 			errChan <- err
// 			return
// 		}

// 		doneChan <- true
// 	}()

// 	select {
// 	case <-doneChan:
// 		j.Done <- true
// 	case err := <-errChan:
// 		j.Err <- err
// 	case <-ctx.Done():
// 		return
// 	}
// }
