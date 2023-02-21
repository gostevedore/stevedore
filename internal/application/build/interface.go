package build

import (
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	driverfactory "github.com/gostevedore/stevedore/internal/infrastructure/driver/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
)

// Planner interfaces defines the storage of images
type Planner interface {
	Plan(name string, versions []string) ([]*plan.Step, error)
}

// PlanSteper interface defines the step plan
type PlanSteper interface {
	Image() *image.Image
	Notify()
	Wait()
}

// BuildCommandFactorier interface defines the factory of build commands
type BuildCommandFactorier interface {
	New(repository.BuildDriverer, *image.Image, *image.BuildDriverOptions) command.BuildCommander
}

// JobFactorier interface defines the factory of build jobs
type JobFactorier interface {
	New(job.Commander) scheduler.Jobber
}

// Dispatcher is a dispatcher to build docker images
type Dispatcher interface {
	Enqueue(scheduler.Jobber)
}

// DriverFactorier interface defines the factory to create a build driver
type DriverFactorier interface {
	Get(id string) (driverfactory.BuildDriverFactoryFunc, error)
	Register(id string, driver driverfactory.BuildDriverFactoryFunc) error
}

// Semverser
type Semverser interface {
	GenerateSemverList(version []string, tmpls []string) ([]string, error)
}

// // CredentialsStorer
// type CredentialsStorer interface {
// 	Get(id string) (*credentials.UserPasswordAuth, error)
// }
