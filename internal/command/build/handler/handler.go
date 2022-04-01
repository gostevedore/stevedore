package build

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/engine/build"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
)

const AssignmentTokenSymbol = '='

// Handler is a handler for build commands
type Handler struct {
	planFactory PlanFactorier
	service     ServiceBuilder
}

// NewHandler creates a new handler for build commands
func NewHandler(p PlanFactorier, s ServiceBuilder) *Handler {
	return &Handler{
		planFactory: p,
		service:     s,
	}
}

// Handler handles build commands
func (h *Handler) Handler(ctx context.Context, options *HandlerOptions) error {

	errContext := "(build::Handler)"
	var err error
	var buildPlan build.Planner

	buildServiceOptions := &build.ServiceOptions{}

	if h.planFactory == nil {
		return errors.New(errContext, "Build handler requires a plan factory")
	}

	if h.service == nil {
		return errors.New(errContext, "Build handler requires a service to build images")
	}

	buildServiceOptions.AnsibleConnectionLocal = options.AnsibleConnectionLocal
	buildServiceOptions.AnsibleIntermediateContainerName = options.AnsibleIntermediateContainerName
	buildServiceOptions.AnsibleInventoryPath = options.AnsibleInventoryPath
	buildServiceOptions.AnsibleLimit = options.AnsibleLimit

	// concurrency
	// debug
	// dryrun

	buildServiceOptions.EnableSemanticVersionTags = options.EnableSemanticVersionTags

	buildServiceOptions.ImageFromName = options.ImageFromName
	buildServiceOptions.ImageFromRegistryHost = options.ImageFromRegistryHost
	buildServiceOptions.ImageFromRegistryNamespace = options.ImageFromRegistryNamespace
	buildServiceOptions.ImageFromVersion = options.ImageFromVersion

	// imageName
	buildServiceOptions.ImageRegistryHost = options.ImageRegistryHost
	buildServiceOptions.ImageRegistryNamespace = options.ImageRegistryNamespace
	buildServiceOptions.ImageVersions = append([]string{}, options.Versions...)

	for _, labels := range options.Labels {

		if strings.IndexRune(labels, AssignmentTokenSymbol) < 0 {
			return errors.New(errContext, fmt.Sprintf("Invalid label format '%s'", labels))
		}
		kLabel := labels[:strings.IndexRune(labels, AssignmentTokenSymbol)]
		vLabel := labels[strings.IndexRune(labels, AssignmentTokenSymbol)+1:]

		buildServiceOptions.Labels[kLabel] = vLabel
	}

	for _, persistentVars := range options.PersistentVars {

		if strings.IndexRune(persistentVars, AssignmentTokenSymbol) < 0 {
			return errors.New(errContext, fmt.Sprintf("Invalid label format '%s'", persistentVars))
		}
		kPVar := persistentVars[:strings.IndexRune(persistentVars, AssignmentTokenSymbol)]
		vPVar := persistentVars[strings.IndexRune(persistentVars, AssignmentTokenSymbol)+1:]

		buildServiceOptions.PersistentVars[kPVar] = vPVar
	}

	buildServiceOptions.PullParentImage = options.PullParentImage
	buildServiceOptions.PushImageAfterBuild = options.PushImagesAfterBuild
	buildServiceOptions.RemoveImagesAfterPush = options.RemoveImagesAfterPush

	buildPlan, err = h.getPlan(options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = h.service.Build(ctx, buildPlan, options.ImageName, options.Versions, buildServiceOptions)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

func (h *Handler) getPlan(options *HandlerOptions) (plan.Planner, error) {
	errContext := "(build::getPlan)"

	var err error
	var plan build.Planner

	planParameters := map[string]interface{}{}
	planType := "single"

	if h.planFactory == nil {
		return nil, errors.New(errContext, "To create a build plan, is required a plan factory")
	}

	if options == nil {
		return nil, errors.New(errContext, "To create a build plan, is required a service options")
	}

	if options.BuildOnCascade {
		planType = "cascade"
		planParameters["depth"] = options.CascadeDepth

		err = validateCascadePlanOptions(options)
		if err != nil {
			return nil, errors.New(errContext, err.Error())
		}
	}

	plan, err = h.planFactory.NewPlan(planType, planParameters)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return plan, nil

}

// validateCascadePlanOptions returns an error if the options are not valid for cascade plan
func validateCascadePlanOptions(options *HandlerOptions) error {
	errContext := "(build::validateCascadePlanOptions)"

	if options.AnsibleIntermediateContainerName != "" {
		return errors.New(errContext, "Cascade plan does not support intermediate containers name. It could cause an unpredictable result")
	}

	if options.AnsibleInventoryPath != "" {
		return errors.New(errContext, "Cascade plan does not support ansible inventory path. It could cause an unpredictable result")
	}

	if options.AnsibleLimit != "" {
		return errors.New(errContext, "Cascade plan does not support ansible limit. It could cause an unpredictable result")
	}

	return nil
}
