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
		app             *CreateConfigurationApplication
		config          *configuration.Configuration
		prepareMockFunc func(a *CreateConfigurationApplication)
		err             error
	}{
		{
			desc: "Testing application create configuration",
			app: NewCreateConfigurationApplication(
				WithWrite(output.NewConfigurationMockOutput()),
			),
			config: &configuration.Configuration{
				BuildersPath: "mystevedore.yaml",
				Concurrency:  10,
				ImagesPath:   "mystevedore.yaml",
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      "local",
					LocalStoragePath: "mycredentials",
					Format:           "json",
				},
				LogPathFile:                  "mystevedore.log",
				PushImages:                   true,
				EnableSemanticVersionTags:    true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(a *CreateConfigurationApplication) {
				a.write.(*output.ConfigurationMockOutput).On(
					"Write",
					&configuration.Configuration{
						BuildersPath: "mystevedore.yaml",
						Concurrency:  10,
						ImagesPath:   "mystevedore.yaml",
						Credentials: &configuration.CredentialsConfiguration{
							StorageType:      "local",
							LocalStoragePath: "mycredentials",
							Format:           "json",
						},
						LogPathFile:                  "mystevedore.log",
						PushImages:                   true,
						EnableSemanticVersionTags:    true,
						SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
					},
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

			err := test.app.Run(context.TODO(), test.config)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.app.write.(*output.ConfigurationMockOutput).AssertExpectations(t)
			}
		})
	}
}
