package command

import "github.com/gostevedore/stevedore/internal/driver"

// BuildCommandFactory is a factory to create a build command
type BuildCommandFactory struct{}

func NewBuildCommandFactory() BuildCommandFactory {
	return BuildCommandFactory{}
}

// New returns a new build command constructor
func (f *BuildCommandFactory) New(driver BuildDriverer, options *driver.BuildDriverOptions) *BuildCommand {
	return NewBuildCommand(driver, options)
}
