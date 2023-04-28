package credentials

const (
	// StoreAuthProvider provider which reads auth directly from the credential
	StoreAuthProvider = "store"
	// AWSECRSAuthProvider provider which uses aws ecr get-login-password to get user/password auth
	AWSECRSAuthProvider = "aws-ecr"
	// MockAuthProvider is a mocked auth provider
	MockAuthProvider = "mock"
)
