package build

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/engine/build"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {}

func TestGetPlan(t *testing.T) {
	errContext := "(build::getPlan)"

	tests := []struct {
		desc              string
		handler           *Handler
		options           *HandlerOptions
		res               plan.Planner
		err               error
		prepareAssertFunc func(PlanFactorier)
		assertFunc        func(PlanFactorier)
	}{
		{
			desc: "Testing error when plan factory is not defined",
			handler: &Handler{
				planFactory: nil,
			},
			err: errors.New(errContext, "To create a build plan, is required a plan factory"),
		},
		{
			desc: "Testing error when options are is not defined",
			handler: &Handler{
				planFactory: plan.NewMockPlanFactory(),
			},
			options: nil,
			err:     errors.New(errContext, "To create a build plan, is required a service options"),
		},
		{
			desc:    "Testing get cascade plan",
			handler: NewHandler(dispatch.NewMockDispatch(), plan.NewMockPlanFactory(), build.NewMockService()),
			options: &HandlerOptions{
				BuildOnCascade: true,
				CascadeDepth:   5,
			},
			res: nil,
			err: &errors.Error{},
			prepareAssertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).On("NewPlan", "cascade", map[string]interface{}{
					"depth": 5,
				}).Return(plan.NewMockPlan(), &errors.Error{})
			},
			assertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).AssertExpectations(t)
			},
		},
		{
			desc:    "Testing get default (single) plan",
			handler: NewHandler(dispatch.NewMockDispatch(), plan.NewMockPlanFactory(), build.NewMockService()),
			options: &HandlerOptions{},
			res:     nil,
			err:     &errors.Error{},
			prepareAssertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).On("NewPlan", "single", map[string]interface{}{}).Return(plan.NewMockPlan(), &errors.Error{})
			},
			assertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).AssertExpectations(t)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.handler.planFactory)
			}

			_, err := test.handler.getPlan(test.options)
			if test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(test.handler.planFactory)
			}
		})
	}
}

func TestValidateCascadePlanOptions(t *testing.T) {
	errContext := "(build::validateCascadePlanOptions)"

	tests := []struct {
		desc    string
		options *HandlerOptions
		err     error
	}{
		{
			desc: "Testing not valid cascade plan options when ansible intermediate container name is defined",
			options: &HandlerOptions{
				AnsibleIntermediateContainerName: "name",
			},
			err: errors.New(errContext, "Cascade plan does not support intermediate containers name. It could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when ansible inventory path is defined",
			options: &HandlerOptions{
				AnsibleInventoryPath: "path",
			},
			err: errors.New(errContext, "Cascade plan does not support ansible inventory path. It could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when ansible limit is defined",
			options: &HandlerOptions{
				AnsibleLimit: "limit",
			},
			err: errors.New(errContext, "Cascade plan does not support ansible limit. It could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when options are nil",
			err:  errors.New(errContext, "Options to be validated are required"),
		},
		{
			desc: "Testing not valid cascade plan options when image name is defined",
			options: &HandlerOptions{
				ImageName: "name",
			},
			err: errors.New(errContext, "Cascade plan does not support image name. It could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when image from name is defined",
			options: &HandlerOptions{
				ImageFromName: "name",
			},
			err: errors.New(errContext, "Cascade plan does not support image from name. It could cause an unpredictable result"),
		},
		{
			desc:    "Testing valid options for cascade plan",
			options: &HandlerOptions{},
			err:     &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := validateCascadePlanOptions(test.options)
			if err != nil {
				assert.Equal(t, err.Error(), test.err.Error())
			} else {
				assert.Empty(t, err)
			}

		})
	}
}
