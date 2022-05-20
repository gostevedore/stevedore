package command

import (
	"github.com/gostevedore/stevedore/internal/core/domain/driver"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// BuildCommandFactory is a factory to create a build command
type BuildCommandFactory struct{}

// NewBuildCommandFactory creates a new build command factory
func NewBuildCommandFactory() *BuildCommandFactory {
	return &BuildCommandFactory{}
}

// New returns a new build command constructor
func (f *BuildCommandFactory) New(driver repository.BuildDriverer, image *image.Image, options *driver.BuildDriverOptions) BuildCommander {
	return NewBuildCommand(driver, image, options)
}
