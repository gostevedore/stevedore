package credentials

type Options struct {
	// // AskPassword is true if the password should be asked
	// AskPassword bool
	// // AskAWSSecretAccessKey is true if the AWS secret access key should be asked
	// AskAWSSecretAccessKey bool
	// LocalStoragePath is the location of local storage
	LocalStoragePath string
	// DEPRECATEDRegistryHost is the registry host used as credentials id
	DEPRECATEDRegistryHost string
	// ForceCreate forces to create a credential
	ForceCreate bool
}
