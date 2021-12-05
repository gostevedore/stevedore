package gitcontext

import (
	gitauth "github.com/apenella/go-docker-builder/pkg/auth/git"
	"github.com/stretchr/testify/mock"
)

// MockGitBuildContext defines a build context from a git repository
type MockGitBuildContext struct {
	mock.Mock
}

// NewMockGitBuildContext creates a new git build context
func NewMockGitBuildContext() *MockGitBuildContext {
	return &MockGitBuildContext{}
}

// WithRepository sets the repository
func (c *MockGitBuildContext) WithRepository(repository string) {
	c.Called(repository)
}

// WithReference sets the reference
func (c *MockGitBuildContext) WithReference(reference string) {
	c.Called(reference)
}

// WithPath sets the path inside the repository where is located the context
func (c *MockGitBuildContext) WithPath(path string) {
	c.Called(path)
}

// WithAuth sets the authentication
func (c *MockGitBuildContext) WithAuth(auth gitauth.GitAuther) {
	c.Called(auth)
}
