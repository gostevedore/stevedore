package builders

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	errContext := "(handler::get::builders::Handler)"
	tests := []struct {
		desc            string
		handler         *GetBuildersHandler
		options         *Options
		prepareMockFunc func(Applicationer)
		err             error
	}{
		{
			desc:    "Testing error on get builders handler when options is not defined",
			handler: &GetBuildersHandler{},
			err:     errors.New(errContext, "Get builders handler requires handler options"),
		},
		{
			desc:    "Testing error get builders when application is not defined",
			handler: NewGetBuildersHandler(),
			options: &Options{},
			err:     errors.New(errContext, "Get builders handler requires an application"),
		},
		{
			desc: "Testing get builders handler",
			handler: NewGetBuildersHandler(
				WithApplication(application.NewMockGetBuildersApplication()),
			),
			options: &Options{
				Filter: []string{
					"name=a",
				},
			},
			prepareMockFunc: func(a Applicationer) {
				a.(*application.MockGetBuildersApplication).On("Run",
					context.TODO(),
					&application.Options{
						Filter: []string{
							"name=a",
						},
					},
					// application OptionsFunc
					mock.AnythingOfType("[]builders.OptionsFunc"),
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
