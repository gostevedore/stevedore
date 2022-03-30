package plan

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
)

const (
	// CascadePlanID is the id for the cascade plan
	CascadePlanID = "cascade"
	// SinglePlanID is the id for the single plan
	SinglePlanID = "single"
)

// PlanFactory is a factory to create Planner
type PlanFactory struct{}

// NewPlanFactory creates a new PlanFactory
func NewPlanFactory() *PlanFactory {
	return &PlanFactory{}
}

// NewPlan creates a new Planner
func (f *PlanFactory) NewPlan(id string, parameters map[string]interface{}) (Planner, error) {
	var store ImagesStorer
	var exists bool
	var depth int

	errContext := "(PlanFactory::NewPlan)"

	switch id {
	case CascadePlanID:
		store, exists = parameters["store"].(ImagesStorer)
		if !exists || store == nil {
			return nil, errors.New(errContext, "To create a cascade plan, is required a store")
		}

		depth, exists = parameters["depth"].(int)
		if !exists {
			return nil, errors.New(errContext, "To create a cascade plan, is required a depth")
		}

		return NewCascadePlan(store, depth), nil

	case SinglePlanID:
		store, exists = parameters["store"].(ImagesStorer)
		if !exists || store == nil {
			return nil, errors.New(errContext, "To create a single plan, is required a depth")
		}

		return NewSinglePlan(store), nil
	default:
		return nil, errors.New(errContext, fmt.Sprintf("Plan '%s' has not been registered", id))
	}
}
