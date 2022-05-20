package command

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/driver"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// BuildCommand contains details to build a docker image
type BuildCommand struct {
	driver  repository.BuildDriverer
	image   *image.Image
	options *driver.BuildDriverOptions
}

// NewBuildCommand creates a command to build docker images
func NewBuildCommand(driver repository.BuildDriverer, i *image.Image, options *driver.BuildDriverOptions) *BuildCommand {
	return &BuildCommand{
		driver:  driver,
		image:   i,
		options: options,
	}
}

// Execute performs the action
func (c *BuildCommand) Execute(ctx context.Context) error {
	errContext := "(command::Execute)"

	if c.image == nil {
		return errors.New(errContext, "An image is required to execute a command")
	}

	if c.options == nil {
		return errors.New(errContext, "Options are required to execute a command")
	}

	return c.driver.Build(ctx, c.image, c.options)
}
