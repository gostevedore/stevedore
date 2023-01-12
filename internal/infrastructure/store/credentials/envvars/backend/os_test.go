package backend

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/structs"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/assert"
)

func TestAchieveBadge(t *testing.T) {
	tests := []struct {
		desc              string
		backend           *OSEnvvarsBackend
		id                string
		prepareAssertFunc func(*OSEnvvarsBackend, string)
		cleanupFunc       func(*OSEnvvarsBackend, string)
		res               *credentials.Badge
		err               error
	}{
		{
			desc:    "Testing get credentials from envvars credentials store",
			id:      "STEVEDORE_ENVVARS_CREDENTIALS_MYREGISTRY_TEST_5000",
			backend: NewOSEnvvarsBackend(),
			res: &credentials.Badge{
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
				DEPRECATEDPassword:            "docker_login_password",
				DEPRECATEDUsername:            "docker_login_username",
			},
			prepareAssertFunc: func(backend *OSEnvvarsBackend, id string) {
				badge := &credentials.Badge{}
				badgeFields := structs.Fields(badge)
				for _, field := range badgeFields {

					attribute := field.Tag("mapstructure")
					if attribute == "" {
						continue
					}

					key := strings.ToUpper(
						strings.Join(
							[]string{id, attribute},
							"_",
						),
					)
					if attribute == "aws_use_default_credentials_chain" || attribute == "use_ssh_agent" {
						attribute = "1"
					}
					backend.Setenv(key, attribute)

				}
			},
			cleanupFunc: func(backend *OSEnvvarsBackend, id string) {
				badge := &credentials.Badge{}

				badgeFields := structs.Fields(badge)
				for _, field := range badgeFields {

					attribute := field.Tag("mapstructure")
					if attribute == "" {
						continue
					}
					key := strings.ToUpper(
						strings.Join(
							[]string{id, attribute},
							"_",
						),
					)
					backend.Unsetenv(key)
				}
			},
		},
	}

	for _, test := range tests {

		if test.prepareAssertFunc != nil && test.backend != nil {
			test.prepareAssertFunc(test.backend, test.id)
			defer test.cleanupFunc(test.backend, test.id)
		}

		res, err := test.backend.AchieveBadge(test.id)
		if err != nil {
			fmt.Println(">>>>>", err.Error())
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res, res)
		}
	}
}
