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
type PlanFactory struct {
	imagesStore ImagesStorer
}

// NewPlanFactory creates a new PlanFactory
func NewPlanFactory(store ImagesStorer) *PlanFactory {
	return &PlanFactory{
		imagesStore: store,
	}
}

// NewPlan creates a new Planner
func (f *PlanFactory) NewPlan(id string, parameters map[string]interface{}) (Planner, error) {
	var exists bool
	var depth int

	errContext := "(PlanFactory::NewPlan)"

	if f.imagesStore == nil {
		return nil, errors.New(errContext, "To create a build plan, is required a store")
	}

	switch id {
	case CascadePlanID:

		depth, exists = parameters["depth"].(int)
		if !exists {
			return nil, errors.New(errContext, "To create a cascade plan, is required a depth")
		}

		return NewCascadePlan(f.imagesStore, depth), nil

	case SinglePlanID:
		return NewSinglePlan(f.imagesStore), nil
	default:
		return nil, errors.New(errContext, fmt.Sprintf("Plan '%s' has not been registered", id))
	}
}
