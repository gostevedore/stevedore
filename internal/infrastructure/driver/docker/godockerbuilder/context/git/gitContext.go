package gitcontext

import (
	gitauth "github.com/apenella/go-docker-builder/pkg/auth/git"
	context "github.com/apenella/go-docker-builder/pkg/build/context/git"
)

// GitBuildContext defines a build context from a git repository
type GitBuildContext struct {
	// context contains build context
	context.GitBuildContext
}

// NewGitBuildContext creates a new git build context
func NewGitBuildContext() *GitBuildContext {
	return &GitBuildContext{}
}

// WithRepository sets the repository
func (c *GitBuildContext) WithRepository(repository string) {
	c.Repository = repository
}

// WithReference sets the reference
func (c *GitBuildContext) WithReference(reference string) {
	c.Reference = reference
}

// WithPath sets the path inside the repository where is located the context
func (c *GitBuildContext) WithPath(path string) {
	c.Path = path
}

// WithAuth sets the authentication
func (c *GitBuildContext) WithAuth(auth gitauth.GitAuther) {
	c.Auth = auth
}
