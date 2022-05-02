package build

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/service/build"
	"github.com/gostevedore/stevedore/internal/service/build/plan"
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
func (h *Handler) Handler(ctx context.Context, imageName string, options *Options) error {

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

	buildServiceOptions.DryRun = options.DryRun

	buildServiceOptions.EnableSemanticVersionTags = options.EnableSemanticVersionTags

	buildServiceOptions.ImageFromName = options.ImageFromName
	buildServiceOptions.ImageFromRegistryHost = options.ImageFromRegistryHost
	buildServiceOptions.ImageFromRegistryNamespace = options.ImageFromRegistryNamespace
	buildServiceOptions.ImageFromVersion = options.ImageFromVersion

	buildServiceOptions.ImageName = options.ImageName
	buildServiceOptions.ImageRegistryHost = options.ImageRegistryHost
	buildServiceOptions.ImageRegistryNamespace = options.ImageRegistryNamespace
	buildServiceOptions.ImageVersions = append([]string{}, options.Versions...)

	buildServiceOptions.Labels = make(map[string]string)
	for _, label := range options.Labels {

		if strings.IndexRune(label, AssignmentTokenSymbol) < 0 {
			return errors.New(errContext, fmt.Sprintf("Invalid label format '%s'", label))
		}
		kLabel := label[:strings.IndexRune(label, AssignmentTokenSymbol)]
		vLabel := label[strings.IndexRune(label, AssignmentTokenSymbol)+1:]

		buildServiceOptions.Labels[kLabel] = vLabel
	}

	buildServiceOptions.PersistentVars = make(map[string]interface{})
	for _, persistentVars := range options.PersistentVars {

		if strings.IndexRune(persistentVars, AssignmentTokenSymbol) < 0 {
			return errors.New(errContext, fmt.Sprintf("Invalid persistent variable format '%s'", persistentVars))
		}
		kPVar := persistentVars[:strings.IndexRune(persistentVars, AssignmentTokenSymbol)]
		vPVar := persistentVars[strings.IndexRune(persistentVars, AssignmentTokenSymbol)+1:]

		buildServiceOptions.PersistentVars[kPVar] = vPVar
	}

	buildServiceOptions.PullParentImage = options.PullParentImage
	buildServiceOptions.PushImageAfterBuild = options.PushImagesAfterBuild
	buildServiceOptions.RemoveImagesAfterPush = options.RemoveImagesAfterPush

	buildPlan, err = h.createBuildPlan(options)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = h.service.Build(
		ctx,
		buildPlan,
		imageName,
		options.Versions,
		buildServiceOptions,
	)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

func (h *Handler) createBuildPlan(options *Options) (plan.Planner, error) {
	errContext := "(build::createBuildPlan)"

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
		err = validateCascadePlanOptions(options)
		if err != nil {
			return nil, errors.New(errContext, err.Error())
		}

		planType = "cascade"
		planParameters["depth"] = options.CascadeDepth
	}

	plan, err = h.planFactory.NewPlan(planType, planParameters)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return plan, nil

}

// validateCascadePlanOptions returns an error if the options are not valid for cascade plan
func validateCascadePlanOptions(options *Options) error {
	errContext := "(build::validateCascadePlanOptions)"

	if options == nil {
		return errors.New(errContext, "Options to be validated are required")
	}

	if options.AnsibleIntermediateContainerName != "" {
		return errors.New(errContext, "Cascade plan does not support intermediate containers name, it could cause an unpredictable result")
	}

	if options.AnsibleInventoryPath != "" {
		return errors.New(errContext, "Cascade plan does not support ansible inventory path, it could cause an unpredictable result")
	}

	if options.AnsibleLimit != "" {
		return errors.New(errContext, "Cascade plan does not support ansible limit, it could cause an unpredictable result")
	}

	if options.ImageName != "" {
		return errors.New(errContext, "Cascade plan does not support image name, it could cause an unpredictable result")
	}

	if options.ImageFromName != "" {
		return errors.New(errContext, "Cascade plan does not support image from name, it could cause an unpredictable result")
	}

	return nil
}
