package credentials

// Options is the options for the create credentials handler
type Options struct {
	AllowUseSSHAgent              bool
	AWSAccessKeyID                string
	AWSProfile                    string
	AWSRegion                     string
	AWSRoleARN                    string
	AWSSecretAccessKey            string
	AWSSharedConfigFiles          []string
	AWSSharedCredentialsFiles     []string
	AWSUseDefaultCredentialsChain bool
	GitSSHUser                    string
	Password                      string
	PrivateKeyFile                string
	PrivateKeyPassword            string
	Username                      string
}
