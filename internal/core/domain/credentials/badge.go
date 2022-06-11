package credentials

// Badge containes information to access to a credentials
type Badge struct {
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSRegion          string `json:"aws_region"`
	AWSRoleARN         string `json:"aws_role_arn"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
	DEPRECATEDPassword string `json:"docker_login_password"`
	DEPRECATEDUsername string `json:"docker_login_username"`
	Password           string `json:"password"`
	Provider           string `json:"provider"`
	Username           string `json:"username"`
	PrivateKeyFile     string `json:"private_key_file"`
	PrivateKeyPassword string `json:"private_key_password"`
	GitSSHUser         string `json:"git_ssh_user"`
}
