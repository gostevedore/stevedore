package build

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/application/build"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {

	errContext := "(handler::build::Handler)"

	tests := []struct {
		desc              string
		handler           *Handler
		imageName         string
		options           *Options
		err               error
		prepareAssertFunc func(string, PlanFactorier, BuildApplication)
		assertFunc        func(PlanFactorier, BuildApplication)
	}{
		{
			desc: "Testing error when plan factory is not defined",
			handler: &Handler{
				planFactory: nil,
			},
			err: errors.New(errContext, "Build handler requires a plan factory"),
		},
		{
			desc: "Testing error when build service is not defined",
			handler: &Handler{
				planFactory: plan.NewMockPlanFactory(),
				app:         nil,
			},
			err: errors.New(errContext, "Build handler requires a service to build images"),
		},
		{
			desc: "Testing error when received label format is not valid",
			handler: &Handler{
				planFactory: plan.NewMockPlanFactory(),
				app:         build.NewMockApplication(),
			},
			options: &Options{
				Labels: []string{"invalid_label"},
			},
			err: errors.New(errContext, "Invalid label format 'invalid_label'"),
		},
		{
			desc: "Testing error when received persistent variable format is not valid",
			handler: &Handler{
				planFactory: plan.NewMockPlanFactory(),
				app:         build.NewMockApplication(),
			},
			options: &Options{
				PersistentVars: []string{"invalid_persistent_var"},
			},
			err: errors.New(errContext, "Invalid persistent variable format 'invalid_persistent_var'"),
		},
		{
			desc: "Testing handler build with all options",
			handler: &Handler{
				planFactory: plan.NewMockPlanFactory(),
				app:         build.NewMockApplication(),
			},
			options: &Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "ansible-intermediate-container",
				AnsibleInventoryPath:             "ansible-inventory",
				AnsibleLimit:                     "ansible-limit",
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label-1=value-label1"},
				PersistentLabels:                 []string{"plabel-1=pvalue-label1"},
				PersistentVars:                   []string{"persistent-var-1=value-persistent-var1"},
				PullParentImage:                  true,
				PushImagesAfterBuild:             true,
				RemoveImagesAfterPush:            true,
				SemanticVersionTagsTemplates:     []string{"{{ .Major }}.{{ .Minor }}"},
				Versions:                         []string{"version-1", "version-2"},
				Vars:                             []string{"var-1=value-var1"},

				BuildOnCascade: false,
				CascadeDepth:   5,
			},
			err: &errors.Error{},
			prepareAssertFunc: func(name string, p PlanFactorier, s BuildApplication) {
				p.(*plan.MockPlanFactory).On(
					"NewPlan",
					"single",
					map[string]interface{}{},
				).Return(plan.NewMockPlan(), nil)

				s.(*build.MockApplication).On(
					"Build",
					context.TODO(),
					plan.NewMockPlan(), name,
					[]string{"version-1", "version-2"},
					&build.Options{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "ansible-intermediate-container",
						AnsibleInventoryPath:             "ansible-inventory",
						AnsibleLimit:                     "ansible-limit",
						EnableSemanticVersionTags:        true,
						ImageFromName:                    "image-from-name",
						ImageFromRegistryHost:            "image-from-registry-host",
						ImageFromRegistryNamespace:       "image-from-registry-namespace",
						ImageFromVersion:                 "image-from-version",
						ImageName:                        "image-name",
						ImageRegistryHost:                "image-registry-host",
						ImageRegistryNamespace:           "image-registry-namespace",
						ImageVersions:                    []string{"version-1", "version-2"},
						Labels:                           map[string]string{"label-1": "value-label1"},
						PersistentLabels:                 map[string]string{"plabel-1": "pvalue-label1"},
						PersistentVars:                   map[string]interface{}{"persistent-var-1": "value-persistent-var1"},
						PullParentImage:                  true,
						PushImageAfterBuild:              true,
						RemoveImagesAfterPush:            true,
						SemanticVersionTagsTemplates:     []string{"{{ .Major }}.{{ .Minor }}"},
						Vars:                             map[string]interface{}{"var-1": "value-var1"},
					},
					mock.AnythingOfType("[]build.OptionsFunc"),
				).Return(nil)
			},
			assertFunc: func(p PlanFactorier, s BuildApplication) {
				s.(*build.MockApplication).AssertExpectations(t)
				p.(*plan.MockPlanFactory).AssertExpectations(t)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.imageName, test.handler.planFactory, test.handler.app)
			}

			err := test.handler.Handler(context.TODO(), test.imageName, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(test.handler.planFactory, test.handler.app)
			}
		})
	}
}

func TestCreateBuildPlan(t *testing.T) {
	errContext := "(handler::build::createBuildPlan)"

	tests := []struct {
		desc              string
		handler           *Handler
		options           *Options
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
			handler: NewHandler(plan.NewMockPlanFactory(), build.NewMockApplication()),
			options: &Options{
				BuildOnCascade: true,
				CascadeDepth:   5,
				ImageName:      image.UndefinedStringValue,
				ImageFromName:  image.UndefinedStringValue,
			},
			res: nil,
			err: nil,
			prepareAssertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).On("NewPlan", "cascade", map[string]interface{}{
					"depth": 5,
				}).Return(plan.NewMockPlan(), nil)
			},
			assertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).AssertExpectations(t)
			},
		},
		{
			desc:    "Testing get default (single) plan",
			handler: NewHandler(plan.NewMockPlanFactory(), build.NewMockApplication()),
			options: &Options{},
			res:     nil,
			err:     nil,
			prepareAssertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).On("NewPlan", "single", map[string]interface{}{}).Return(plan.NewMockPlan(), nil)
			},
			assertFunc: func(p PlanFactorier) {
				p.(*plan.MockPlanFactory).AssertExpectations(t)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.handler.planFactory)
			}

			_, err := test.handler.createBuildPlan(test.options)
			if test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(test.handler.planFactory)
			}
		})
	}
}

func TestValidateCascadePlanOptions(t *testing.T) {
	errContext := "(handler::build::validateCascadePlanOptions)"

	tests := []struct {
		desc    string
		options *Options
		err     error
	}{
		{
			desc: "Testing not valid cascade plan options when ansible intermediate container name is defined",
			options: &Options{
				AnsibleIntermediateContainerName: "name",
			},
			err: errors.New(errContext, "Cascade plan does not support intermediate containers name, it could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when ansible inventory path is defined",
			options: &Options{
				AnsibleInventoryPath: "path",
			},
			err: errors.New(errContext, "Cascade plan does not support ansible inventory path, it could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when ansible limit is defined",
			options: &Options{
				AnsibleLimit: "limit",
			},
			err: errors.New(errContext, "Cascade plan does not support ansible limit, it could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when options are nil",
			err:  errors.New(errContext, "Options to be validated are required"),
		},
		{
			desc: "Testing not valid cascade plan options when image name is defined",
			options: &Options{
				ImageName: "name",
			},
			err: errors.New(errContext, "Cascade plan does not support image name, it could cause an unpredictable result"),
		},
		{
			desc: "Testing not valid cascade plan options when image from name is defined",
			options: &Options{
				ImageName:     image.UndefinedStringValue,
				ImageFromName: "name",
			},
			err: errors.New(errContext, "Cascade plan does not support image from name, it could cause an unpredictable result"),
		},
		{
			desc: "Testing valid options for cascade plan",
			options: &Options{
				ImageName:     image.UndefinedStringValue,
				ImageFromName: image.UndefinedStringValue,
			},
			err: &errors.Error{},
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
