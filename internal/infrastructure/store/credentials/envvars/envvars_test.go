package envvars

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	credentialsjsonformater "github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/json"
	credentialsformater "github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/mock"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
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
			err:   errors.New(errContext, "Envvars credentials store requires a console writer to store a badge"),
		},
		{
			desc: "Testing error storing envvars credentials when env vars formater is not provided",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
			),
			err: errors.New(errContext, "Envvars credentials store requires a formater to store a badge"),
		},
		{
			desc: "Testing error storing envvars credentials when encryption is not provided",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
				WithFormater(credentialsformater.NewMockFormater()),
			),
			err: errors.New(errContext, "Envvars credentials store requires encryption to store a badge"),
		},
		{
			desc: "Testing error storing envvars credentials when ID is not provided",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
				WithFormater(credentialsformater.NewMockFormater()),
				WithEncryption(encryption.NewMockEncryption()),
			),
			err: errors.New(errContext, "To store credentials badege, is required an ID"),
		},
		{
			desc: "Testing store envvars credentials",
			store: NewEnvvarsStore(
				WithConsole(console.NewMockConsole()),
				WithFormater(credentialsformater.NewMockFormater()),
				WithEncryption(encryption.NewMockEncryption()),
			),
			id: "myregistry.test:5000",
			badge: &credentials.Badge{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(s *EnvvarsStore) {

				s.formater.(*credentialsformater.MockFormater).On("Marshal",
					&credentials.Badge{
						ID:       "myregistry.test:5000",
						Username: "username",
						Password: "password",
					},
				).Return(`{
  "ID": "myregistry.test:5000",
  "aws_access_key_id": "",
  "aws_region": "",
  "aws_role_arn": "",
  "aws_secret_access_key": "",
  "aws_profile": "",
  "aws_shared_credentials_files": null,
  "aws_shared_config_files": null,
  "aws_use_default_credentials_chain": false,
  "docker_login_password": "",
  "docker_login_username": "",
  "password": "password",
  "username": "username",
  "private_key_file": "",
  "private_key_password": "",
  "git_ssh_user": "",
  "use_ssh_agent": false
}`, nil)

				s.encyption.(*encryption.MockEncription).On("Encrypt", `{
  "ID": "myregistry.test:5000",
  "aws_access_key_id": "",
  "aws_region": "",
  "aws_role_arn": "",
  "aws_secret_access_key": "",
  "aws_profile": "",
  "aws_shared_credentials_files": null,
  "aws_shared_config_files": null,
  "aws_use_default_credentials_chain": false,
  "docker_login_password": "",
  "docker_login_username": "",
  "password": "password",
  "username": "username",
  "private_key_file": "",
  "private_key_password": "",
  "git_ssh_user": "",
  "use_ssh_agent": false
}`).Return("encrypted-text", nil)

				s.console.(*console.MockConsole).On("Warn", []interface{}{"You must create the following environment variable to use the recently created credentials:"})
				s.console.(*console.MockConsole).On("Warn", []interface{}{" STEVEDORE_ENVVARS_CREDENTIALS_E3A70918293EEFC49419599C9D8B5ABC=encrypted-text"})
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
				if test.store.encyption != nil {
					test.store.encyption.(*encryption.MockEncription).AssertExpectations(t)
				}
				if test.store.console != nil {
					test.store.console.(*console.MockConsole).AssertExpectations(t)
				}

			}
		})
	}
}

func TestGet(t *testing.T) {
	errContextGet := "(store::credentials::envvars::Get)"
	errContextPrivGet := "(store::credentials::envvars::get)"

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
			desc:  "Testing error getting envvars credentials when ID is not provided",
			store: NewEnvvarsStore(),
			err:   errors.New(errContextGet, "To get credentials badge, is required an ID"),
		},
		{
			desc:  "Testing error getting envvars credentials when backend is not provided",
			store: NewEnvvarsStore(),
			id:    "myregistry.test:5000",
			err: errors.New(errContextGet, "Error getting credentials badge 'myregistry.test:5000'",
				errors.New(errContextPrivGet, "Envvars credentials store requires a backend to get credentials badge")),
		},
		{
			desc: "Testing error getting envvars credentials when env vars formater is not provided",
			store: NewEnvvarsStore(
				WithBackend(backend.NewMockEnvvarsBackend()),
			),
			id: "myregistry.test:5000",
			err: errors.New(errContextGet, "Error getting credentials badge 'myregistry.test:5000'",
				errors.New(errContextPrivGet, "Envvars credentials store requires a formater to get credentials badge")),
		},
		{
			desc: "Testing error getting envvars credentials when encryption is not provided",
			store: NewEnvvarsStore(
				WithBackend(backend.NewMockEnvvarsBackend()),
				WithFormater(credentialsformater.NewMockFormater()),
			),
			id: "myregistry.test:5000",
			err: errors.New(errContextGet, "Error getting credentials badge 'myregistry.test:5000'",
				errors.New(errContextPrivGet, "Envvars credentials store requires encryption to get credentials badge")),
		},
		{
			desc: "Testing get envvars credentials badge",
			store: NewEnvvarsStore(
				WithBackend(backend.NewMockEnvvarsBackend()),
				WithFormater(credentialsjsonformater.NewJSONFormater()),
				WithEncryption(encryption.NewEncryption(
					encryption.WithKey("encryption-key"),
				)),
			),
			id:  "myregistry.test:5000",
			err: &errors.Error{},
			prepareAssertFunc: func(s *EnvvarsStore) {
				s.backend.(*backend.MockEnvvarsBackend).On("Getenv", "STEVEDORE_ENVVARS_CREDENTIALS_E3A70918293EEFC49419599C9D8B5ABC").Return("b39h6nnNl5ALc8zqq809KDa7dFFSqaYV7FvMm+tyX/KZrt7XwdWwZiQJIsyBaiqvZPgT3ljkKGpIo9/2Xdb/goUT+w4p9ug3i7rJEeX5m0JFXe6uHuxfMvPD2yXcyXaNqCfPUjKwzd/XUbJqRaJw/STwcMy4AN2vhxCiFkhWQaIl9vDREmIlJa2PZP9igIUqU2Pc0GENylj5VNCzqDjspSlrEaiRMq2WV4dSGwD1NI5azS7Ok9ITsqhjkpsRKUwiRYGBPLStKYb3lv6San2FRMZE8MaG0rm3W/71cwSTvI5Jzt8qrgqyCC87cib61jt8F7GsldtMSY+ZUTd52ryTq5hrntq5xJ8qLo+/xebO/6Td5Mv9qL6kEL3Zqo1SUBIkufS3Wo/okEuTmgl/U9nVfWIA0OgCg/fDWQVw+zYM7fNkDliWk63ypzQFIQyeFmOVBlBeFGXzShkJcBw1/AABgylOBQEdw3pzGuODUuXEFQ8EdiEFLe7EEQ9CA6E9yuJXInDgiBbgZBD9H67riCPagKIYYaDmfC7XD/fCVdzmwiSJpiQ2ySf5iU7tqGzeO7qIqkyZXiJes9xyVHCBqjOdZMwMydnjz6qIi2aFXYvXI5Sy4VqRtpL3AfJ9+c4SQZR//wspuod0h3Ek9dV+bg1voqNnMwpvAOnUMoGB7yhrmEkjUfyKpJo7koAl7FBakDlZ36UAHT+jDXIWzQ==")
			},
			res: &credentials.Badge{
				ID:                            "myregistry.test:5000",
				AllowUseSSHAgent:              false,
				AWSAccessKeyID:                "",
				AWSProfile:                    "",
				AWSRegion:                     "",
				AWSRoleARN:                    "",
				AWSSecretAccessKey:            "",
				AWSSharedConfigFiles:          []string{""},
				AWSSharedCredentialsFiles:     []string{""},
				AWSUseDefaultCredentialsChain: false,
				GitSSHUser:                    "",
				Password:                      "password",
				PrivateKeyFile:                "",
				PrivateKeyPassword:            "",
				Username:                      "username",
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

// func TestGenerateResult(t *testing.T) {

// 	store := NewEnvvarsStore(
// 		WithConsole(console.NewConsole(os.Stdout, nil)),
// 		WithFormater(credentialsjsonformater.NewJSONFormater()),
// 		WithEncryption(encryption.NewEncryption(
// 			encryption.WithKey("encryption-key"),
// 		)),
// 	)

// 	id := "myregistry.test:5000"
// 	badge := &credentials.Badge{
// 		AllowUseSSHAgent:              false,
// 		AWSAccessKeyID:                "",
// 		AWSProfile:                    "",
// 		AWSRegion:                     "",
// 		AWSRoleARN:                    "",
// 		AWSSecretAccessKey:            "",
// 		AWSSharedConfigFiles:          []string{""},
// 		AWSSharedCredentialsFiles:     []string{""},
// 		AWSUseDefaultCredentialsChain: false,
// 		GitSSHUser:                    "",
// 		Password:                      "password",
// 		PrivateKeyFile:                "",
// 		PrivateKeyPassword:            "",
// 		Username:                      "username",
// 	}
// 	err := store.Store(id, badge)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

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

				// s.backend.(*backend.MockEnvvarsBackend).On("Environ").Return(
				// 	[]string{
				// 		"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY1_TEST_5000_ATTR_USERNAME=username",
				// 		"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY1_TEST_5000_ATTR_PASSWORD=password",
				// 		"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY2_TEST_5000_ATTR_USERNAME=username",
				// 		"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY2_TEST_5000_ATTR_PASSWORD=password",
				// 		"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY3_TEST_5000_ATTR_USERNAME=username",
				// 		"STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY3_TEST_5000_ATTR_PASSWORD=password",
				// 	},
				// )

				// 		s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY1_TEST_5000_ATTR").Return(
				// 			&credentials.Badge{
				// 				ID:                            "myregistry1_test_5000",
				// 				AllowUseSSHAgent:              false,
				// 				AWSAccessKeyID:                "",
				// 				AWSProfile:                    "",
				// 				AWSRegion:                     "",
				// 				AWSRoleARN:                    "",
				// 				AWSSecretAccessKey:            "",
				// 				AWSSharedConfigFiles:          []string{},
				// 				AWSSharedCredentialsFiles:     []string{},
				// 				AWSUseDefaultCredentialsChain: false,
				// 				GitSSHUser:                    "",
				// 				Password:                      "password",
				// 				PrivateKeyFile:                "",
				// 				PrivateKeyPassword:            "",
				// 				Username:                      "username",
				// 			}, nil)

				// 		s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY2_TEST_5000_ATTR").Return(
				// 			&credentials.Badge{
				// 				ID:                            "myregistry2_test_5000",
				// 				AllowUseSSHAgent:              false,
				// 				AWSAccessKeyID:                "",
				// 				AWSProfile:                    "",
				// 				AWSRegion:                     "",
				// 				AWSRoleARN:                    "",
				// 				AWSSecretAccessKey:            "",
				// 				AWSSharedConfigFiles:          []string{},
				// 				AWSSharedCredentialsFiles:     []string{},
				// 				AWSUseDefaultCredentialsChain: false,
				// 				GitSSHUser:                    "",
				// 				Password:                      "password",
				// 				PrivateKeyFile:                "",
				// 				PrivateKeyPassword:            "",
				// 				Username:                      "username",
				// 			}, nil)

				// 		s.backend.(*backend.MockEnvvarsBackend).On("AchieveBadge", "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY3_TEST_5000_ATTR").Return(
				// 			&credentials.Badge{
				// 				ID:                            "myregistry3_test_5000",
				// 				AllowUseSSHAgent:              false,
				// 				AWSAccessKeyID:                "",
				// 				AWSProfile:                    "",
				// 				AWSRegion:                     "",
				// 				AWSRoleARN:                    "",
				// 				AWSSecretAccessKey:            "",
				// 				AWSSharedConfigFiles:          []string{},
				// 				AWSSharedCredentialsFiles:     []string{},
				// 				AWSUseDefaultCredentialsChain: false,
				// 				GitSSHUser:                    "",
				// 				Password:                      "password",
				// 				PrivateKeyFile:                "",
				// 				PrivateKeyPassword:            "",
				// 				Username:                      "username",
				// 			}, nil)

				// 	},
				// 	res: []*credentials.Badge{
				// 		{
				// 			ID:                            "myregistry1_test_5000",
				// 			AllowUseSSHAgent:              false,
				// 			AWSAccessKeyID:                "",
				// 			AWSProfile:                    "",
				// 			AWSRegion:                     "",
				// 			AWSRoleARN:                    "",
				// 			AWSSecretAccessKey:            "",
				// 			AWSSharedConfigFiles:          []string{},
				// 			AWSSharedCredentialsFiles:     []string{},
				// 			AWSUseDefaultCredentialsChain: false,
				// 			GitSSHUser:                    "",
				// 			Password:                      "password",
				// 			PrivateKeyFile:                "",
				// 			PrivateKeyPassword:            "",
				// 			Username:                      "username",
				// 		},
				// 		{
				// 			ID:                            "myregistry2_test_5000",
				// 			AllowUseSSHAgent:              false,
				// 			AWSAccessKeyID:                "",
				// 			AWSProfile:                    "",
				// 			AWSRegion:                     "",
				// 			AWSRoleARN:                    "",
				// 			AWSSecretAccessKey:            "",
				// 			AWSSharedConfigFiles:          []string{},
				// 			AWSSharedCredentialsFiles:     []string{},
				// 			AWSUseDefaultCredentialsChain: false,
				// 			GitSSHUser:                    "",
				// 			Password:                      "password",
				// 			PrivateKeyFile:                "",
				// 			PrivateKeyPassword:            "",
				// 			Username:                      "username",
				// 		},
				// 		{
				// 			ID:                            "myregistry3_test_5000",
				// 			AllowUseSSHAgent:              false,
				// 			AWSAccessKeyID:                "",
				// 			AWSProfile:                    "",
				// 			AWSRegion:                     "",
				// 			AWSRoleARN:                    "",
				// 			AWSSecretAccessKey:            "",
				// 			AWSSharedConfigFiles:          []string{},
				// 			AWSSharedCredentialsFiles:     []string{},
				// 			AWSUseDefaultCredentialsChain: false,
				// 			GitSSHUser:                    "",
				// 			Password:                      "password",
				// 			PrivateKeyFile:                "",
				// 			PrivateKeyPassword:            "",
				// 			Username:                      "username",
				// 		},
			},
		},
	}

	for _, test := range tests {
		if test.prepareAssertFunc != nil && test.store != nil {
			test.prepareAssertFunc(test.store)
		}

		res := test.store.All()
		assert.ElementsMatch(t, test.res, res)

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
