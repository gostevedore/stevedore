package command

import (
	"context"

	"github.com/gostevedore/stevedore/internal/driver"
)

// BuildCommand contains details to build a docker image
type BuildCommand struct {
	driver  BuildDriverer
	options *driver.BuildDriverOptions
}

// NewBuildCommand creates a command to build docker images
func NewBuildCommand(driver BuildDriverer, options *driver.BuildDriverOptions) *BuildCommand {
	return &BuildCommand{
		driver:  driver,
		options: options,
	}
}

// Execute performs the action
func (c *BuildCommand) Execute(ctx context.Context) error {
	return c.driver.Build(ctx, c.options)
}
