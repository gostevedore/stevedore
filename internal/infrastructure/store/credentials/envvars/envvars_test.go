package envvars

import (
	"os"
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

				s.encryption.(*encryption.MockEncription).On("Encrypt", `{
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
				if test.store.encryption != nil {
					test.store.encryption.(*encryption.MockEncription).AssertExpectations(t)
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
				s.backend.(*backend.MockEnvvarsBackend).On("Getenv", "STEVEDORE_ENVVARS_CREDENTIALS_E3A70918293EEFC49419599C9D8B5ABC").Return("3e2e6af012fa7712fcfa268363801aa4bba063ca4f8ed59dcdb30fd2da215575a22b84dbf24c17e6872ec04ac8b48e40da7e97267fcc761dd0993e483b1e6de0d2967a8e96c3054ec2f06c0126d43b6029067a1cf51ea870fc92746ddeb4eaa5df556263cf67af89a2e6fb45218d7619eff5a2abcd29856e1fd79972a0d54b9eaa085df2e8d49de2b147c7ed11c5130f6d988b7fc2b43652733c691b12e96e7715e797eb96cf45bedd48da2b7ea22ca50d18f86bc7d318fcb76924e94b88f539b3896356e71c9e0f2a8cb83fc34d26159c2f8dac2eef044be8f3ee369b41f04c7c5a0e433f0839be052ce2fa3af9b917ac7c34ea56722ca16e66c3e4dea43389db6b15d0ec10b9218365a202f7083e65ec7c150f0736b52c52accf58bfff905575f53318594b7bc0558b7f330ca7bdf306e042f15166955e9b2fd77ed981a913a6e0fa111986705856d6fccedb693a08cee5e8cfa85a45d2c702a2389b6e3ae0f88884e835bf709c327bbb0d60f43caa07ab3df576227d2797b2a97d6b55d774a8d9c3002c884b12bfe0ecb6fc6e8912b7cb886c9a85cfb319a1d8204aa091ca449f8c3bcda7d4ff4e26ca6f73633d9ff5c2ceaaff77c8700816856434dd7bf6c057cda7fb6bb0014dcba110ab7887f1edaa15f4fdd4383b5a53d1635fb50e7fb5ec9cd0a7c4162a6c18ffe8ebef5aa7df93bbc329801f46fc044a94e3d5203a796079bab07c665252d61419a34c994db38373f0175ff57261f1")
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

func TestAll(t *testing.T) {
	tests := []struct {
		desc              string
		store             *EnvvarsStore
		prepareAssertFunc func(*EnvvarsStore)
		res               []*credentials.Badge
		err               error
	}{
		{
			desc: "Testing achieving all badges from envvars store",
			store: NewEnvvarsStore(
				WithFormater(credentialsjsonformater.NewJSONFormater()),
				WithBackend(backend.NewMockEnvvarsBackend()),
				WithEncryption(encryption.NewEncryption(
					encryption.WithKey("encryption-key"),
				)),
			),
			prepareAssertFunc: func(s *EnvvarsStore) {
				s.backend.(*backend.MockEnvvarsBackend).On("Environ").Return(
					[]string{
						"STEVEDORE_ENVVARS_CREDENTIALS_E3A70918293EEFC49419599C9D8B5ABC=3e2e6af012fa7712fcfa268363801aa4bba063ca4f8ed59dcdb30fd2da215575a22b84dbf24c17e6872ec04ac8b48e40da7e97267fcc761dd0993e483b1e6de0d2967a8e96c3054ec2f06c0126d43b6029067a1cf51ea870fc92746ddeb4eaa5df556263cf67af89a2e6fb45218d7619eff5a2abcd29856e1fd79972a0d54b9eaa085df2e8d49de2b147c7ed11c5130f6d988b7fc2b43652733c691b12e96e7715e797eb96cf45bedd48da2b7ea22ca50d18f86bc7d318fcb76924e94b88f539b3896356e71c9e0f2a8cb83fc34d26159c2f8dac2eef044be8f3ee369b41f04c7c5a0e433f0839be052ce2fa3af9b917ac7c34ea56722ca16e66c3e4dea43389db6b15d0ec10b9218365a202f7083e65ec7c150f0736b52c52accf58bfff905575f53318594b7bc0558b7f330ca7bdf306e042f15166955e9b2fd77ed981a913a6e0fa111986705856d6fccedb693a08cee5e8cfa85a45d2c702a2389b6e3ae0f88884e835bf709c327bbb0d60f43caa07ab3df576227d2797b2a97d6b55d774a8d9c3002c884b12bfe0ecb6fc6e8912b7cb886c9a85cfb319a1d8204aa091ca449f8c3bcda7d4ff4e26ca6f73633d9ff5c2ceaaff77c8700816856434dd7bf6c057cda7fb6bb0014dcba110ab7887f1edaa15f4fdd4383b5a53d1635fb50e7fb5ec9cd0a7c4162a6c18ffe8ebef5aa7df93bbc329801f46fc044a94e3d5203a796079bab07c665252d61419a34c994db38373f0175ff57261f1",
						"STEVEDORE_ENVVARS_CREDENTIALS_33A9F775E178D5617883676370230761=9515a6797bf56797cd88d985499e9039e822d5c7b143a54f53f8ee70592313cb08d76def6761b8af92e1ed4efe682eeeeb36fa924ce8bc39166b71fb7e8f254ece752640ffaa9ee5870dc85a1e3a642ce53d1064e35c35dfd2119b2da8a6e5f16f0f4f6ad930d143fae7c4bea654d9b0a13d7f2d53d6aaac2a3aa636262644bbd51991566a4b4aa3ed6a7325bdc7158b4c061393838cdf91d8be11aaa57ae96b68e3d01e1a015fe96d2d3eb80e6294c999fbd11e22009ac13678460aaa3a9987926215f2ba53dbdebd38894bd394ac256fdd23b4397e21e93b4cbc7f9b8e5f23fa2f23746b52fc92115f3b71954d900e7d6501e33f26014868455020442bc07fc9b40af63e1a303cc922f5e511d380fdd599280a30131af14e3d230e644f32984f2b6a845749ec86d2bbc3891cf6cd91d3ac83653a35d39587ff297429649782912b6312ce57b2427c5feb8fc51df429084aea443a49cb3d66242d1daca2d283ed569155dbff522f4d0171fc7b2807d39166c99eeb115ce7e296b804eb950ebe172145e75cea52e109e3686ef67f8043b6e23c47158a9ea614a35afe96b34ed03363aae9540cf5b87050913283e03c6f78b09cb3d41d3112cc95a73eb691ddd4e7f7157438232ddfb352d33ba11531d271482f24022fbad0fe30fb51f9c3d2609af55b640739e282ea1c02399f4343ce010741461fa93320628eb917e0d5065601f9e3e95518c23419daa4e0b733aac7a71057cf49db571a403862",
						"STEVEDORE_ENVVARS_CREDENTIALS_37E49154AAF0C364891E66923C4B5B7D=9312dca35d21dfdc97a6acb1109aeb419ac099571db56210c245cd112f8bdfc7d065f7dda53834b5b474093e2b589878a09c58d2dbbae8eda6cb6882e5b9d2c189663407f03327d3404d81ca8b1ea61042db15b859119811d833e75c5dc98d7644df5b73f65e95230d9e8a102c9d5d05284ed3ffca175271a843f3c820e1cc585c2c8829607968869fdbb4c7b299fdf365c8295699d0b23fe5826a60079ed23e3914ed02fe80e149b7310aff0c27d53edb5fb71d74e4c3ceba0286d31d4e19b705c1dad938c1e912ee5befd2a22b1ec2e9c96553150f0ee5a87c62548bc8941fd5e420a8374c361dbd3d4b5d2707f408b0c44f9e56eb9927d46122a517970b124c91cc63482cb637ed6dc6b58a493973aa13290c918fe0cee10a2f528c1891ad47c5a28a37f7fc1157b9bc602bdc146cd9160a6edee07106ae4f20734d8e27a2e13bf3c434199e7c3f67296fcf7eef37147ff36ebe8039076a0e4f2fe9ccc49079a41a60dba744216e89b8b777c06c080ba7726e236c09511e456be0093fc0cb8e71f32153d7bb79c5e0aee96e9921bba47f7de7bd5d7e217d3cf5c822deef29d7da1cbe9b174e9e371fff8c373e14ffc6f8582e043f203a3cf732dd3edc7860edff436a8af22f8af9718f2a7eae2f1b1c17263920c79cb7f8ea7f533ed1ac2523b395342adc8005910567a6355fc4a11f0322ffcfbfa90a3189b1e64761513bcd88c9b25f6399e8c25d5df6df1c24f8e31efec5b1ec3286c8a894",
					},
				)

				s.backend.(*backend.MockEnvvarsBackend).On("Getenv", "STEVEDORE_ENVVARS_CREDENTIALS_E3A70918293EEFC49419599C9D8B5ABC").Return("3e2e6af012fa7712fcfa268363801aa4bba063ca4f8ed59dcdb30fd2da215575a22b84dbf24c17e6872ec04ac8b48e40da7e97267fcc761dd0993e483b1e6de0d2967a8e96c3054ec2f06c0126d43b6029067a1cf51ea870fc92746ddeb4eaa5df556263cf67af89a2e6fb45218d7619eff5a2abcd29856e1fd79972a0d54b9eaa085df2e8d49de2b147c7ed11c5130f6d988b7fc2b43652733c691b12e96e7715e797eb96cf45bedd48da2b7ea22ca50d18f86bc7d318fcb76924e94b88f539b3896356e71c9e0f2a8cb83fc34d26159c2f8dac2eef044be8f3ee369b41f04c7c5a0e433f0839be052ce2fa3af9b917ac7c34ea56722ca16e66c3e4dea43389db6b15d0ec10b9218365a202f7083e65ec7c150f0736b52c52accf58bfff905575f53318594b7bc0558b7f330ca7bdf306e042f15166955e9b2fd77ed981a913a6e0fa111986705856d6fccedb693a08cee5e8cfa85a45d2c702a2389b6e3ae0f88884e835bf709c327bbb0d60f43caa07ab3df576227d2797b2a97d6b55d774a8d9c3002c884b12bfe0ecb6fc6e8912b7cb886c9a85cfb319a1d8204aa091ca449f8c3bcda7d4ff4e26ca6f73633d9ff5c2ceaaff77c8700816856434dd7bf6c057cda7fb6bb0014dcba110ab7887f1edaa15f4fdd4383b5a53d1635fb50e7fb5ec9cd0a7c4162a6c18ffe8ebef5aa7df93bbc329801f46fc044a94e3d5203a796079bab07c665252d61419a34c994db38373f0175ff57261f1", nil)
				s.backend.(*backend.MockEnvvarsBackend).On("Getenv", "STEVEDORE_ENVVARS_CREDENTIALS_33A9F775E178D5617883676370230761").Return("9515a6797bf56797cd88d985499e9039e822d5c7b143a54f53f8ee70592313cb08d76def6761b8af92e1ed4efe682eeeeb36fa924ce8bc39166b71fb7e8f254ece752640ffaa9ee5870dc85a1e3a642ce53d1064e35c35dfd2119b2da8a6e5f16f0f4f6ad930d143fae7c4bea654d9b0a13d7f2d53d6aaac2a3aa636262644bbd51991566a4b4aa3ed6a7325bdc7158b4c061393838cdf91d8be11aaa57ae96b68e3d01e1a015fe96d2d3eb80e6294c999fbd11e22009ac13678460aaa3a9987926215f2ba53dbdebd38894bd394ac256fdd23b4397e21e93b4cbc7f9b8e5f23fa2f23746b52fc92115f3b71954d900e7d6501e33f26014868455020442bc07fc9b40af63e1a303cc922f5e511d380fdd599280a30131af14e3d230e644f32984f2b6a845749ec86d2bbc3891cf6cd91d3ac83653a35d39587ff297429649782912b6312ce57b2427c5feb8fc51df429084aea443a49cb3d66242d1daca2d283ed569155dbff522f4d0171fc7b2807d39166c99eeb115ce7e296b804eb950ebe172145e75cea52e109e3686ef67f8043b6e23c47158a9ea614a35afe96b34ed03363aae9540cf5b87050913283e03c6f78b09cb3d41d3112cc95a73eb691ddd4e7f7157438232ddfb352d33ba11531d271482f24022fbad0fe30fb51f9c3d2609af55b640739e282ea1c02399f4343ce010741461fa93320628eb917e0d5065601f9e3e95518c23419daa4e0b733aac7a71057cf49db571a403862", nil)
				s.backend.(*backend.MockEnvvarsBackend).On("Getenv", "STEVEDORE_ENVVARS_CREDENTIALS_37E49154AAF0C364891E66923C4B5B7D").Return("9312dca35d21dfdc97a6acb1109aeb419ac099571db56210c245cd112f8bdfc7d065f7dda53834b5b474093e2b589878a09c58d2dbbae8eda6cb6882e5b9d2c189663407f03327d3404d81ca8b1ea61042db15b859119811d833e75c5dc98d7644df5b73f65e95230d9e8a102c9d5d05284ed3ffca175271a843f3c820e1cc585c2c8829607968869fdbb4c7b299fdf365c8295699d0b23fe5826a60079ed23e3914ed02fe80e149b7310aff0c27d53edb5fb71d74e4c3ceba0286d31d4e19b705c1dad938c1e912ee5befd2a22b1ec2e9c96553150f0ee5a87c62548bc8941fd5e420a8374c361dbd3d4b5d2707f408b0c44f9e56eb9927d46122a517970b124c91cc63482cb637ed6dc6b58a493973aa13290c918fe0cee10a2f528c1891ad47c5a28a37f7fc1157b9bc602bdc146cd9160a6edee07106ae4f20734d8e27a2e13bf3c434199e7c3f67296fcf7eef37147ff36ebe8039076a0e4f2fe9ccc49079a41a60dba744216e89b8b777c06c080ba7726e236c09511e456be0093fc0cb8e71f32153d7bb79c5e0aee96e9921bba47f7de7bd5d7e217d3cf5c822deef29d7da1cbe9b174e9e371fff8c373e14ffc6f8582e043f203a3cf732dd3edc7860edff436a8af22f8af9718f2a7eae2f1b1c17263920c79cb7f8ea7f533ed1ac2523b395342adc8005910567a6355fc4a11f0322ffcfbfa90a3189b1e64761513bcd88c9b25f6399e8c25d5df6df1c24f8e31efec5b1ec3286c8a894", nil)

			},
			err: &errors.Error{},
			res: []*credentials.Badge{
				{
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
				{
					ID:                            "myregistry1.test:5000",
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
				{
					ID:                            "myregistry2.test:5000",
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
		},
	}

	for _, test := range tests {
		if test.prepareAssertFunc != nil && test.store != nil {
			test.prepareAssertFunc(test.store)
		}

		res, err := test.store.All()
		if err != nil {
			assert.Equal(t, test.err, err)
		}
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

func TestGenerateResult(t *testing.T) {

	store := NewEnvvarsStore(
		WithConsole(console.NewConsole(os.Stdout, nil)),
		WithFormater(credentialsjsonformater.NewJSONFormater()),
		WithEncryption(encryption.NewEncryption(
			encryption.WithKey("encryption-key"),
		)),
	)

	id := "id"
	badge := &credentials.Badge{
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
	}
	err := store.Store(id, badge)
	if err != nil {
		t.Error(err)
	}
}
