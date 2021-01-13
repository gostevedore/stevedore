package gitcontext

// GitContext defines a build context from a git repository
type GitContext struct {
	// Repository which will be used as docker build context
	Repository string `yaml:"repository"`
	// Reference is the name of the branch to clone. By default is used 'master'
	Reference string `yaml:"reference"`
}
