package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/application/promote"
)

type Handler struct {
	service ServicePromoter
}

func NewHandler(p ServicePromoter) *Handler {
	return &Handler{
		service: p,
	}
}

func (h *Handler) Handler(ctx context.Context, options *Options) error {

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

	serviceOptions.DryRun = options.DryRun
	serviceOptions.PromoteSourceImageTag = options.PromoteSourceImageTag
	serviceOptions.EnableSemanticVersionTags = options.EnableSemanticVersionTags
	serviceOptions.RemoveTargetImageTags = options.RemoveTargetImageTags
	serviceOptions.RemoteSourceImage = options.RemoteSourceImage
	serviceOptions.SemanticVersionTagsTemplates = options.SemanticVersionTagsTemplates

	err := h.service.Promote(ctx, serviceOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
