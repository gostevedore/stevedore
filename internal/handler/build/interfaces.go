package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
	"github.com/gostevedore/stevedore/internal/service/build"
	"github.com/gostevedore/stevedore/internal/service/build/plan"
)

// PlanFactorier interface defines the execution plan
type PlanFactorier interface {
	NewPlan(id string, parameters map[string]interface{}) (plan.Planner, error)
}

// ServiceBuilder is the service for build commands
type ServiceBuilder interface {
	Build(ctx context.Context, buildPlan build.Planner, name string, version []string, options *build.ServiceOptions, optionsFunc ...build.OptionsFunc) error
}

// Dispatcher is a dispatcher for build commands
type Dispatcher interface {
	Start(ctx context.Context, opts ...dispatch.OptionsFunc) error
	Enqueue(schedule.Jobber)
}

// DriverFactorier interface defines the factory to create a build driver
type DriverFactorier interface {
	Get(id string) (driver.BuildDriverer, error)
	Register(id string, driver driver.BuildDriverer) error
}
