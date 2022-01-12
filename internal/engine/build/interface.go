package build

import (
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/engine/build/command"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/schedule/job"
)

// BuildDriverer interface defines which methods are used to build a docker image
// type BuildDriverer interface {
// 	Build(context.Context, *driver.BuildDriverOptions) error
// }

// // Imager interface defines the image
// type Imager interface {
// 	Parent() Imager
// 	Children() []Imager
// }

// // ImagesStorer interfaces defines the storage of images
// type ImagesStorer interface {
// 	Find(string, string) (*image.Image, error)
// }

// Steper interface defines the step plan
type Steper interface {
	Image() *image.Image
	Notify()
	Wait()
}

// Planner interface defines the execution plan
type Planner interface {
	Plan(string, []string) ([]*plan.Step, error)
}

// BuildersStorer interface defines the storage of builders
type BuildersStorer interface {
	Find(name string) (*builder.Builder, error)
}

// BuildCommandFactorier interface defines the factory of build commands
type BuildCommandFactorier interface {
	New(driver.BuildDriverer, *driver.BuildDriverOptions) command.BuildCommander
}

// BuildCommander interface defines the command to build a docker image
// type BuildCommander interface {
// 	Execute(context.Context) error
// }

// JobFactorier interface defines the factory of build jobs
type JobFactorier interface {
	New(job.Commander) schedule.Jobber
}

// Dispatcher is a dispatcher to build docker images
type Dispatcher interface {
	Enqueue(schedule.Jobber)
}

// DriverFactorier interface defines the factory to create a build driver
type DriverFactorier interface {
	Get(id string) (driver.BuildDriverer, error)
}

// Semverser
type Semverser interface {
	GenerateSemverList(version []string, tmpls []string) ([]string, error)
}

// CredentialsStorer
type CredentialsStorer interface {
	GetCredentials(registy string) (*credentials.RegistryUserPassAuth, error)
}
