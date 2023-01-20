package configuration

import (
	"context"
	"io"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/create/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	errContext := "(handler::create::configuration::Handler)"
	tests := []struct {
		desc            string
		handler         *CreateConfigurationHandler
		options         *Options
		prepareMockFunc func(Applicationer)
		err             error
	}{
		{
			desc: "Testing error on create configuration handler when options parameters is not provided",
			handler: NewCreateConfigurationHandler(
				WithApplication(application.NewMockCreateConfigurationApplication()),
			),
			err: errors.New(errContext, "Create configuration handler requires the options parameter"),
		},
		{
			desc: "Testing create configuration handler",
			handler: NewCreateConfigurationHandler(
				WithApplication(application.NewMockCreateConfigurationApplication()),
			),
			options: &Options{
				BuildersPath:                 "builderspath",
				Concurrency:                  10,
				CredentialsEncryptionKey:     "credentialsencryptionkey",
				CredentialsFormat:            "credentialsformat",
				CredentialsLocalStoragePath:  "credentialslocalstoragepath",
				CredentialsStorageType:       "credentialsstoragetype",
				EnableSemanticVersionTags:    true,
				ImagesPath:                   "imagespath",
				LogPathFile:                  "logpathfile",
				PushImages:                   true,
				SemanticVersionTagsTemplates: []string{"tmpl1"},
			},
			prepareMockFunc: func(a Applicationer) {
				a.(*application.MockCreateConfigurationApplication).On(
					"Run",
					context.TODO(),
					&configuration.Configuration{
						BuildersPath: "builderspath",
						Concurrency:  10,
						Credentials: &configuration.CredentialsConfiguration{
							EncryptionKey:    "credentialsencryptionkey",
							Format:           "credentialsformat",
							LocalStoragePath: "credentialslocalstoragepath",
							StorageType:      "credentialsstoragetype",
						},
						EnableSemanticVersionTags:    true,
						ImagesPath:                   "imagespath",
						LogPathFile:                  "logpathfile",
						LogWriter:                    io.Discard,
						PushImages:                   true,
						SemanticVersionTagsTemplates: []string{"tmpl1"},
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
