package pathcontext

import (
	context "github.com/apenella/go-docker-builder/pkg/build/context/path"
)

// PathBuildContext defines a docker build context from local filesystem
type PathBuildContext struct {
	// context contains build context
	context.PathBuildContext
}

// NewPathBuildContext creates a new PathBuildContext
func NewPathBuildContext() *PathBuildContext {
	return &PathBuildContext{}
}

// WithPath sets the path of the build context
func (c *PathBuildContext) WithPath(path string) {
	c.Path = path
}
