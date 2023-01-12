package envvars

import (
	"io"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars/backend"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	errContext := "(store::credentials::envvars::Store)"

	tests := []struct {
		desc              string
		id                string
		store             *EnvvarsStore
		badge             *credentials.Badge
		prepareAssertFunc func(*EnvvarsStore)
		err               error
	}{
		{
			desc:  "Testing error storing envvars credentials when console is not provided",
			store: NewEnvvarsStore(),
			err:   errors.New(errContext, "Envvars credentials store requires a console writer"),
		},
		{
			desc: "Testing error storing envvars credentials when env vars backend is not provided",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
			),
			err: errors.New(errContext, "Envvars credentials store requires a backend to store envvars"),
		},
		{
			desc: "Testing store credentials to envvars credentials store",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
				WithBackend(backend.NewMockEnvvarsBackend()),
			),
			id: "myregistry.test:5000",
			badge: &credentials.Badge{
				Username:                      "username",
				Password:                      "password",
				AWSUseDefaultCredentialsChain: true,
				AWSSharedCredentialsFiles:     []string{"file1", "file2"},
			},
			prepareAssertFunc: func(s *EnvvarsStore) {
				s.backend.(*backend.MockEnvvarsBackend).On("Setenv", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_AWS_SHARED_CREDENTIALS_FILES", "file1,file2")
				s.backend.(*backend.MockEnvvarsBackend).On("Setenv", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_AWS_USE_DEFAULT_CREDENTIALS_CHAIN", "1")
				s.backend.(*backend.MockEnvvarsBackend).On("Setenv", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_PASSWORD", "password")
				s.backend.(*backend.MockEnvvarsBackend).On("Setenv", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_USERNAME", "username")

				s.console.(*console.MockConsole).On("Info", []interface{}{"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_AWS_SHARED_CREDENTIALS_FILES=file1,file2"})
				s.console.(*console.MockConsole).On("Info", []interface{}{"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_AWS_USE_DEFAULT_CREDENTIALS_CHAIN=1"})
				s.console.(*console.MockConsole).On("Info", []interface{}{"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_PASSWORD=password"})
				s.console.(*console.MockConsole).On("Info", []interface{}{"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_USERNAME=username"})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil && test.store != nil {
				test.prepareAssertFunc(test.store)
			}

			err := test.store.Store(test.id, test.badge)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				if test.store.backend != nil {
					test.store.backend.(*backend.MockEnvvarsBackend).AssertExpectations(t)
				}
				if test.store.console != nil {
					test.store.console.(*console.MockConsole).AssertExpectations(t)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		desc              string
		store             *EnvvarsStore
		id                string
		prepareAssertFunc func(*EnvvarsStore)
		cleanupFunc       func()
		res               *credentials.Badge
		err               error
	}{
		{
			desc: "Testing get credentials from envvars credentials store",
			store: NewEnvvarsStore(
				WithConsole(console.NewConsole(io.Discard, nil)),
				WithBackend(backend.NewMockEnvvarsBackend()),
			),
			id: "myregistry.test:5000",
			res: &credentials.Badge{
				ID:                            "myregistry.test:5000",
				AllowUseSSHAgent:              true,
				AWSAccessKeyID:                "aws_access_key_id",
				AWSProfile:                    "aws_profile",
				AWSRegion:                     "aws_region",
				AWSRoleARN:                    "aws_role_arn",
				AWSSecretAccessKey:            "aws_secret_access_key",
				AWSSharedConfigFiles:          []string{"aws_shared_config_files"},
				AWSSharedCredentialsFiles:     []string{"aws_shared_credentials_files"},
				AWSUseDefaultCredentialsChain: true,
				GitSSHUser:                    "git_ssh_user",
				Password:                      "password",
				PrivateKeyFile:                "private_key_file",
				PrivateKeyPassword:            "private_key_password",
				Username:                      "username",
			},
			prepareAssertFunc: func(s *EnvvarsStore) {
				s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR").Return(
					&credentials.Badge{
						AllowUseSSHAgent:              true,
						AWSAccessKeyID:                "aws_access_key_id",
						AWSProfile:                    "aws_profile",
						AWSRegion:                     "aws_region",
						AWSRoleARN:                    "aws_role_arn",
						AWSSecretAccessKey:            "aws_secret_access_key",
						AWSSharedConfigFiles:          []string{"aws_shared_config_files"},
						AWSSharedCredentialsFiles:     []string{"aws_shared_credentials_files"},
						AWSUseDefaultCredentialsChain: true,
						GitSSHUser:                    "git_ssh_user",
						Password:                      "password",
						PrivateKeyFile:                "private_key_file",
						PrivateKeyPassword:            "private_key_password",
						Username:                      "username",
					}, nil)
			},
		},
	}

	for _, test := range tests {

		if test.prepareAssertFunc != nil && test.store != nil {
			test.prepareAssertFunc(test.store)
		}

		res, err := test.store.Get(test.id)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res, res)
		}
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		desc              string
		store             *EnvvarsStore
		prepareAssertFunc func(*EnvvarsStore)
		res               []*credentials.Badge
	}{
		{
			desc: "Testing achieving all badges from envvars store",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
				WithBackend(backend.NewMockEnvvarsBackend()),
			),
			prepareAssertFunc: func(s *EnvvarsStore) {
				s.backend.(*backend.MockEnvvarsBackend).On("Environ").Return(
					[]string{
						"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY1_TEST_5000_ATTR_USERNAME=username",
						"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY1_TEST_5000_ATTR_PASSWORD=password",
						"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY2_TEST_5000_ATTR_USERNAME=username",
						"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY2_TEST_5000_ATTR_PASSWORD=password",
						"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY3_TEST_5000_ATTR_USERNAME=username",
						"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY3_TEST_5000_ATTR_PASSWORD=password",
					},
				)

				s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY1_TEST_5000_ATTR").Return(
					&credentials.Badge{
						ID:                            "myregistry1_test_5000",
						AllowUseSSHAgent:              false,
						AWSAccessKeyID:                "",
						AWSProfile:                    "",
						AWSRegion:                     "",
						AWSRoleARN:                    "",
						AWSSecretAccessKey:            "",
						AWSSharedConfigFiles:          []string{},
						AWSSharedCredentialsFiles:     []string{},
						AWSUseDefaultCredentialsChain: false,
						GitSSHUser:                    "",
						Password:                      "password",
						PrivateKeyFile:                "",
						PrivateKeyPassword:            "",
						Username:                      "username",
					}, nil)

				s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY2_TEST_5000_ATTR").Return(
					&credentials.Badge{
						ID:                            "myregistry2_test_5000",
						AllowUseSSHAgent:              false,
						AWSAccessKeyID:                "",
						AWSProfile:                    "",
						AWSRegion:                     "",
						AWSRoleARN:                    "",
						AWSSecretAccessKey:            "",
						AWSSharedConfigFiles:          []string{},
						AWSSharedCredentialsFiles:     []string{},
						AWSUseDefaultCredentialsChain: false,
						GitSSHUser:                    "",
						Password:                      "password",
						PrivateKeyFile:                "",
						PrivateKeyPassword:            "",
						Username:                      "username",
					}, nil)

				s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY3_TEST_5000_ATTR").Return(
					&credentials.Badge{
						ID:                            "myregistry3_test_5000",
						AllowUseSSHAgent:              false,
						AWSAccessKeyID:                "",
						AWSProfile:                    "",
						AWSRegion:                     "",
						AWSRoleARN:                    "",
						AWSSecretAccessKey:            "",
						AWSSharedConfigFiles:          []string{},
						AWSSharedCredentialsFiles:     []string{},
						AWSUseDefaultCredentialsChain: false,
						GitSSHUser:                    "",
						Password:                      "password",
						PrivateKeyFile:                "",
						PrivateKeyPassword:            "",
						Username:                      "username",
					}, nil)

			},
			res: []*credentials.Badge{
				{
					ID:                            "myregistry1_test_5000",
					AllowUseSSHAgent:              false,
					AWSAccessKeyID:                "",
					AWSProfile:                    "",
					AWSRegion:                     "",
					AWSRoleARN:                    "",
					AWSSecretAccessKey:            "",
					AWSSharedConfigFiles:          []string{},
					AWSSharedCredentialsFiles:     []string{},
					AWSUseDefaultCredentialsChain: false,
					GitSSHUser:                    "",
					Password:                      "password",
					PrivateKeyFile:                "",
					PrivateKeyPassword:            "",
					Username:                      "username",
				},
				{
					ID:                            "myregistry2_test_5000",
					AllowUseSSHAgent:              false,
					AWSAccessKeyID:                "",
					AWSProfile:                    "",
					AWSRegion:                     "",
					AWSRoleARN:                    "",
					AWSSecretAccessKey:            "",
					AWSSharedConfigFiles:          []string{},
					AWSSharedCredentialsFiles:     []string{},
					AWSUseDefaultCredentialsChain: false,
					GitSSHUser:                    "",
					Password:                      "password",
					PrivateKeyFile:                "",
					PrivateKeyPassword:            "",
					Username:                      "username",
				},
				{
					ID:                            "myregistry3_test_5000",
					AllowUseSSHAgent:              false,
					AWSAccessKeyID:                "",
					AWSProfile:                    "",
					AWSRegion:                     "",
					AWSRoleARN:                    "",
					AWSSecretAccessKey:            "",
					AWSSharedConfigFiles:          []string{},
					AWSSharedCredentialsFiles:     []string{},
					AWSUseDefaultCredentialsChain: false,
					GitSSHUser:                    "",
					Password:                      "password",
					PrivateKeyFile:                "",
					PrivateKeyPassword:            "",
					Username:                      "username",
				},
			},
		},
	}

	for _, test := range tests {
		if test.prepareAssertFunc != nil && test.store != nil {
			test.prepareAssertFunc(test.store)
		}

		res := test.store.All()
		assert.Equal(t, test.res, res)

	}
}

func TestConvertFieldValueToString(t *testing.T) {
	errContext := "(store::credentials::envvars::convertFieldValueToString)"

	tests := []struct {
		desc  string
		field interface{}
		res   string
		err   error
	}{
		{
			desc:  "Testing convert string field value to string",
			field: "value",
			res:   "value",
			err:   &errors.Error{},
		},
		{
			desc:  "Testing convert []string field value to string",
			field: []string{"value1", "value2"},
			res:   "value1,value2",
			err:   &errors.Error{},
		},
		{
			desc:  "Testing convert bool field value to string",
			field: true,
			res:   "1",
			err:   &errors.Error{},
		},
		{
			desc:  "Testing error converting invalid field type",
			field: byte(0),
			err:   errors.New(errContext, "Field could not be converted to string"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := convertFieldValueToString(test.field)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}

}

func TestGenerateEnvvarKey(t *testing.T) {
	tests := []struct {
		desc  string
		items []string
		res   string
	}{
		{
			desc:  "Testing generate envvars key",
			items: []string{envvarsCredentialsPrefix, "myregistry.test:5000", envvarsCredentialsAttributePrefix, "attribute"},
			res:   "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000_ATTR_ATTRIBUTE",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := generateEnvvarKey(test.items...)
			assert.Equal(t, test.res, res)
		})
	}
}
