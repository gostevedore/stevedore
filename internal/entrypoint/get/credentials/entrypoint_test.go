package credentials

import (
	"context"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	errContext := "(get::credentials::entrypoint::Execute)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
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
				WithWriter(ioutil.Discard),
				WithFileSystem(afero.NewMemMapFs()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      credentials.LocalStore,
					LocalStoragePath: "./test/credentials",
					Format:           "json",
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf)
			if err != nil {
				assert.Equal(t, test.err, err)
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
		res        *credentialslocalstore.LocalStore
		err        error
	}{
		{
			desc:       "Testing error creating credentials local storage on get credentials when configuration is not defined",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To create credentials local store, credentials configuration is required"),
		},
		{
			desc:       "Testing error creating credentials local storage on get credentials when credentials format is not defined",
			entrypoint: &Entrypoint{},
			conf:       &configuration.CredentialsConfiguration{},
			err:        errors.New(errContext, "To create credentials local store, credentials format must be specified"),
		},
		{
			desc:       "Testing error creating credentials local storage on get credentials when local storage path is not defined",
			entrypoint: &Entrypoint{},
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
				Format:      "json",
			},
			err: errors.New(errContext, "To create credentials local store, local storage path is required"),
		},
		{
			desc: "Testing error creating credentials filter on get credentials when storage type in not defined",
			entrypoint: &Entrypoint{
				fs: afero.NewMemMapFs(),
			},
			conf: &configuration.CredentialsConfiguration{
				Format:      "json",
				StorageType: "unknown",
			},
			err: errors.New(errContext, "Unsupported credentials storage type 'unknown'"),
		},
		{
			desc:       "Testing create credentials local storage on get credentials",
			entrypoint: &Entrypoint{},
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
		res        *credentialslocalstore.LocalStore
		err        error
	}{
		{
			desc:       "Testing error creating credentials filter on get credentials when file system is not defined",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To create the credentials filter, a file system is required"),
		},
		{
			desc: "Testing error creating credentials filter on get credentials when configuration is not defined",
			entrypoint: &Entrypoint{
				fs: afero.NewMemMapFs(),
			},
			err: errors.New(errContext, "To create the credentials filter, configuration is required"),
		},
		{
			desc: "Testing error creating credentials filter on get credentials when credentials configuration is not defined",
			entrypoint: &Entrypoint{
				fs: afero.NewMemMapFs(),
			},
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create the credentials filter, credentials configuration is required"),
		},
		{
			desc: "Testing create credentials filter on get credentials",
			entrypoint: &Entrypoint{
				fs: afero.NewMemMapFs(),
			},
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      credentials.LocalStore,
					LocalStoragePath: "./test/credentials",
					Format:           "json",
				},
			},
			res: &credentialslocalstore.LocalStore{},
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
