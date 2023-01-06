package configuration

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	handler "github.com/gostevedore/stevedore/internal/handler/create/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	output "github.com/gostevedore/stevedore/internal/infrastructure/configuration/output/file"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		desc            string
		entrypoint      *CreateConfigurationEntrypoint
		options         *Options
		prepareMockFunc func()
		err             error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.entrypoint.Execute(context.TODO(), test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			}
		})
	}
}

func TestPrepareHnadlerOptions(t *testing.T) {
	errContext := "(entrypoint::create::configuration::prepareHandlerOptions)"

	tests := []struct {
		desc       string
		entrypoint *CreateConfigurationEntrypoint
		options    *Options
		res        *handler.Options
		err        error
	}{
		{
			desc: "Testing create configuration entrypoint error preparing handler options when options are not provided",
			entrypoint: NewCreateConfigurationEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			err: errors.New(errContext, "Create configuration entrypoint requires options to prepare handler options"),
		},
		{
			desc:       "Testing prepare handler options into create configuration entrypoint",
			entrypoint: NewCreateConfigurationEntrypoint(),
			options: &Options{
				BuildersPath:                 "builderspath",
				Concurrency:                  5,
				CredentialsFormat:            "credentialsformat",
				CredentialsLocalStoragePath:  "credentialslocalstoragepath",
				CredentialsStorageType:       "credentialsstoragetype",
				EnableSemanticVersionTags:    true,
				Force:                        true,
				ImagesPath:                   "imagespath",
				LogPathFile:                  "logpathfile",
				PushImages:                   true,
				SemanticVersionTagsTemplates: []string{"tmpl1"},
			},
			res: &handler.Options{
				BuildersPath:                 "builderspath",
				Concurrency:                  5,
				CredentialsFormat:            "credentialsformat",
				CredentialsLocalStoragePath:  "credentialslocalstoragepath",
				CredentialsStorageType:       "credentialsstoragetype",
				EnableSemanticVersionTags:    true,
				ImagesPath:                   "imagespath",
				LogPathFile:                  "logpathfile",
				PushImages:                   true,
				SemanticVersionTagsTemplates: []string{"tmpl1"},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.prepareHandlerOptions(test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestGetConfigurationFileName(t *testing.T) {
	errContext := "(entrypoint::create::configuration::getConfigurationFileName)"

	tests := []struct {
		desc       string
		entrypoint *CreateConfigurationEntrypoint
		options    *Options
		res        string
		err        error
	}{
		{
			desc: "Testing create configuration entrypoint error getting configuration file name when options are not provided",
			entrypoint: NewCreateConfigurationEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			err: errors.New(errContext, "Create configuration entrypoint requires options to get configuration file name"),
		},
		{
			desc:       "Testing default file in get configuration file name into create configuratino entrypoint",
			entrypoint: NewCreateConfigurationEntrypoint(),
			options:    &Options{},
			res:        "stevedore.yaml",
			err:        &errors.Error{},
		},
		{
			desc:       "Testing custom configuration file in get configuration file name into create configuratino entrypoint",
			entrypoint: NewCreateConfigurationEntrypoint(),
			options: &Options{
				ConfigurationFilePath: "custom.yaml",
			},
			res: "custom.yaml",
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.getConfigurationFileName(test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestCreateOutputWriter(t *testing.T) {

	errContext := "(entrypoint::create::configuration::createOutputWriter)"

	tests := []struct {
		desc       string
		entrypoint *CreateConfigurationEntrypoint
		options    *Options
		res        configuration.ConfigurationWriter
		err        error
	}{
		{
			desc:       "Testing create configuration entrypoint error creating output writer when filesystem is not provided",
			entrypoint: NewCreateConfigurationEntrypoint(),
			err:        errors.New(errContext, "Create configuration entrypoint requires a filesystem to create the output writer"),
		},
		{
			desc: "Testing create configuration entrypoint error creating output writer when options are not provided",
			entrypoint: NewCreateConfigurationEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			err: errors.New(errContext, "Create configuration entrypoint requires options to create the output writer"),
		},
		{
			desc: "Testing create configuration entrypoint output writer",
			entrypoint: NewCreateConfigurationEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			options: &Options{},
			err:     &errors.Error{},
			res:     &output.ConfigurationFileSafePersist{},
		},
		{
			desc: "Testing create configuration entrypoint output writer with force enabled",
			entrypoint: NewCreateConfigurationEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			options: &Options{
				Force: true,
			},
			err: &errors.Error{},
			res: &output.ConfigurationFilePersist{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.createOutputWriter(test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}
}
