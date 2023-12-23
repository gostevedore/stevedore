package credentials

import (
	"context"
	"io"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	credentialsenvvarsstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	errContext := "(get::credentials::entrypoint::Execute)"

	filesystem := afero.NewMemMapFs()
	_ = filesystem.MkdirAll("/credentials", 0755)

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *Options
		args       []string
		conf       *configuration.Configuration
		err        error
	}{
		{
			desc:       "Testing error executing get credentials entrypoint when writer is not defined",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To execute the entrypoint, a writer is required"),
		},
		{
			desc: "Testing execute get credentials entrypoint",
			entrypoint: NewEntrypoint(
				WithWriter(console.NewConsole(io.Discard, nil)),
				WithFileSystem(filesystem),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			options: &Options{},
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      credentials.LocalStore,
					LocalStoragePath: "/credentials",
					Format:           "json",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error executing get credentials entrypoint when credentials configuration file is not defined",
			entrypoint: NewEntrypoint(
				WithWriter(console.NewConsole(io.Discard, nil)),
				WithFileSystem(filesystem),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			options: &Options{},
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      credentials.LocalStore,
					LocalStoragePath: "/unknown",
					Format:           "json",
				},
			},
			err: errors.New(errContext, "Error reading credentials file '/unknown'\n open /unknown: file does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}

func TestCreateCredentialsLocalStore(t *testing.T) {

	errContext := "(get::credentials::entrypoint::createCredentialsLocalStore)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.CredentialsConfiguration
		format     repository.Formater
		res        *credentialslocalstore.LocalStore
		err        error
	}{
		{
			desc:       "Testing error creating credentials local storage on get credentials entrypoint when configuration is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create credentials local store in the entrypoint, a file system is required"),
		},
		{
			desc: "Testing error creating credentials local storage on get credentials entrypoint when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			err: errors.New(errContext, "To create credentials local store in the entrypoint, credentials configuration is required"),
		},
		{
			desc: "Testing error creating credentials local storage on get credentials entrypoint when compatibilitier format is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			conf: &configuration.CredentialsConfiguration{
				Format: "json",
			},
			err: errors.New(errContext, "To create credentials local store in the entrypoint, compatibilitier is required"),
		},
		{
			desc: "Testing error creating credentials local storage on get credentials entrypoint when local storage path is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
				Format:      "json",
			},
			err: errors.New(errContext, "To create credentials local store in the entrypoint, local storage path is required"),
		},
		{
			desc: "Testing error creating credentials filter on get credentials entrypoint when storage type in not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				Format:      "json",
				StorageType: "unknown",
			},
			err: errors.New(errContext, "Unsupported credentials storage type 'unknown'"),
		},
		{
			desc: "Testing create credentials local storage on get credentials",
			entrypoint: NewEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType:      credentials.LocalStore,
				LocalStoragePath: "./test/credentials",
				Format:           "json",
			},
			res: &credentialslocalstore.LocalStore{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createCredentialsLocalStore(test.conf)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, store)
			}

		})
	}
}

func TestCredentialsFilter(t *testing.T) {

	errContext := "(get::credentials::entrypoint::createCredentialsFilter)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		res        repository.CredentialsFilterer
		err        error
	}{
		{
			desc:       "Testing error creating credentials filter on get credentials when configuration is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create the credentials filter in the entrypoint, configuration is required"),
		},
		{
			desc:       "Testing error creating credentials filter on get credentials when credentials configuration is not defined",
			entrypoint: NewEntrypoint(),
			conf:       &configuration.Configuration{},
			err:        errors.New(errContext, "To create the credentials filter in the entrypoint, credentials configuration is required"),
		},
		{
			desc: "Testing create credentials filter on get credentials",
			entrypoint: NewEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      credentials.LocalStore,
					LocalStoragePath: "./test/credentials",
					Format:           credentials.JSONFormat,
				},
			},
			res: &credentialslocalstore.LocalStore{},
		},
		{
			desc:       "Testing create credentials filter on get credentials",
			entrypoint: NewEntrypoint(),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType: credentials.EnvvarsStore,
					Format:      credentials.JSONFormat,
				},
			},
			res: &credentialsenvvarsstore.EnvvarsStore{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			store, err := test.entrypoint.createCredentialsFilter(test.conf)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, store)
			}

		})
	}
}
