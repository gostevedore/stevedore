package configuration

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/get/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		desc            string
		handler         *GetConfigurationHandler
		options         *Options
		prepareMockFunc func(Applicationer)
		err             error
	}{
		{
			desc: "Testing get configuration handler",
			handler: NewGetConfigurationHandler(
				WithApplication(
					application.NewMockGetConfigurationApplication(),
				),
			),
			options: &Options{
				Configuration: &configuration.Configuration{},
			},
			prepareMockFunc: func(a Applicationer) {
				a.(*application.MockGetConfigurationApplication).On(
					"Run",
					context.TODO(),
					&application.Options{
						Configuration: &configuration.Configuration{},
					},
					// application OptionsFunc
					mock.AnythingOfType("[]configuration.OptionsFunc"),
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
