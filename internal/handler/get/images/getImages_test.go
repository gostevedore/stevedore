package images

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/images"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {

	errContext := "(handler::get::images::Handler)"

	tests := []struct {
		desc            string
		handler         *GetImagesHandler
		prepareMockFunc func(Applicationer)
		options         *Options
		err             error
	}{
		{
			desc:    "Testing error on get images handler when options is not defined",
			handler: &GetImagesHandler{},
			err:     errors.New(errContext, "Get images handler requires handler options"),
		},
		{
			desc:    "Testing error get images when application is not defined",
			handler: NewGetImagesHandler(),
			options: &Options{},
			err:     errors.New(errContext, "Get images handler requires an application"),
		},
		{
			desc: "Testing get images handler",
			handler: NewGetImagesHandler(
				WithApplication(application.NewMockGetImagesApplication()),
			),
			options: &Options{
				Filter: []string{
					"name=a",
				},
			},
			prepareMockFunc: func(a Applicationer) {
				a.(*application.MockGetImagesApplication).On("Run",
					context.TODO(),
					&application.Options{
						Filter: []string{
							"name=a",
						},
					},
					// application OptionsFunc
					mock.AnythingOfType("[]images.OptionsFunc"),
				).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil && test.handler.app != nil {
				test.prepareMockFunc(test.handler.app)
			}

			err := test.handler.Handler(context.TODO(), test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
