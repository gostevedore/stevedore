package credentials

// createCredentialsFlagOptions is the options for the create credentials command
type createCredentialsFlagOptions struct {
	// AllowUseSSHAgent
	AllowUseSSHAgent bool
	// AskPrivateKeyPassword
	AskPrivateKeyPassword bool
	// AWSAccessKeyID
	AWSAccessKeyID string
	// AWSProfile
	AWSProfile string
	// AWSRegion
	AWSRegion string
	// AWSRoleARN
	AWSRoleARN string
	// AWSSharedConfigFiles
	AWSSharedConfigFiles []string
	// AWSSharedCredentialsFiles
	AWSSharedCredentialsFiles []string
	// AWSUseDefaultCredentialsChain
	AWSUseDefaultCredentialsChain bool
	// GitSSHUser
	GitSSHUser string
	// LocalStoragePath
	LocalStoragePath string
	// PrivateKeyFile
	PrivateKeyFile string
	// Username
	Username string
	// Force
	Force bool

	DEPRECATEDRegistryHost                 string
	DEPRECATEDDockerRegistryCredentialsDir string
}
