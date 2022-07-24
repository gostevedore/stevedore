package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	tests := []struct {
		desc string
		app  *CreateCredentialsApplication
		err  error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.app.Run(context.TODO())
			assert.Equal(t, test.err.Error(), err.Error())
		})
	}
	assert.True(t, false)
}
