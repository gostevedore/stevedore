package build

import (
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"github.com/gostevedore/stevedore/internal/service/build/command"
	"github.com/gostevedore/stevedore/internal/service/build/plan"
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

// BuildersStorer interface defines the storage of builders
type BuildersStorer interface {
	Find(name string) (*builder.Builder, error)
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

// // DriverFactorier interface defines the factory to create a build driver
// type DriverFactorier interface {
// 	Get(id string) (repository.BuildDriverer, error)
// }

// Semverser
type Semverser interface {
	GenerateSemverList(version []string, tmpls []string) ([]string, error)
}

// CredentialsStorer
type CredentialsStorer interface {
	Get(id string) (*credentials.UserPasswordAuth, error)
}
