package credentials

const (
	// EnvvarsStore is a store backend that uses environment variables to store credentials
	EnvvarsStore = "envvars"
	// LocalStore is a store backend which stores credentials in local file system
	LocalStore = "local"
	// MockStore is a mocked backend store
	MockStore = "mock"
)
