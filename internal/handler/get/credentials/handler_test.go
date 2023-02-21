package credentials

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {

	errContext := "(get::credentials::Handler)"

	tests := []struct {
		desc              string
		handler           *Handler
		prepareAssertFunc func(handler *Handler)
		err               error
	}{
		{
			desc:    "Testing error getting credentials without an application defined",
			handler: NewHandler(),
			err:     errors.New(errContext, "Handler application is not configured"),
		},
		{
			desc: "Testing get credentials application",
			handler: NewHandler(
				WithApplication(
					application.NewMockApplication(),
				),
			),
			prepareAssertFunc: func(handler *Handler) {
				handler.app.(*application.MockApplication).On("Run", context.TODO(), mock.Anything).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.handler)
			}

			err := test.handler.Handler(context.TODO())
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.handler.app.(*application.MockApplication).AssertExpectations(t)
			}
		})
	}
}
