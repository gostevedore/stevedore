package credentials

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	application "github.com/gostevedore/stevedore/internal/application/create/credentials"
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

// func TestExecute(t *testing.T) {
// 	tests := []struct {
// 		desc              string
// 		entrypoint        *CreateCredentialsEntrypoint
// 		entrypointOptions *Options
// 		handlerOptions    *handler.Options
// 		args              []string
// 		conf              *configuration.Configuration
// 		err               error
// 	}{}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf, test.entrypointOptions, test.handlerOptions)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			}
// 		})
// 	}
// 	assert.True(t, false)
// }

func TestPrepareCredentialsId(t *testing.T) {
	errContext := "(create::credentials::entrypoint:::prepareCredentialsId)"

	tests := []struct {
		desc              string
		entrypoint        *CreateCredentialsEntrypoint
		args              []string
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
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.entrypoint)
			}

			res, err := test.entrypoint.prepareCredentialsId(test.args)
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
			desc:       "Testing error on create credentials entrypoint get password method when options is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint"),
		},
		{
			desc:       "Testing error on create credentials entrypoint get password method when console is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			options:    &Options{},
			err:        errors.New(errContext, "Console must be provided to execute create credentials entrypoint"),
		},
		{
			desc: "Testing create credentials entrypoint get password",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			options: &Options{},
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getPasswordInputMessage).Return("p4ssw0rd", nil)
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

			res, err := test.entrypoint.getPassword(test.options)
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
		options           *Options
		prepareAssertFunc func(*CreateCredentialsEntrypoint)
		res               string
		err               error
	}{
		{
			desc:       "Testing error on create credentials entrypoint get aws secret access key method when options is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint"),
		},
		{
			desc:       "Testing error on create credentials entrypoint get aws secret access key method when console is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			options:    &Options{},
			err:        errors.New(errContext, "Console must be provided to execute create credentials entrypoint"),
		},
		{
			desc: "Testing create credentials entrypoint get aws secret access key",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			options: &Options{},
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getAWSSecretAccessKeyInputMessage).Return("s3cret", nil)
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

			res, err := test.entrypoint.getAWSSecretAccessKey(test.options)
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
		entrypointOptions *Options
		handlerOptions    *handler.Options
		prepareAssertFunc func(*CreateCredentialsEntrypoint)
		res               *handler.Options
		err               error
	}{
		{
			desc:       "Testing error on create credentials entrypoint prepare handler options method when entrypoint options are not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "Entrypoint options must be provided to execute create credentials entrypoint"),
		},
		{
			desc:              "Testing error on create credentials entrypoint prepare handler options method when handler options are not provided",
			entrypoint:        NewCreateCredentialsEntrypoint(),
			entrypointOptions: &Options{},
			err:               errors.New(errContext, "Handler options must be provided to execute create credentials entrypoint"),
		},
		{
			desc: "Testing create credentials entrypoint prepare handler options method when ask for password is enable",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			entrypointOptions: &Options{
				AskPassword: true,
			},
			handlerOptions: &handler.Options{
				Username: "username",
			},
			res: &handler.Options{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(e *CreateCredentialsEntrypoint) {
				e.console.(*console.MockConsole).On("ReadPassword", getPasswordInputMessage).Return("password", nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing create credentials entrypoint prepare handler options method when ask for aws secret access key is enable",
			entrypoint: NewCreateCredentialsEntrypoint(
				WithConsole(console.NewMockConsole()),
			),
			entrypointOptions: &Options{
				AskAWSSecretAccessKey: true,
			},
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
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.entrypoint)
			}

			res, err := test.entrypoint.prepareHandlerOptions(test.entrypointOptions, test.handlerOptions)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestPrepareConfiguration(t *testing.T) {
	tests := []struct {
		desc string
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
		})
	}

	assert.True(t, false)
}

func TestCreateCredentialsStore(t *testing.T) {

	errContext := "(create::credentials::entrypoint:::createCredentialsLocalStore)"

	tests := []struct {
		desc       string
		entrypoint *CreateCredentialsEntrypoint
		conf       *configuration.Configuration
		err        error
		res        application.CredentialsStorer
	}{
		{
			desc:       "Testing error creating credentials store on create credentials entrypoint when compatibilitier is not provided",
			entrypoint: NewCreateCredentialsEntrypoint(),
			err:        errors.New(errContext, "To create the credentials store, compatibilitier is required"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.createCredentialsStore(test.conf)
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
