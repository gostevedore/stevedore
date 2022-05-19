package command

import (
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/images/image"
)

// BuildCommandFactory is a factory to create a build command
type BuildCommandFactory struct{}

// NewBuildCommandFactory creates a new build command factory
func NewBuildCommandFactory() *BuildCommandFactory {
	return &BuildCommandFactory{}
}

// New returns a new build command constructor
func (f *BuildCommandFactory) New(driver driver.BuildDriverer, image *image.Image, options *driver.BuildDriverOptions) BuildCommander {
	return NewBuildCommand(driver, image, options)
}
