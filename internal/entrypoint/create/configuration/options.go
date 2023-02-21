package configuration

type Options struct {
	BuildersPath                     string
	Concurrency                      int
	ConfigurationFilePath            string
	CredentialsEncryptionKey         string
	CredentialsFormat                string
	CredentialsLocalStoragePath      string
	CredentialsStorageType           string
	EnableSemanticVersionTags        bool
	Force                            bool
	GenerateCredentialsEncryptionKey bool
	ImagesPath                       string
	LogPathFile                      string
	PushImages                       bool
	SemanticVersionTagsTemplates     []string
}
