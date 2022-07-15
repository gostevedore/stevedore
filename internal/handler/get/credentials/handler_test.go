package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		desc    string
		handler *Handler
		err     error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)
			err := test.handler.Handler(context.TODO())
			assert.Equal(t, test.err.Error(), err.Error())
		})
	}
}
