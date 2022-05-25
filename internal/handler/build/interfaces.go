package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/dispatch"
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
	Enqueue(scheduler.Jobber)
}
