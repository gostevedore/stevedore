package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/engine/promote"
)

type Handler struct {
	service ServicePromoter
}

func NewHandler(p ServicePromoter) *Handler {
	return &Handler{
		service: p,
	}
}

func (h *Handler) Handler(ctx context.Context, options *HandlerOptions) error {
	//func (h *Handler) Handler(ctx context.Context, options *HandlerOptions) func(cmd *cobra.Command, args []string) error {
	//	return func(cmd *cobra.Command, args []string) error {

	errContext := "(promote::Handler)"

	if options.SourceImageName == "" {
		return errors.New(errContext, "Source images name must be provided")
	}

	serviceOptions := &promote.ServiceOptions{
		SourceImageName:       options.SourceImageName,
		RemoveTargetImageTags: options.RemoveTargetImageTags,
	}

	if options.TargetImageName != "" {
		serviceOptions.TargetImageName = options.TargetImageName
	}

	if options.TargetImageRegistryNamespace != "" {
		serviceOptions.TargetImageRegistryNamespace = options.TargetImageRegistryNamespace
	}

	if options.TargetImageRegistryHost != "" {
		serviceOptions.TargetImageRegistryHost = options.TargetImageRegistryHost
	}

	if len(options.TargetImageTags) > 0 {
		serviceOptions.TargetImageTags = options.TargetImageTags
	}

	serviceOptions.PromoteSourceImageTag = options.PromoteSourceImageTag
	serviceOptions.EnableSemanticVersionTags = options.EnableSemanticVersionTags
	serviceOptions.RemoveTargetImageTags = options.RemoveTargetImageTags
	serviceOptions.RemoteSourceImage = options.RemoteSourceImage
	serviceOptions.SemanticVersionTagsTemplates = options.SemanticVersionTagsTemplates

	promoteType := "docker"
	if options.DryRun {
		promoteType = "dry-run"
	}

	err := h.service.Promote(ctx, serviceOptions, promoteType)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
