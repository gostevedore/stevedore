package plan

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/store"
	"github.com/stretchr/testify/assert"
)

func TestNewPlan(t *testing.T) {
	errContext := "(PlanFactory::NewPlan)"

	tests := []struct {
		desc       string
		factory    *PlanFactory
		id         string
		parameters map[string]interface{}
		res        Planner
		err        error
	}{
		{
			desc:    "Testing new plan error when unknown id",
			factory: NewPlanFactory(),
			id:      "unknown",
			err:     errors.New(errContext, "Plan 'unknown' has not been registered"),
		},
		{
			desc:    "Testing new plan error when store is not provided on cascade plan",
			factory: NewPlanFactory(),
			id:      "cascade",
			err:     errors.New(errContext, "To create a cascade plan, is required a store"),
		},
		{
			desc:    "Testing new plan error when depth is not provided on cascade plan",
			factory: NewPlanFactory(),
			id:      "cascade",
			parameters: map[string]interface{}{
				"store": store.NewMockImageStore(),
			},
			err: errors.New(errContext, "To create a cascade plan, is required a depth"),
		},
		{
			desc:    "Testing new plan error when store is not provided on single plan",
			factory: NewPlanFactory(),
			id:      "single",
			err:     errors.New(errContext, "To create a single plan, is required a depth"),
		},
		{
			desc:    "Testing new plan that returns a cascade plan",
			factory: NewPlanFactory(),
			id:      "cascade",
			parameters: map[string]interface{}{
				"store": store.NewMockImageStore(),
				"depth": -1,
			},
			res: &CascadePlan{},
			err: &errors.Error{},
		},
		{
			desc:    "Testing new plan that returns a single plan",
			factory: NewPlanFactory(),
			id:      "single",
			parameters: map[string]interface{}{
				"store": store.NewMockImageStore(),
			},
			res: &SinglePlan{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			plan, err := test.factory.NewPlan(test.id, test.parameters)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, test.res, plan)
			}
		})
	}
}
