package schedule

import (
	"context"
)

// MockJobber
type MockJobber struct {
	run bool
}

func (j *MockJobber) Run(ctx context.Context) {
	j.run = true
}
