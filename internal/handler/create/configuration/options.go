package configuration

// Options for create configuration handler
type Options struct {
	BuildersPath                 string
	Concurrency                  int
	CredentialsFormat            string
	CredentialsLocalStoragePath  string
	CredentialsStorageType       string
	EnableSemanticVersionTags    bool
	ImagesPath                   string
	LogPathFile                  string
	PushImages                   bool
	SemanticVersionTagsTemplates []string
}
