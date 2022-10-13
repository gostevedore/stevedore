package credentials

import (
	"fmt"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/create/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/create/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	credentialscompatibilitiy "github.com/gostevedore/stevedore/internal/infrastructure/credentials/compatibility"
	credentialsformat "github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/mock"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestPrepareCredentialsId(t *testing.T) {
	errContext := "(create::credentials::entrypoint:::prepareCredentialsId)"

	tests := []struct {
		desc              string
		entrypoint        *CreateCredentialsEntrypoint
		args              []string
		options           *Options
		res               string
		prepareAssertFunc func(*CreateCredentialsEntrypoint)
		err               error
	}{
		{
			desc:       "Testing error on prepare credentials id into create credentials when args are nil",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "To execute the create credentials entrypoint, an argument with credential id is required"),
		},
		{
			desc:       "Testing prepare credentials id into create credentials",
			entrypoint: NewCreateCredentialsEntrypoint(),
			args:       []string{"id"},
			res:        "id",
		},
		{
			desc: "Testing prepare credentials id into create credentials",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			args: []string{"id", "ignored_id"},
			res:  "id",
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("Warn", []interface{}{"Ignoring extra arguments: [ignored_id]\n"})
			},
		},
		{
			desc: "Testing prepare credentials id into create credentials using deprecated registry host",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			args: []string{},
			options: &Options{
				DEPRECATEDRegistryHost: "registry-host",
			},
			res: "registry-host",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.entrypoint)
			}

			res, err := test.entrypoint.prepareCredentialsId(test.args, test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestGetPassword(t *testing.T) {
	errContext := "(create::credentials::entrypoint::getPassword)"

	tests := []struct {
		desc              string
		entrypoint        *CreateCredentialsEntrypoint
		options           *Options
		prepareAssertFunc func(*CreateCredentialsEntrypoint)
		res               string
		err               error
	}{
		{
			desc:       "Testing error on create credentials entrypoint get password method when console is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Console must be provided to execute create credentials entrypoint"),
		},
		{
			desc: "Testing create credentials entrypoint get password",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getPasswordInputMessage).Return("p4ssw0rd", nil)
				e.console.(*console.MockConsole).On("Write", []byte(fmt.Sprintln())).Return(0, nil)
			},
			res: "p4ssw0rd",
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.entrypoint)
			}

			res, err := test.entrypoint.getPassword()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestGetAWSSecretAccessKey(t *testing.T) {
	errContext := "(create::credentials::entrypoint::getAWSSecretAccessKey)"

	tests := []struct {
		desc              string
		entrypoint        *CreateCredentialsEntrypoint
		prepareAssertFunc func(*CreateCredentialsEntrypoint)
		res               string
		err               error
	}{
		{
			desc:       "Testing error on create credentials entrypoint get aws secret access key method when console is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Console must be provided to execute create credentials entrypoint"),
		},
		{
			desc: "Testing create credentials entrypoint get aws secret access key",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getAWSSecretAccessKeyInputMessage).Return("s3cret", nil)
				e.console.(*console.MockConsole).On("Write", []byte(fmt.Sprintln())).Return(0, nil)
			},
			res: "s3cret",
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.entrypoint)
			}

			res, err := test.entrypoint.getAWSSecretAccessKey()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestPrepareHandlerOptions(t *testing.T) {
	errContext := "(create::credentials::entrypoint::prepareHandlerOptions)"

	tests := []struct {
		desc              string
		entrypoint        *CreateCredentialsEntrypoint
		handlerOptions    *handler.Options
		prepareAssertFunc func(*CreateCredentialsEntrypoint)
		res               *handler.Options
		err               error
	}{
		{
			desc:       "Testing error on create credentials entrypoint prepare handler options method when handler options are not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Handler options must be provided to execute create credentials entrypoint"),
		},
		{
			desc: "Testing create credentials entrypoint prepare handler options method when ask for password is enabled",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			handlerOptions: &handler.Options{
				Username: "username",
			},
			res: &handler.Options{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getPasswordInputMessage).Return("password", nil)
				e.console.(*console.MockConsole).On("Write", []byte(fmt.Sprintln())).Return(0, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing create credentials entrypoint prepare handler options method when ask for aws secret access key is enabled",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			handlerOptions: &handler.Options{
				AWSAccessKeyID:            "AWSAccessKeyID",
				AWSSharedConfigFiles:      []string{"file"},
				AWSSharedCredentialsFiles: []string{"file"},
			},
			res: &handler.Options{
				AWSAccessKeyID:            "AWSAccessKeyID",
				AWSSecretAccessKey:        "AWSSecretAccessKey",
				AWSSharedConfigFiles:      []string{"file"},
				AWSSharedCredentialsFiles: []string{"file"},
			},
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getAWSSecretAccessKeyInputMessage).Return("AWSSecretAccessKey", nil)
				e.console.(*console.MockConsole).On("Write", []byte(fmt.Sprintln())).Return(0, nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.entrypoint)
			}

			res, err := test.entrypoint.prepareHandlerOptions(test.handlerOptions)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestPrepareConfiguration(t *testing.T) {
	errContext := "(create::credentials::entrypoint::prepareConfiguration)"

	tests := []struct {
		desc          string
		entrypoint    *CreateCredentialsEntrypoint
		options       *Options
		configuration *configuration.Configuration
		err           error
		res           *configuration.Configuration
	}{
		{
			desc:       "Testing error preparing configuration on create credentials entrypoint when options are not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Entrypoint options must be provided to prepare configuration"),
		},
		{
			desc:       "Testing error preparing configuration on create credentials entrypoint when configuration is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			options:    &Options{},
			err:        errors.New(errContext, "Configuration must be provided to prepare configuration"),
		},
		{
			desc:          "Testing error preparing configuration on create credentials entrypoint when configuration credentials are not provided",
			entrypoint:    NewCreateCredentialsEntrypoint(),
			options:       &Options{},
			configuration: &configuration.Configuration{},
			err:           errors.New(errContext, "Configuration credentials must be provided to prepare configuration"),
		},
		{
			desc:       "Testing error preparing configuration on create credentials entrypoint when credentials storage type is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			options:    &Options{},
			configuration: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{},
			},
			res: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					LocalStoragePath: "path",
				},
			},
			err: errors.New(errContext, "Credentials storage type must be provided to prepare configuration"),
		},
		{
			desc:       "Testing prepare configuration on create credentials entrypoint",
			entrypoint: NewCreateCredentialsEntrypoint(),
			options: &Options{
				LocalStoragePath: "path",
			},
			configuration: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType: credentials.LocalStore,
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			conf, err := test.entrypoint.prepareConfiguration(test.configuration, test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.configuration, conf)
			}
		})
	}
}

func TestCreateCredentialsStore(t *testing.T) {

	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	tests := []struct {
		desc       string
		entrypoint *CreateCredentialsEntrypoint
		conf       *configuration.Configuration
		options    *Options
		err        error
		res        application.CredentialsStorer
	}{
		{
			desc:       "Testing error creating credentials store on create credentials entrypoint when compatibilitier is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "To create the credentials store, compatibilitier is required"),
		},
		// {
		// 	desc: "Testing error creating credentials store on create credentials entrypoint when configuration is not provided",
		// 	entrypoint: NewCreateCredentialsEntrypoint(
		// 		WithCompatibilitier(compatibility.NewMockCompatibility()),
		// 	),
		// 	conf: nil,
		// 	err:  errors.New(errContext, "To create the credentials store, configuration is required"),
		// },
		{
			desc: "Testing error creating credentials store on create credentials entrypoint when credentials configuration is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithCompatibilitier(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create the credentials store, credentials configuration is required"),
		},
		{
			desc: "Testing error creating credentials store on create credentials entrypoint when credentials format is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithCompatibilitier(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{},
			},
			err: errors.New(errContext, "To create the credentials store, credentials format must be defined"),
		},
		{
			desc: "Testing error creating credentials store on create credentials entrypoint when credentials storage type is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithCompatibilitier(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					Format: credentials.JSONFormat,
				},
			},
			err: errors.New(errContext, "To create the credentials store, credentials storage type must be defined"),
		},

		{
			desc: "Testing error creating credentials store on create credentials entrypoint when options are not provided",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithCompatibilitier(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					Format:      credentials.JSONFormat,
					StorageType: credentials.LocalStore,
				},
			},
			err: errors.New(errContext, "To create the credentials store, options are required"),
		},
		{
			desc: "Testing create a local credentials store on create credentials entrypoint",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithCompatibilitier(compatibility.NewMockCompatibility()),
				WithFileSystem(afero.NewMemMapFs()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					Format:           credentials.JSONFormat,
					StorageType:      credentials.LocalStore,
					LocalStoragePath: "path",
				},
			},
			options: &Options{
				ForceCreate: true,
			},
			err: &errors.Error{},
			res: local.NewLocalStore(
				afero.NewMemMapFs(),
				"path",
				credentialsformat.NewMockFormater(),
				credentialscompatibilitiy.NewCredentialsCompatibility(compatibility.NewMockCompatibility()),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.createCredentialsStore(test.conf, test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}

}

func TestCreateCredentialsLocalStore(t *testing.T) {
	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	tests := []struct {
		desc            string
		entrypoint      *CreateCredentialsEntrypoint
		compatibilitier credentialslocalstore.CredentialsCompatibilier
		conf            *configuration.CredentialsConfiguration
		format          repository.Formater
		res             *local.LocalStore
		err             error
	}{
		{
			desc:       "Testing error creating credentials local store on create credentials entrypoint when compatibilitier is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "To create the credentials local store, credentials compatibilitier is required"),
		},
		{
			desc:            "Testing error creating credentials local store on create credentials entrypoint when credentials configuration is not provided",
			entrypoint:      NewCreateCredentialsEntrypoint(),
			compatibilitier: credentialscompatibilitiy.NewCredentialsCompatibility(compatibility.NewMockCompatibility()),
			err:             errors.New(errContext, "To create the credentials local store, credentials configuration is required"),
		},
		{
			desc:            "Testing error creating credentials local store on create credentials entrypoint when local storage path is not provided",
			entrypoint:      NewCreateCredentialsEntrypoint(),
			conf:            &configuration.CredentialsConfiguration{},
			compatibilitier: credentialscompatibilitiy.NewCredentialsCompatibility(compatibility.NewMockCompatibility()),
			err:             errors.New(errContext, "To create the credentials local store, local storage path is required"),
		},
		{
			desc:            "Testing error creating credentials local store on create credentials entrypoint when credentials formater is not provided",
			entrypoint:      NewCreateCredentialsEntrypoint(),
			compatibilitier: credentialscompatibilitiy.NewCredentialsCompatibility(compatibility.NewMockCompatibility()),
			conf: &configuration.CredentialsConfiguration{
				LocalStoragePath: "path",
			},
			err: errors.New(errContext, "To create the credentials local store, formater is required"),
		},
		{
			desc:            "Testing error creating credentials local store on create credentials entrypoint when filesystem is not provided",
			entrypoint:      NewCreateCredentialsEntrypoint(),
			compatibilitier: credentialscompatibilitiy.NewCredentialsCompatibility(compatibility.NewMockCompatibility()),
			conf: &configuration.CredentialsConfiguration{
				LocalStoragePath: "path",
			},
			format: credentialsformat.NewMockFormater(),
			err:    errors.New(errContext, "To create the credentials local store, filesystem is required"),
		},
		{
			desc: "Testing create credentials local store on create credentials entrypoint",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithFileSystem(afero.NewMemMapFs()),
			),
			compatibilitier: credentialscompatibilitiy.NewCredentialsCompatibility(compatibility.NewMockCompatibility()),
			conf: &configuration.CredentialsConfiguration{
				LocalStoragePath: "path",
			},
			format: credentialsformat.NewMockFormater(),
			err:    errors.New(errContext, "To create the credentials local store, filesystem is required"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.createCredentialsLocalStore(test.compatibilitier, test.conf, test.format)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}

}
