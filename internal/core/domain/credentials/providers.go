package credentials

const (
	// BadgeCredentialsProvider provider which reads credentials directly from badge
	BadgeCredentialsProvider = "badge"
	// AWSECRSCredentialsProvider provider which uses aws ecr get-login-password to get user/password credentials
	AWSECRSCredentialsProvider = "aws-ecr"
	// MockCredentialsProvider is a mocked credentials provider
	MockCredentialsProvider = "mock"
)
