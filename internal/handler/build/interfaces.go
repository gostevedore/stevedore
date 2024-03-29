package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/application/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/dispatch"
)

// PlanFactorier interface defines the execution plan
type PlanFactorier interface {
	NewPlan(id string, parameters map[string]interface{}) (plan.Planner, error)
}

// BuildApplication is the service for build commands
type BuildApplication interface {
	Build(ctx context.Context, buildPlan build.Planner, name string, version []string, options *build.Options, optionsFunc ...build.OptionsFunc) error
}

// Dispatcher is a dispatcher for build commands
type Dispatcher interface {
	Start(ctx context.Context, opts ...dispatch.OptionsFunc) error
	Enqueue(scheduler.Jobber)
}
