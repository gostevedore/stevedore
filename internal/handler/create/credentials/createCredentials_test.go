package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		desc    string
		handler *CreateCredentialsHandler
		err     error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.handler.Handler(context.TODO())
			if err != nil {
				assert.Equal(t, test.err, err)
			}
		})
	}
	assert.True(t, false)
}
