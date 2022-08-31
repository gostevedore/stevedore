package credentials

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/create/credentials"
	handler "github.com/gostevedore/stevedore/internal/handler/create/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		desc            string
		config          *configuration.Configuration
		entrypoint      Entrypointer
		compatibility   Compatibilitier
		prepareMockFunc func(Entrypointer, Compatibilitier, *configuration.Configuration)
		args            []string
		err             error
	}{
		{
			desc:          "Testing run create credentials command",
			config:        &configuration.Configuration{},
			compatibility: compatibility.NewMockCompatibility(),
			entrypoint:    entrypoint.NewMockCreateCredentialsEntrypoint(),
			args: []string{
				"credential-id",
				"--allow-use-ssh-agent",
				"--aws-secret-access-key",
				"--password",
				"--aws-shared-config-files",
				"aws-shared-config-file1",
				"--aws-shared-config-files",
				"aws-shared-config-file2",
				"--aws-shared-credentials-files",
				"aws-shared-credentials-file1",
				"--aws-shared-credentials-files",
				"aws-shared-credentials-file2",
				"--aws-access-key-id",
				"aws-access-key-id",
				"--aws-profile",
				"aws-profile",
				"--aws-region",
				"aws-region",
				"--aws-role-arn",
				"aws-role-arn",
				"--git-ssh-user",
				"git-ssh-user",
				"--local-storage-path",
				"local-storage-path",
				"--private-key-file",
				"private-key-file",
				"--private-key-password",
				"private-key-password",
				"--username",
				"username",
			},
			prepareMockFunc: func(e Entrypointer, comp Compatibilitier, conf *configuration.Configuration) {
				e.(*entrypoint.MockCreateCredentialsEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{"credential-id"},
					conf,
					&entrypoint.Options{
						AskAWSSecretAccessKey: true,
						AskPassword:           true,
						LocalStoragePath:      "local-storage-path",
					},
					&handler.Options{
						AllowUseSSHAgent: true,
						AWSSharedConfigFiles: []string{
							"aws-shared-config-file1",
							"aws-shared-config-file2",
						},
						AWSSharedCredentialsFiles: []string{
							"aws-shared-credentials-file1",
							"aws-shared-credentials-file2",
						},
						AWSAccessKeyID:     "aws-access-key-id",
						AWSProfile:         "aws-profile",
						AWSRegion:          "aws-region",
						AWSRoleARN:         "aws-role-arn",
						GitSSHUser:         "git-ssh-user",
						PrivateKeyFile:     "private-key-file",
						PrivateKeyPassword: "private-key-password",
						Username:           "username",
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:          "Testing run create credentials command with deprecated flags",
			config:        &configuration.Configuration{},
			compatibility: compatibility.NewMockCompatibility(),
			entrypoint:    entrypoint.NewMockCreateCredentialsEntrypoint(),
			args: []string{
				"--credentials-dir",
				"credentials-dir",
				"--registry-host",
				"registry-host",
			},
			prepareMockFunc: func(e Entrypointer, comp Compatibilitier, conf *configuration.Configuration) {
				e.(*entrypoint.MockCreateCredentialsEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{},
					conf,
					&entrypoint.Options{
						DEPRECATEDRegistryHost: "registry-host",
						LocalStoragePath:       "credentials-dir",
					},
					&handler.Options{},
				).Return(nil)

				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageRegistryHost}).Return(nil)
				comp.(*compatibility.MockCompatibility).On("AddDeprecated", []string{DeprecatedFlagMessageDockerRegistryCredentialsDir}).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.entrypoint, test.compatibility, test.config)
			}

			cmd := NewCommand(context.TODO(), test.compatibility, test.config, test.entrypoint)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}
