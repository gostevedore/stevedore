package configuration

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	output "github.com/gostevedore/stevedore/internal/infrastructure/configuration/output/mock"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	tests := []struct {
		desc            string
		app             *GetConfigurationApplication
		options         *Options
		prepareMockFunc func(a *GetConfigurationApplication)
		err             error
	}{
		{
			desc: "Testing application get configuration",
			app: NewGetConfigurationApplication(
				WithWrite(output.NewConfigurationMockOutput()),
			),
			options: &Options{
				Configuration: &configuration.Configuration{},
			},
			prepareMockFunc: func(a *GetConfigurationApplication) {
				a.write.(*output.ConfigurationMockOutput).On(
					"Write",
					&configuration.Configuration{},
				).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil && test.app != nil {
				test.prepareMockFunc(test.app)
			}

			err := test.app.Run(context.TODO(), test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.app.write.(*output.ConfigurationMockOutput).AssertExpectations(t)
			}
		})
	}
}
