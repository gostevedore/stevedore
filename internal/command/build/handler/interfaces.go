package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/engine/build"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
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
