package credentials

// createCredentialsFlagOptions is the options for the create credentials command
type createCredentialsFlagOptions struct {
	AllowUseSSHAgent              bool
	AskAWSSecretAccessKey         bool
	AskPassword                   bool
	AWSAccessKeyID                string
	AWSProfile                    string
	AWSRegion                     string
	AWSRoleARN                    string
	AWSSharedConfigFiles          []string
	AWSSharedCredentialsFiles     []string
	AWSUseDefaultCredentialsChain bool
	GitSSHUser                    string
	ID                            string
	PrivateKeyFile                string
	PrivateKeyPassword            string
	Username                      string

	DEPRECATEDRegistryHost                 string
	DEPRECATEDDockerRegistryCredentialsDir string
}
