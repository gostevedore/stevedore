package build

import (
	"context"

	"github.com/gostevedore/stevedore/internal/engine/build"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
)

// PlanFactorier interface defines the execution plan
type PlanFactorier interface {
	NewPlan(id string, parameters map[string]interface{}) (plan.Planner, error)
}

// ServiceBuilder is the service for build commands
type ServiceBuilder interface {
	Build(ctx context.Context, buildPlan build.Planner, name string, version []string, options *build.ServiceOptions) error
}
