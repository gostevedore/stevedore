package credentials

import (
	errors "github.com/apenella/go-common-utils/error"
)

// Badge containes information to access to a credentials
type Badge struct {
	// Badge id
	ID string
	// AWSAccessKeyID is the access key ID for the AWS account
	AWSAccessKeyID string `json:"aws_access_key_id" yaml:"aws_access_key_id"`
	// AWSRegion is the region for the AWS account
	AWSRegion string `json:"aws_region" yaml:"aws_region"`
	// AWSRoleARN defines the ARN of the role to assume. It is expected to be used when AWSUseDefaultCredentialsChain is true
	AWSRoleARN string `json:"aws_role_arn" yaml:"aws_role_arn"`
	// AWSSecretAccessKey is the secret access key for the AWS account
	AWSSecretAccessKey string `json:"aws_secret_access_key" yaml:"aws_secret_access_key"`
	// AWSProfile is the name of the AWS profile to use
	AWSProfile string `json:"aws_profile" yaml:"aws_profile"`
	// AWSSharedCredentialsFiles is a list of shared credentials files to use. It is used when AWSUseDefaultCredentialsChain is true
	AWSSharedCredentialsFiles []string `json:"aws_shared_credentials_files" yaml:"aws_shared_credentials_files"`
	// AWSSharedConfigFiles is a list of shared config files to use. It is used when AWSUseDefaultCredentialsChain is true
	AWSSharedConfigFiles []string `json:"aws_shared_config_files" yaml:"aws_shared_config_files"`
	// AWSUseDefaultCredentialsChain must be set to true when you want to use the sdk default's credentials chain described at https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials
	AWSUseDefaultCredentialsChain bool `json:"aws_use_default_credentials_chain" yaml:"aws_use_default_credentials_chain"`
	// DEPRECATEDPassword password for basic auth method
	DEPRECATEDPassword string `json:"docker_login_password" yaml:"docker_login_password"`
	// DEPRECATEDUsername username for basic auth method
	DEPRECATEDUsername string `json:"docker_login_username" yaml:"docker_login_username"`
	// Password for basic auth method. It could be used to authenticate to either docker registry or git server
	Password string `json:"password" yaml:"password"`
	// Username for basic auth method. It could be used to authenticate to either docker registry or git server
	Username string `json:"username" yaml:"username"`
	// PrivateKeyFile is the path to the private key file. It could be used to authenticate to git server
	PrivateKeyFile string `json:"private_key_file" yaml:"private_key_file"`
	// PrivateKeyPassword is the password for the private key file. It could be used to authenticate to git server
	PrivateKeyPassword string `json:"private_key_password" yaml:"private_key_password"`
	// GitSSHUser is the username for the git ssh. It could be used to authenticate to git server
	GitSSHUser string `json:"git_ssh_user" yaml:"git_ssh_user"`
	// AllowUseSSHAgent must be set to true when you allow to use the ssh-agent to authenticate to the git server
	AllowUseSSHAgent bool `json:"use_ssh_agent" yaml:"use_ssh_agent"`
}

// IsValid return whether a badge is valid, otherwise returns an error with the invalid reason
func (badge *Badge) IsValid() (bool, error) {

	errContext := "(core::domain::credentials::IsValid)"

	if badge == nil {
		return false, errors.New(errContext, "Invalid badge. Badge is nil")
	}

	if badge.Username != "" && badge.Password != "" {
		return true, nil
	}

	if badge.AWSAccessKeyID != "" && badge.AWSSecretAccessKey != "" {
		return true, nil
	}

	if badge.AWSUseDefaultCredentialsChain {
		return true, nil
	}

	if badge.PrivateKeyFile != "" {
		return true, nil
	}

	if badge.AllowUseSSHAgent {
		return true, nil
	}

	return false, errors.New(errContext, "Invalid badge. Unknown reason")
}
